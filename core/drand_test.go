package core

import (
	"context"
	"fmt"
	"io/ioutil"
	gnet "net"
	"os"
	"path"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"

	"github.com/BurntSushi/toml"
	"github.com/drand/drand/beacon"
	"github.com/drand/drand/key"
	"github.com/drand/drand/log"
	"github.com/drand/drand/net"
	"github.com/drand/drand/protobuf/drand"
	"github.com/drand/drand/test"

	//"github.com/drand/kyber"
	clock "github.com/jonboulle/clockwork"
	"github.com/kabukky/httpscerts"
	"github.com/stretchr/testify/require"
)

var testBeaconOffset = "2s"
var testDkgTimeout = "2s"

func TestDrandDKGFresh(t *testing.T) {
	n := 4
	beaconPeriod := 1 * time.Second
	var offsetGenesis = 1 * time.Second
	genesis := clock.NewFakeClock().Now().Add(offsetGenesis).Unix()
	dt := NewDrandTest(t, n, key.DefaultThreshold(n), beaconPeriod, genesis)
	defer dt.Cleanup()
	dt.RunDKG()
	fmt.Println(" --- DKG FINISHED ---")
	// make the last node fail
	lastID := dt.ids[n-1]
	dt.StopDrand(lastID)
	//lastOne.Stop()
	fmt.Printf("\n--- lastOne STOPPED --- \n\n")

	// move time to genesis
	dt.MoveTime(offsetGenesis)
	// two = genesis + 1st round (happens at genesis)
	dt.TestBeaconLength(2, dt.ids[:n-1]...)
	fmt.Println(" --- Test BEACON LENGTH --- ")
	// start last one
	dt.StartDrand(lastID, true)
	// leave some room to do the catchup
	time.Sleep(100 * time.Millisecond)
	fmt.Println(" --- STARTED BEACON DRAND ---")
	dt.MoveTime(beaconPeriod)
	dt.TestBeaconLength(3, dt.ids...)
	dt.TestPublicBeacon(dt.ids[0])
}

func TestDrandDKGReshareTimeout(t *testing.T) {
	oldN := 4
	newN := 4
	oldThr := 3
	newThr := 3
	timeoutStr := "1s"
	timeout, _ := time.ParseDuration(timeoutStr)
	beaconPeriod := 2 * time.Second
	// starts at 2sec
	var offsetGenesis = 2 * time.Second
	genesisTime := clock.NewFakeClock().Now().Add(offsetGenesis).Unix()
	// after 2 beacon period, we switch to new network
	offsetTransition := beaconPeriod * 2
	transitionTime := clock.NewFakeClock().Now().Add(offsetGenesis).Add(offsetTransition).Unix()
	offline := 1

	dt := NewDrandTest(t, oldN, oldThr, beaconPeriod, genesisTime)
	defer dt.Cleanup()
	dt.RunDKG()

	// move to genesis time - so nodes start to make a round
	dt.MoveTime(offsetGenesis)
	// two = genesis + 1st round (happens at genesis)
	dt.TestBeaconLength(2, dt.ids...)

	// + offline makes sure t
	toKeep := oldN - offline
	toAdd := newN - toKeep
	dt.SetupReshare(toKeep, toAdd, newThr, transitionTime)

	fmt.Println("SETUP RESHARE DONE")
	//// run the resharing
	var doneReshare = make(chan bool)
	go func() {
		dt.RunReshare(toKeep, toAdd, timeoutStr)
		doneReshare <- true
	}()
	checkDone := func() bool {
		select {
		case <-doneReshare:
			return true
		default:
			return false
		}
	}
	// check it is not done yet
	time.Sleep(getSleepDuration())
	require.False(t, checkDone())

	// leave time for the dkg to proceed
	// for go routines and such
	time.Sleep(getSleepDuration())
	// advance time to the timeout
	fmt.Printf("\n -- Move to timeout time !! -- \n")
	// after this, the resharing should be done
	dt.MoveTime(timeout)
	//give time to finish for the go routines and such
	time.Sleep(getSleepDuration())
	fmt.Printf("\n -- Move to timeout time  AFTER - genesis: %d - current %d !! -- \n", genesisTime, dt.myClock.Now().Unix())
	require.True(t, checkDone())

	fmt.Printf("\n--- TEST RESHARING: Move time to 1sec\n\n")

	// move 1 second to pass to the next round time
	// it is still the "old network - offline" that is running this
	dt.MoveTime(1 * time.Second)
	// genesisTime: genesis + round 1
	// genesisTime + beaconPeriod: round 2
	dt.TestBeaconLength(3, dt.ids...)

	fmt.Printf("\n--- TEST RESHARING: Move time to Transition phase \n\n")

	// move 1 second to come to the transition phase: transition time - 1s is
	// where old nodes stop and new nodes takes over
	dt.MoveTime(1 * time.Second)
	// no new beacon created yet, but we check only for the nodes that are in
	// the new group
	dt.TestBeaconLength(3, dt.reshareIds...)

	fmt.Printf("\n -- TEST RESHARING: New group round time \n\n")
	// move 1s to pass to first round of the new group
	dt.MoveTime(1 * time.Second)
	dt.TestBeaconLength(4, dt.reshareIds...)
}

