package client

import (
	"context"
	"errors"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/drand/drand/chain"
	"github.com/drand/drand/log"
)

const (
	defaultRequestTimeout     = time.Second * 5
	defaultSpeedTestInterval  = time.Minute * 5
	defaultRequestConcurrency = 2
)

// NewOptimizingClient creates a drand client that measures the speed of clients
// and uses the fastest ones.
//
// Clients passed to the optimising client are ordered by speed and calls to
// `Get` race the 2 fastest clients (by default) for the result. If a client
// errors then it is moved to the back of the list.
//
// A speed test is performed periodically in the background every 5 minutes (by
// default) to ensure we're still using the fastest clients. A negative speed
// test interval will disable testing.
//
// Calls to `Get` actually iterate over the speed-ordered client list with a
// concurrency of 2 (by default) until a result is retrieved. It means that the
// optimising client will fallback to using the other slower clients in the
// event of failure(s).
//
// Additionally, calls to Get are given a timeout of 5 seconds (by default) to
// ensure no unbounded blocking occurs.
func NewOptimizingClient(
	clients []Client,
	chainInfo *chain.Info,
	requestTimeout time.Duration,
	requestConcurrency int,
	speedTestInterval time.Duration,
) (Client, error) {
	if len(clients) == 0 {
		return nil, errors.New("missing clients")
	}
	if chainInfo == nil {
		return nil, errors.New("missing chain info")
	}
	stats := make([]*requestStat, len(clients))
	now := time.Now()
	for i, c := range clients {
		stats[i] = &requestStat{client: c, rtt: 0, startTime: now}
	}
	done := make(chan struct{})
	if requestTimeout <= 0 {
		requestTimeout = defaultRequestTimeout
	}
	if requestConcurrency <= 0 {
		requestConcurrency = defaultRequestConcurrency
	}
	if speedTestInterval == 0 {
		speedTestInterval = defaultSpeedTestInterval
	}
	oc := &optimizingClient{
		clients:            clients,
		stats:              stats,
		requestTimeout:     requestTimeout,
		requestConcurrency: requestConcurrency,
		speedTestInterval:  speedTestInterval,
		chainInfo:          chainInfo,
		log:                log.DefaultLogger,
		done:               done,
	}
	if speedTestInterval > 0 {
		go oc.testSpeed()
	}
	return oc, nil
}

type optimizingClient struct {
	sync.RWMutex
	clients            []Client
	stats              []*requestStat
	requestTimeout     time.Duration
	requestConcurrency int
	speedTestInterval  time.Duration
	chainInfo          *chain.Info
	log                log.Logger
	done               chan struct{}
}

type requestStat struct {
	// client is the client used to make the request.
	client Client
	// rtt is the time it took to make the request.
	rtt time.Duration
	// startTime is the time at which the request was started.
	startTime time.Time
}

type requestResult struct {
	// result is the return value from the call to Get.
	result Result
	// err is the error that occurred from a call to Get (not including context error).
	err error
	// stat is stats from the call to Get.
	stat *requestStat
}

func (oc *optimizingClient) testSpeed() {
	for {
		stats := []*requestStat{}
		ctx, cancel := context.WithCancel(context.Background())
		ch := parallelGet(ctx, oc.clients, 1, oc.requestTimeout, oc.requestConcurrency)

	LOOP:
		for {
			select {
			case rr, ok := <-ch:
				if !ok {
					cancel()
					break LOOP
				}
				stats = append(stats, rr.stat)
			case <-oc.done:
				cancel()
				return
			}
		}

		oc.updateStats(stats)

		t := time.NewTimer(oc.speedTestInterval)
		select {
		case <-t.C:
		case <-oc.done:
			t.Stop()
			return
		}
	}
}

// SetLog configures the client log output.
func (oc *optimizingClient) SetLog(l log.Logger) {
	oc.log = l
}

