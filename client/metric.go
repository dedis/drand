package client

import (
	"context"
	"time"

	"github.com/drand/drand/chain"
	"github.com/drand/drand/metrics"
)

func newWatchLatencyMetricClient(base Client, info *chain.Info) Client {
	c := &watchLatencyMetricClient{
		Client:    base,
		chainInfo: info,
	}
	go c.startObserve(context.Background())
	return c
}

type watchLatencyMetricClient struct {
	Client
	chainInfo *chain.Info
}

func (c *watchLatencyMetricClient) startObserve(ctx context.Context) {
	rch := c.Watch(ctx)
	for {
		select {
		case result, ok := <-rch:
			if !ok {
				return
			}
			// compute the latency metric
			actual := time.Now().Unix()
			expected := chain.TimeOfRound(c.chainInfo.Period, c.chainInfo.GenesisTime, result.Round())
			// the labels of the gauge vec must already be set at the registerer level
			metrics.ClientWatchLatency.Set(float64(actual - expected))
		case <-ctx.Done():
			return
		}
	}
}