type DrandTest struct {
	sync.Mutex
	t            *testing.T
	n            int
	thr          int
	dir          string
	newDir       string
	drands       map[string]*Drand
	newDrands    map[string]*Drand
	group        *key.Group
	newGroup     *key.Group
	groupPath    string
	newGroupPath string
	period       time.Duration
	genesis      int64
	ids          []string
	newIds       []string
	reshareIds   []string
	shares       map[string]*key.Share
	myClock      clock.FakeClock
	certPaths    []string
	newCertPaths []string
}

func (d *DrandTest) Cleanup() {
	os.RemoveAll(d.dir)
	os.RemoveAll(d.newDir)
}

func (d *DrandTest) GetBeacon(id string, round int) (*beacon.Beacon, error) {
	dd, ok := d.drands[id]
	require.True(d.t, ok)

	return dd.beacon.Store().Get(uint64(round))
}

// returns new ids generated
func (d *DrandTest) SetupReshare(keepOld, addNew, newThr int, transitionTime int64) []string {
	newN := keepOld + addNew
	ids := make([]*key.Identity, 0, newN)
	newAddr := make([]string, addNew)
	newDrands, _, newDir, newCertPaths := BatchNewDrand(addNew, false,
		WithCallOption(grpc.FailFast(true)), WithLogLevel(log.LogDebug))
	d.newDir = newDir
	d.newDrands = make(map[string]*Drand)
	// add old participants
	for _, id := range d.ids[:keepOld] {
		drand := d.drands[id]
		ids = append(ids, drand.priv.Public)
		for _, cp := range newCertPaths {
			drand.opts.certmanager.Add(cp)
		}

	}
	// add new participants
	for i, drand := range newDrands {
		ids = append(ids, drand.priv.Public)
		newAddr[i] = drand.priv.Public.Address()
		d.newDrands[drand.priv.Public.Address()] = drand
		d.setClock(newAddr[i])
		for _, cp := range d.certPaths {
			drand.opts.certmanager.Add(cp)
		}
	}

	d.newIds = newAddr

	d.newGroup = key.NewGroup(ids, newThr, d.group.GenesisTime)
	d.newGroup.Period = d.period
	d.newGroup.TransitionTime = transitionTime
	d.newGroup.GenesisSeed = d.group.GenesisSeed
	fmt.Println("RESHARE GROUP:\n", d.newGroup.String())
	d.newGroupPath = path.Join(newDir, "newgroup.toml")
	require.NoError(d.t, key.Save(d.newGroupPath, d.newGroup, false))
	return newAddr
}