// Get returns the randomness at `round` or an error.
func (oc *optimizingClient) Get(ctx context.Context, round uint64) (res Result, err error) {
	oc.RLock()
	// copy the current ordered client list so we iterate over a stable slice
	var clients []Client
	for _, s := range oc.stats {
		clients = append(clients, s.client)
	}
	oc.RUnlock()

	stats := []*requestStat{}
	defer oc.updateStats(stats)
	ch := raceGet(ctx, clients, round, oc.requestTimeout, oc.requestConcurrency)

LOOP:
	for {
		select {
		case rr, ok := <-ch:
			if !ok {
				break LOOP
			}
			stats = append(stats, rr.stat)
			res, err = rr.result, rr.err
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-oc.done:
			return nil, errors.New("client closed")
		}
	}

	return
}

// get calls Get on the passed client and returns a channel that yields a single
// raceResult when the call completes or closes the channel with no result if the
// context is canceled.
func get(ctx context.Context, client Client, round uint64) <-chan *requestResult {
	ch := make(chan *requestResult, 1)
	go func() {
		start := time.Now()
		res, err := client.Get(ctx, round)
		rtt := time.Now().Sub(start)

		if ctx.Err() != nil {
			close(ch)
			return
		}

		// client failure, set a large RTT so it is sent to the back of the list
		if err != nil {
			rtt = math.MaxInt64
		}

		stat := requestStat{client, rtt, start}

		if err != nil {
			ch <- &requestResult{nil, err, &stat}
		} else {
			ch <- &requestResult{res, nil, &stat}
		}
		close(ch)
	}()
	return ch
}

func raceGet(ctx context.Context, clients []Client, round uint64, timeout time.Duration, concurrency int) <-chan *requestResult {
	results := make(chan *requestResult, len(clients))

	go func() {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		defer close(results)
		ch := parallelGet(ctx, clients, round, timeout, concurrency)

		for {
			select {
			case rr, ok := <-ch:
				if !ok {
					return
				}
				results <- rr
				if rr.err == nil { // race is won
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return results
}

func parallelGet(ctx context.Context, clients []Client, round uint64, timeout time.Duration, concurrency int) <-chan *requestResult {
	results := make(chan *requestResult, len(clients))
	token := make(chan struct{}, concurrency)

	for i := 0; i < concurrency; i++ {
		token <- struct{}{}
	}

	go func() {
		wg := sync.WaitGroup{}
	LOOP:
		for len(clients) > 0 {
			c := clients[0]
			clients = clients[1:]
			select {
			case <-token:
				wg.Add(1)
				go func(c Client) {
					defer func() { token <- struct{}{} }()
					defer wg.Done()

					ctx, cancel := context.WithTimeout(ctx, timeout)
					defer cancel()

					ch := get(ctx, c, round)
					select {
					case rr, ok := <-ch:
						if !ok {
							return
						}
						results <- rr
					case <-ctx.Done():
						return
					}
				}(c)
			case <-ctx.Done():
				break LOOP
			}
		}
		wg.Wait()
		close(results)
	}()

	return results
}

func (oc *optimizingClient) updateStats(stats []*requestStat) {
	oc.Lock()
	defer oc.Unlock()

	// update the round trip times with new samples
	for _, next := range stats {
		for _, curr := range oc.stats {
			if curr.client == next.client {
				if curr.startTime.Before(next.startTime) {
					curr.rtt = next.rtt
					curr.startTime = next.startTime
				}
				break
			}
		}
	}

	// sort by fastest
	sort.Slice(oc.stats, func(i, j int) bool {
		return oc.stats[i].rtt < oc.stats[j].rtt
	})
}

// Watch returns new randomness as it becomes available.
func (oc *optimizingClient) Watch(ctx context.Context) <-chan Result {
	return pollingWatcher(ctx, oc, oc.chainInfo, oc.log)
}

// RoundAt will return the most recent round of randomness that will be available
// at time for the current client.
func (oc *optimizingClient) RoundAt(time time.Time) uint64 {
	return oc.clients[0].RoundAt(time)
}

// Close stops the background speed tests and closes the client for further use.
func (oc *optimizingClient) Close() {
	close(oc.done)
}