func (d *DrandTest) RunReshare(oldRun, newRun int, timeout string) {
	fmt.Printf(" -- Running RESHARE with %d/%d old, %d/%d new nodes\n", oldRun, len(d.drands), newRun, len(d.newIds))
	var clientCounter = &sync.WaitGroup{}
	var secret = "thisistheresharing"

	// take list of old nodes present in new groups
	var oldNodes []string
	var oldLeaving []string
	for _, id := range d.ids {
		drand := d.drands[id]
		if d.newGroup.Contains(drand.priv.Public) {
			oldNodes = append(oldNodes, drand.priv.Public.Address())
		} else {
			oldLeaving = append(oldLeaving, id)
		}
	}

	var allIds []string
	// stop the old nodes that we want offline during the resharing
	for _, id := range oldNodes[oldRun:] {
		fmt.Printf("Stopping old-new node %d - %s\n", d.drands[id].index, id)
		d.drands[id].Stop()
	}
	// stop the old nodes that are leaving the group
	for _, id := range oldLeaving {
		fmt.Printf("Stopping old-leaving node %d - %s\n", d.drands[id].index, id)
		d.drands[id].Stop()
	}
	// run the old ones
	// exclude leader
	clientCounter.Add(oldRun - 1)
	for _, id := range oldNodes[1:oldRun] {
		dr := d.drands[id]
		idx, found := d.newGroup.Index(dr.priv.Public)
		if !found {
			panic("old drand not found")
		}
		fmt.Printf("Launching reshare on new %d - %s\n", idx, id)
		go runreshare(d.drands[id], false)
		allIds = append(allIds, id)
	}

	// run the new ones
	clientCounter.Add(newRun)
	for _, id := range d.newIds[:newRun] {
		dr := d.newDrands[id]
		idx, found := d.newGroup.Index(dr.priv.Public)
		if !found {
			panic("new drand not found")
		}
		fmt.Printf("Launching reshare on new  %d - %s\n", idx, id)
		go runreshare(d.newDrands[id], false)
		allIds = append(allIds, id)
	}
	allIds = append(allIds, oldNodes[0])
	d.setDKGCallback(allIds)
	d.reshareIds = allIds
	runreshare := func(dr *Drand, leader bool) {
		// instruct to be ready for a reshare
		client, err := net.NewControlClient(dr.opts.controlPort)
		require.NoError(d.t, err)
		_, err = client.InitReshare(leader, len(allIds), d.newGroup.Threshold, testBeaconOffset, timeout, secret, d.groupPath)
		require.NoError(d.t, err)
		fmt.Printf("\n\nDKG TEST: drand %s DONE RESHARING (leader? %v)\n", dr.priv.Public.Address(), leader)
		clientCounter.Done()
	}

	// run leader
	leader := d.drands[oldNodes[0]]
	idx, found := d.newGroup.Index(leader.priv.Public)
	if !found {
		panic("leader not found")
	}
	fmt.Printf("Launching reshare on (old) root %d - %s\n", idx, oldNodes[0])
	clientCounter.Add(1)
	go runreshare(d.drands[oldNodes[0]], true)
	// wait for the return of the clients
	fmt.Println("\n\n -- Waiting COUNTER for ", oldRun-1+newRun+1, " nodes --")
	checkWait(clientCounter)
	fmt.Printf("\n\n -- TEST FINISHED ALL RESHARE DKG --\n\n")
}

func checkWait(counter *sync.WaitGroup) {
	var doneCh = make(chan bool, 1)
	go func() {
		counter.Wait()
		doneCh <- true
	}()
	select {
	case <-doneCh:
		break
	case <-time.After(11 * time.Second):
		panic("outdated beacon time")
	}
}
func NewDrandTest(t *testing.T, n, thr int, period time.Duration, genesis int64) *DrandTest {
	drands, group, dir, certPaths := BatchNewDrand(n, false,
		WithCallOption(grpc.FailFast(true)),
	)
	group.Period = period
	group.GenesisTime = genesis
	groupPath := path.Join(dir, "dkggroup.toml")
	require.NoError(t, key.Save(groupPath, group, false))
	myClock := clock.NewFakeClock()
	ids := make([]string, n)
	mDrands := make(map[string]*Drand, n)
	for i, d := range drands {
		ids[i] = d.priv.Public.Address()
		mDrands[ids[i]] = d
	}
	return &DrandTest{
		t:         t,
		n:         n,
		thr:       thr,
		drands:    mDrands,
		group:     group,
		groupPath: groupPath,
		period:    period,
		ids:       ids,
		shares:    make(map[string]*key.Share),
		myClock:   myClock,
		certPaths: certPaths,
	}
}

func (d *DrandTest) RunDKG() {
	var wg sync.WaitGroup
	wg.Add(d.n - 1)
	d.setClock(d.ids...)
	d.setDKGCallback(d.ids)
	for _, id := range d.ids[1:] {
		go func(dd *Drand) {
			client, err := net.NewControlClient(dd.opts.controlPort)
			require.NoError(d.t, err)
			_, err = client.InitDKG(d.groupPath, false, "", nil)
			require.NoError(d.t, err)
			wg.Done()
			fmt.Printf("\n\n\n TESTDKG NON-ROOT %s FINISHED\n\n\n", dd.priv.Public.Address())
		}(d.drands[id])
	}

	root := d.drands[d.ids[0]]
	controlClient, err := net.NewControlClient(root.opts.controlPort)
	require.NoError(d.t, err)
	_, err = controlClient.InitDKG(d.groupPath, true, "", nil)
	require.NoError(d.t, err)
	wg.Wait()
	fmt.Printf("\n\n\n TESTDKG ROOT %s FINISHED\n\n\n", d.ids[0])
	resp, err := controlClient.GroupFile()
	require.NoError(d.t, err)
	group := new(key.Group)
	groupToml := new(key.GroupTOML)
	_, err = toml.Decode(resp.GetGroupToml(), groupToml)
	require.NoError(d.t, err)
	require.NoError(d.t, group.FromTOML(groupToml))
	d.group = group
	require.Equal(d.t, d.thr, d.group.Threshold)
	for _, drand := range d.drands {
		require.True(d.t, d.group.Contains(drand.priv.Public))
	}
	require.Len(d.t, d.group.PublicKey.Coefficients, d.thr)
	require.NoError(d.t, key.Save(d.groupPath, d.group, false))
}

func (d *DrandTest) tryBoth(id string, fn func(d *Drand)) {
	if dr, ok := d.drands[id]; ok {
		fn(dr)
	} else if dr, ok = d.newDrands[id]; ok {
		fn(dr)
	} else {
		panic("that should not happen")
	}
}

func (d *DrandTest) setClock(ids ...string) {
	now := d.myClock.Now()
	for _, id := range ids {
		d.tryBoth(id, func(dr *Drand) {
			addr := dr.priv.Public.Address()
			clock := clock.NewFakeClockAt(now)
			dr.opts.clock = clock
			dr.opts.dkgCallback = func(s *key.Share) {
				d.shares[addr] = s
				fmt.Printf("\n\n\n  --- DKG %s FINISHED ---\n\n\n", addr)
			}
		})
	}
}

// first wait for all dkg callbacks to trigger, then update the clock of every
// ids
func (d *DrandTest) setDKGCallback(ids []string) {
	for _, id := range ids {
		d.tryBoth(id, func(dr *Drand) {
			dr.opts.dkgCallback = func(s *key.Share) {
				d.Lock()
				id := dr.priv.Public.Address()
				d.shares[id] = s
				d.Unlock()
				//fmt.Printf("\n\nDKG DONE for %d - %s\n\n", dr.index, id)
			}
		})
	}
}

func (d *DrandTest) GetDrand(id string) *Drand {
	return d.drands[id]
}

func (d *DrandTest) StopDrand(id string) {
	dr := d.drands[id]
	dr.Stop()
	pinger, err := net.NewControlClient(dr.opts.controlPort)
	require.NoError(d.t, err)
	var counter = 1
	fmt.Println(" DRAND ", dr.priv.Public.Address(), " TRYING TO PING")
	for range time.Tick(100 * time.Millisecond) {
		if err := pinger.Ping(); err != nil {
			fmt.Println(" DRAND ", dr.priv.Public.Address(), " TRYING TO PING DONE")
			break
		}
		counter++
		require.LessOrEqual(d.t, counter, 5)
	}
	fmt.Println(" DRAND ", dr.priv.Public.Address(), " STOPPED")
}

func (d *DrandTest) StartDrand(id string, catchup bool) {
	dr, ok := d.drands[id]
	require.True(d.t, ok)
	var err error
	dr, err = LoadDrand(dr.store, dr.opts)
	require.NoError(d.t, err)
	d.drands[id] = dr
	d.setClock(id)
	fmt.Println("--- JUST BEFORE STARTBEACON---")
	dr.StartBeacon(catchup)
	fmt.Println("--- JUST AFTER STARTBEACON---")
}

func (d *DrandTest) MoveTime(p time.Duration) {
	for _, d := range d.drands {
		c := d.opts.clock.(clock.FakeClock)
		c.Advance(p)
	}
	for _, d := range d.newDrands {
		c := d.opts.clock.(clock.FakeClock)
		c.Advance(p)
	}
	d.myClock.Advance(p)
	fmt.Printf(" --- MoveTime: new time is %d \n", d.myClock.Now().Unix())
	time.Sleep(getSleepDuration())
}

func (d *DrandTest) TestBeaconLength(length int, ids ...string) {
	fmt.Printf("\n BEACON LENGTH tests (should be %d):\n", length)
	for _, id := range ids {
		d.tryBoth(id, func(drand *Drand) {
			drand.state.Lock()
			defer drand.state.Unlock()
			if drand.beacon == nil {
				// this drand has stopped for a reason
				return
			}
			fmt.Printf("\n\tTest %s (beacon %p)\n", id, drand.beacon)
			howMany := 0
			drand.beacon.Store().Cursor(func(c beacon.Cursor) {
				for b := c.First(); b != nil; b = c.Next() {
					howMany++
					fmt.Printf("\t %d - %s: beacon %s\n", drand.index, drand.priv.Public.Address(), b)
				}
			})
			require.Equal(d.t, length, drand.beacon.Store().Len(), "id %s - howMany is %d vs Len() %d", id, howMany, drand.beacon.Store().Len())
		})
	}

}

func (d *DrandTest) TestPublicBeacon(id string) {
	dr := d.GetDrand(id)
	client := net.NewGrpcClientFromCertManager(dr.opts.certmanager, dr.opts.grpcOpts...)
	resp, err := client.PublicRand(context.TODO(), test.NewTLSPeer(dr.priv.Public.Addr), &drand.PublicRandRequest{})
	require.NoError(d.t, err)
	require.NotNil(d.t, resp)
}

// Check they all have same public group file after dkg
func TestDrandPublicGroup(t *testing.T) {
	n := 10
	thr := key.DefaultThreshold(n)
	p := 200 * time.Millisecond
	genesisTime := clock.NewFakeClock().Now().Unix()
	dt := NewDrandTest(t, n, thr, p, genesisTime)
	defer dt.Cleanup()
	dt.RunDKG()

	//client := NewGrpcClient()
	cm := dt.drands[dt.ids[0]].opts.certmanager
	client := NewGrpcClientFromCert(cm)
	rest := net.NewRestClientFromCertManager(cm)
	var group *drand.GroupResponse
	for i, id := range dt.ids {
		d := dt.drands[id]
		groupResp, err := client.Group(d.priv.Public.Address(), d.priv.Public.TLS)
		require.NoError(t, err, fmt.Sprintf("idx %d: addr %s", i, id))
		if group == nil {
			group = groupResp
		}
		require.Equal(t, uint32(group.Period), groupResp.Period)
		require.Equal(t, uint32(group.Threshold), groupResp.Threshold)
		require.Equal(t, group.Distkey, groupResp.Distkey)
		require.Len(t, groupResp.Nodes, len(group.Nodes))

		nodes := groupResp.GetNodes()
		for addr, d := range dt.drands {
			var found bool
			for _, n := range nodes {
				sameAddr := n.GetAddress() == addr
				sameKey := n.GetKey() == key.PointToString(d.priv.Public.Key)
				sameTLS := n.GetTLS() == d.priv.Public.TLS
				if sameAddr && sameKey && sameTLS {
					found = true
					break
				}
			}
			require.True(t, found)
		}
		restGroup, err := rest.Group(context.TODO(), d.priv.Public, &drand.GroupRequest{})
		require.NoError(t, err)
		require.Equal(t, groupResp, restGroup)
	}
}

// BatchNewDrand returns n drands, using TLS or not, with the given
// options. It returns the list of Drand structures, the group created,
// the folder where db, certificates, etc are stored. It is the folder
// to delete at the end of the test. As well, it returns a public grpc
// client that can reach any drand node.
func BatchNewDrand(n int, insecure bool, opts ...ConfigOption) ([]*Drand, *key.Group, string, []string) {
	var privs []*key.Pair
	var group *key.Group
	if insecure {
		privs, group = test.BatchIdentities(n)
	} else {
		privs, group = test.BatchTLSIdentities(n)
	}
	ports := test.Ports(n)
	var err error
	drands := make([]*Drand, n, n)
	tmp := os.TempDir()
	dir, err := ioutil.TempDir(tmp, "drand")
	if err != nil {
		panic(err)
	}

	certPaths := make([]string, n)
	keyPaths := make([]string, n)
	if !insecure {
		for i := 0; i < n; i++ {
			certPath := path.Join(dir, fmt.Sprintf("server-%d.crt", i))
			keyPath := path.Join(dir, fmt.Sprintf("server-%d.key", i))
			if httpscerts.Check(certPath, keyPath) != nil {

				h, _, err := gnet.SplitHostPort(privs[i].Public.Address())
				if err != nil {
					panic(err)
				}
				if err := httpscerts.Generate(certPath, keyPath, h); err != nil {
					panic(err)
				}
			}
			certPaths[i] = certPath
			keyPaths[i] = keyPath
		}
	}

	for i := 0; i < n; i++ {
		s := test.NewKeyStore()
		s.SaveKeyPair(privs[i])
		// give each one their own private folder
		dbFolder := path.Join(dir, fmt.Sprintf("db-%d", i))
		confOptions := []ConfigOption{WithDbFolder(dbFolder)}
		if !insecure {
			confOptions = append(confOptions, WithTLS(certPaths[i], keyPaths[i]))
			confOptions = append(confOptions, WithTrustedCerts(certPaths...))
		} else {
			confOptions = append(confOptions, WithInsecure())
		}
		confOptions = append(confOptions, WithControlPort(ports[i]))
		confOptions = append(confOptions, WithLogLevel(log.LogDebug))
		// add options in last so it overwrites the default
		confOptions = append(confOptions, opts...)
		drands[i], err = NewDrand(s, NewConfig(confOptions...))
		if err != nil {
			panic(err)
		}
	}
	return drands, group, dir, certPaths
}

// CloseAllDrands closes all drands
func CloseAllDrands(drands []*Drand) {
	for i := 0; i < len(drands); i++ {
		drands[i].Stop()
		//os.RemoveAll(drands[i].opts.dbFolder)
	}
}

func getSleepDuration() time.Duration {
	if os.Getenv("TRAVIS_BRANCH") != "" {
		return time.Duration(3000) * time.Millisecond
	}
	return time.Duration(500) * time.Millisecond
}
