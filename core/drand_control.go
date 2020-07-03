package core

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/drand/drand/entropy"
	"github.com/drand/drand/key"
	dnet "github.com/drand/drand/net"
	"github.com/drand/drand/protobuf/drand"
	"github.com/drand/kyber/share/dkg"
	vss "github.com/drand/kyber/share/vss/pedersen"
)

// errPreempted is returned on reshares when a subsequent reshare is started concurrently
var errPreempted = errors.New("time out: pre-empted")

// InitDKG take a InitDKGPacket, extracts the informations needed and wait for the
// DKG protocol to finish. If the request specifies this node is a leader, it
// starts the DKG protocol.
func (d *Drand) InitDKG(c context.Context, in *drand.InitDKGPacket) (*drand.GroupPacket, error) {
	isLeader := in.GetInfo().GetLeader()
	d.state.Lock()
	if d.dkgDone {
		d.state.Unlock()
		return nil, errors.New("dkg phase already done - call reshare")
	}
	d.state.Unlock()
	if !isLeader {
		// different logic for leader than the rest
		out, err := d.setupAutomaticDKG(c, in)
		return out, err
	}
	d.log.Info("init_dkg", "begin", "time", d.opts.clock.Now().Unix(), "leader", true)

	// setup the manager
	newSetup := func() (*setupManager, error) {
		return newDKGSetup(d.log, d.opts.clock, d.priv.Public, in.GetBeaconPeriod(), in.GetInfo())
	}

	// expect the group
	group, err := d.leaderRunSetup(newSetup)
	if err != nil {
		return nil, fmt.Errorf("drand: invalid setup configuration: %s", err)
	}

	// send it to everyone in the group nodes
	nodes := group.Nodes
	if err := d.pushDKGInfo([]*key.Node{}, nodes, 0, group, in.GetInfo().GetSecret(), in.GetInfo().GetTimeout()); err != nil {
		return nil, err
	}
	finalGroup, err := d.runDKG(true, group, in.GetInfo().GetTimeout(), in.GetEntropy())
	if err != nil {
		return nil, err
	}
	return finalGroup.ToProto(), nil
}

func (d *Drand) leaderRunSetup(newSetup func() (*setupManager, error)) (group *key.Group, err error) {
	// setup the manager
	d.state.Lock()
	if d.manager != nil {
		d.log.Info("reshare", "already_in_progress", "restart", "reshare")
		d.manager.StopPreemptively()
	}
	manager, err := newSetup()
	if err != nil {
		d.state.Unlock()
		return nil, fmt.Errorf("drand: invalid setup configuration: %s", err)
	}
	go manager.run()

	d.manager = manager
	d.state.Unlock()
	defer func() {
		// don't clear manager if pre-empted
		if err == errPreempted {
			return
		}
		d.state.Lock()
		// set back manager to nil afterwards to be able to run a new setup
		d.manager = nil
		d.state.Unlock()
	}()

	// wait to receive the keys & send them to the other nodes
	var ok bool
	select {
	case group, ok = <-manager.WaitGroup():
		if ok {
			var addr []string
			for _, k := range group.Nodes {
				addr = append(addr, k.Address())
			}
			d.log.Debug("init_dkg", "setup_phase", "keys_received", "["+strings.Join(addr, "-")+"]")
		} else {
			d.log.Debug("init_dkg", "pre-empted")
			return nil, errPreempted
		}
	case <-time.After(MaxWaitPrepareDKG):
		d.log.Debug("init_dkg", "time_out")
		manager.StopPreemptively()
		return nil, errors.New("time outs: no key received")
	}
	return group, nil
}

// runDKG setups the proper structures and protocol to run the DKG and waits
// until it finishes. If leader is true, this node sends the first packet.
func (d *Drand) runDKG(leader bool, group *key.Group, timeout uint32, randomness *drand.EntropyInfo) (*key.Group, error) {
	reader, user := extractEntropy(randomness)
	config := &dkg.Config{
		Suite:          key.KeyGroup.(dkg.Suite),
		NewNodes:       group.DKGNodes(),
		Longterm:       d.priv.Key,
		Reader:         reader,
		UserReaderOnly: user,
		FastSync:       true,
		Threshold:      group.Threshold,
		Nonce:          getNonce(group),
		Auth:           key.AuthScheme,
	}
	phaser := d.getPhaser(timeout)
	board := newBoard(d.log, d.privGateway.ProtocolClient, d.priv.Public, group)
	dkgProto, err := dkg.NewProtocol(config, board, phaser)
	if err != nil {
		return nil, err
	}

	d.state.Lock()
	d.dkgInfo = &dkgInfo{
		target: group,
		board:  board,
		phaser: phaser,
		conf:   config,
		proto:  dkgProto,
	}
	if leader {
		d.dkgInfo.started = true
	}
	d.state.Unlock()

	d.log.Info("init_dkg", "start_dkg")
	if leader {
		// phaser will kick off the first phase for every other nodes so
		// nodes will send their deals
		go phaser.Start()
	}
	finalGroup, err := d.WaitDKG()
	if err != nil {
		return nil, fmt.Errorf("drand: err during DKG: %v", err)
	}
	d.log.Info("init_dkg", "dkg_done", "starting_beacon_time", finalGroup.GenesisTime, "now", d.opts.clock.Now().Unix())
	// beacon will start at the genesis time specified
	go d.StartBeacon(false)
	return finalGroup, nil
}

// runResharing setups all necessary structures to run the resharing protocol
// and waits until it finishes (or timeouts). If leader is true, it sends the
// first packet so other nodes will start as soon as they receive it.
func (d *Drand) runResharing(leader bool, oldGroup, newGroup *key.Group, timeout uint32) (*key.Group, error) {
	oldNode := oldGroup.Find(d.priv.Public)
	oldPresent := oldNode != nil
	if leader && !oldPresent {
		d.log.Error("run_reshare", "invalid", "leader", leader, "old_present", oldPresent)
		return nil, errors.New("can not be a leader if not present in the old group")
	}
	newNode := newGroup.Find(d.priv.Public)
	newPresent := newNode != nil
	config := &dkg.Config{
		Suite:        key.KeyGroup.(dkg.Suite),
		NewNodes:     newGroup.DKGNodes(),
		OldNodes:     oldGroup.DKGNodes(),
		Longterm:     d.priv.Key,
		Threshold:    newGroup.Threshold,
		OldThreshold: oldGroup.Threshold,
		FastSync:     true,
		Nonce:        getNonce(newGroup),
		Auth:         key.AuthScheme,
	}
	err := func() error {
		d.state.Lock()
		defer d.state.Unlock()
		// gives the share to the dkg if we are a current node
		if oldPresent {
			if d.dkgInfo != nil {
				return errors.New("control: can't reshare from old node when DKG not finished first")
			}
			if d.share == nil {
				return errors.New("control: can't reshare without a share")
			}
			dkgShare := dkg.DistKeyShare(*d.share)
			config.Share = &dkgShare
		} else {
			// we are a new node, we want to make sure we reshare from the old
			// group public key
			config.PublicCoeffs = oldGroup.PublicKey.Coefficients
		}
		return nil
	}()
	if err != nil {
		return nil, err
	}
	board := newReshareBoard(d.log, d.privGateway.ProtocolClient, d.priv.Public, oldGroup, newGroup)
	phaser := d.getPhaser(timeout)

	dkgProto, err := dkg.NewProtocol(config, board, phaser)
	if err != nil {
		return nil, err
	}
	info := &dkgInfo{
		target: newGroup,
		board:  board,
		phaser: phaser,
		conf:   config,
		proto:  dkgProto,
	}
	d.state.Lock()
	d.dkgInfo = info
	if leader {
		d.log.Info("dkg_reshare", "leader_start", "target_group", hex.EncodeToString(newGroup.Hash()), "index", newNode.Index)
		d.dkgInfo.started = true
	}
	d.state.Unlock()

	if leader {
		// start the protocol so everyone else follows
		// it sends to all previous and new nodes. old nodes will start their
		// phaser so they will send the deals as soon as they receive this.
		go phaser.Start()
	}

	d.log.Info("init_dkg", "wait_dkg_end")
	finalGroup, err := d.WaitDKG()
	if err != nil {
		return nil, fmt.Errorf("drand: err during DKG: %v", err)
	}
	d.log.Info("dkg_reshare", "finished", "leader", leader)
	// runs the transition of the beacon
	go d.transition(oldGroup, oldPresent, newPresent)
	return finalGroup, nil
}

// This method sends the public key to the denoted leader address and then waits
// to receive the group file. After receiving it, it starts the DKG process in
// "waiting" mode, waiting for the leader to send the first packet.
func (d *Drand) setupAutomaticDKG(_ context.Context, in *drand.InitDKGPacket) (*drand.GroupPacket, error) {
	d.log.Info("init_dkg", "begin", "leader", false)
	// determine the leader's address
	laddr := in.GetInfo().GetLeaderAddress()
	lpeer := dnet.CreatePeer(laddr, in.GetInfo().GetLeaderTls())
	d.state.Lock()
	if d.receiver != nil {
		d.log.Info("dkg_setup", "already_in_progress", "restart", "dkg")
		d.receiver.stop()
	}
	receiver, err := newSetupReceiver(d.log, d.opts.clock, d.privGateway.ProtocolClient, in.GetInfo())
	if err != nil {
		d.log.Error("setup", "fail", "err", err)
		return nil, err
	}
	d.receiver = receiver
	d.state.Unlock()

	defer func() {
		d.state.Lock()
		d.receiver.stop()
		d.receiver = nil
		d.state.Unlock()
	}()
	// send public key to leader
	id := d.priv.Public.ToProto()
	prep := &drand.SignalDKGPacket{
		Node:        id,
		SecretProof: in.GetInfo().GetSecret(),
	}

	d.log.Debug("init_dkg", "send_key", "leader", lpeer.Address())
	err = d.privGateway.ProtocolClient.SignalDKGParticipant(context.Background(), lpeer, prep)
	if err != nil {
		return nil, fmt.Errorf("drand: err when receiving group: %s", err)
	}

	d.log.Debug("init_dkg", "wait_group")

	group, dkgTimeout, err := d.receiver.WaitDKGInfo()
	if err != nil {
		return nil, err
	}
	if group == nil {
		d.log.Debug("init_dkg", "wait_group", "canceled", "nil_group")
		return nil, errors.New("canceled operation")
	}

	now := d.opts.clock.Now().Unix()
	if group.GenesisTime < now {
		d.log.Error("genesis", "invalid", "given", group.GenesisTime)
		return nil, errors.New("control: group with genesis time in the past")
	}

	node := group.Find(d.priv.Public)
	if node == nil {
		d.log.Error("init_dkg", "absent_public_key_in_received_group")
		return nil, errors.New("drand: public key not found in group")
	}
	d.state.Lock()
	d.index = int(node.Index)
	d.state.Unlock()

	// run the dkg
	finalGroup, err := d.runDKG(false, group, dkgTimeout, in.GetEntropy())
	if err != nil {
		return nil, err
	}
	return finalGroup.ToProto(), nil
}

// similar to setupAutomaticDKG but with additional verification and information
// w.r.t. to the previous group
func (d *Drand) setupAutomaticResharing(c context.Context, oldGroup *key.Group, in *drand.InitResharePacket) (*drand.GroupPacket, error) {
	oldHash := oldGroup.Hash()
	// determine the leader's address
	laddr := in.GetInfo().GetLeaderAddress()
	lpeer := dnet.CreatePeer(laddr, in.GetInfo().GetLeaderTls())
	d.state.Lock()
	if d.receiver != nil {
		d.log.Info("reshare_setup", "already_in_progress", "restart", "reshare")
		d.receiver.stop()
	}

	receiver, err := newSetupReceiver(d.log, d.opts.clock, d.privGateway.ProtocolClient, in.GetInfo())
	if err != nil {
		d.log.Error("setup", "fail", "err", err)
		return nil, err
	}
	d.receiver = receiver
	d.state.Unlock()

	defer func() {
		d.state.Lock()
		d.receiver.stop()
		d.receiver = nil
		d.state.Unlock()
	}()
	// send public key to leader
	id := d.priv.Public.ToProto()
	prep := &drand.SignalDKGPacket{
		Node:              id,
		SecretProof:       in.GetInfo().GetSecret(),
		PreviousGroupHash: oldHash,
	}

	// we wait only a certain amount of time for the prepare phase
	nc, cancel := context.WithTimeout(c, MaxWaitPrepareDKG)
	defer cancel()

	d.log.Info("setup_reshare", "signaling_key_to_leader")
	err = d.privGateway.ProtocolClient.SignalDKGParticipant(nc, lpeer, prep)
	if err != nil {
		return nil, fmt.Errorf("drand: err when receiving group: %s", err)
	}

	newGroup, dkgTimeout, err := d.receiver.WaitDKGInfo()
	if err != nil {
		return nil, err
	}

	// some assertions that should be true but never too safe
	if oldGroup.GenesisTime != newGroup.GenesisTime {
		return nil, errors.New("control: old and new group have different genesis time")
	}

	if oldGroup.Period != newGroup.Period {
		return nil, errors.New("control: old and new group have different period - unsupported feature at the moment")
	}

	if !bytes.Equal(oldGroup.GetGenesisSeed(), newGroup.GetGenesisSeed()) {
		return nil, errors.New("control: old and new group have different genesis seed")
	}
	now := d.opts.clock.Now().Unix()
	if newGroup.TransitionTime < now {
		d.log.Error("setup_reshare", "invalid_transition", "given", newGroup.TransitionTime, "now", now)
		return nil, errors.New("control: new group with transition time in the past")
	}

	node := newGroup.Find(d.priv.Public)
	if node == nil {
		// It is ok to not have our key found in the new group since we may just
		// be a node that is leaving the network, but leaving gracefully, by
		// still participating in the resharing.
		d.log.Info("setup_reshare", "not_found_in_new_group")
	} else {
		d.state.Lock()
		d.index = int(node.Index)
		d.state.Unlock()
		d.log.Info("setup_reshare", "participate_newgroup", "index", node.Index)
	}

	// run the dkg !
	finalGroup, err := d.runResharing(false, oldGroup, newGroup, dkgTimeout)
	if err != nil {
		return nil, err
	}
	return finalGroup.ToProto(), nil
}

func (d *Drand) extractGroup(old *drand.GroupInfo) (oldGroup *key.Group, err error) {
	d.state.Lock()
	if oldGroup, err = extractGroup(old); err != nil {
		// try to get the current group
		if d.group == nil {
			d.state.Unlock()
			return nil, errors.New("drand: can't init-reshare if no old group provided")
		}
		d.log.With("module", "control").Debug("init_reshare", "using_stored_group")
		oldGroup = d.group
		err = nil
	}
	d.state.Unlock()
	return
}

// InitReshare receives information about the old and new group from which to
// operate the resharing protocol.
func (d *Drand) InitReshare(c context.Context, in *drand.InitResharePacket) (*drand.GroupPacket, error) {
	oldGroup, err := d.extractGroup(in.Old)
	if err != nil {
		return nil, err
	}

	if !in.GetInfo().GetLeader() {
		d.log.Info("init_reshare", "begin", "leader", false)
		return d.setupAutomaticResharing(c, oldGroup, in)
	}
	d.log.Info("init_reshare", "begin", "leader", true, "time", d.opts.clock.Now())

	newSetup := func() (*setupManager, error) {
		return newReshareSetup(d.log, d.opts.clock, d.priv.Public, oldGroup, in)
	}

	newGroup, err := d.leaderRunSetup(newSetup)
	if err != nil {
		return nil, fmt.Errorf("drand: invalid setup configuration: %s", err)
	}
	// some assertions that should always be true but never too safe
	if oldGroup.GenesisTime != newGroup.GenesisTime {
		return nil, errors.New("control: old and new group have different genesis time")
	}
	if oldGroup.GenesisTime > d.opts.clock.Now().Unix() {
		return nil, errors.New("control: genesis time is in the future")
	}
	if oldGroup.Period != newGroup.Period {
		return nil, errors.New("control: old and new group have different period - unsupported feature at the moment")
	}
	if newGroup.TransitionTime < d.opts.clock.Now().Unix() {
		return nil, errors.New("control: group with transition time in the past")
	}
	if !bytes.Equal(newGroup.GetGenesisSeed(), oldGroup.GetGenesisSeed()) {
		return nil, errors.New("control: old and new group have different genesis seed")
	}

	// send it to everyone in the group nodes
	if err := d.pushDKGInfo(oldGroup.Nodes, newGroup.Nodes,
		oldGroup.Threshold,
		newGroup,
		in.GetInfo().GetSecret(),
		in.GetInfo().GetTimeout()); err != nil {
		d.log.Error("push_group", err)
		return nil, errors.New("fail to push new group")
	}

	finalGroup, err := d.runResharing(true, oldGroup, newGroup, in.GetInfo().GetTimeout())
	if err != nil {
		return nil, err
	}
	return finalGroup.ToProto(), nil
}

// PingPong simply responds with an empty packet, proving that this drand node
// is up and alive.
func (d *Drand) PingPong(c context.Context, in *drand.Ping) (*drand.Pong, error) {
	return &drand.Pong{}, nil
}

// Share is a functionality of Control Service defined in protobuf/control that requests the private share of the drand node running locally
func (d *Drand) Share(ctx context.Context, in *drand.ShareRequest) (*drand.ShareResponse, error) {
	share, err := d.store.LoadShare()
	if err != nil {
		return nil, err
	}
	id := uint32(share.Share.I)
	buff, err := share.Share.V.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &drand.ShareResponse{Index: id, Share: buff}, nil
}

// PublicKey is a functionality of Control Service defined in protobuf/control
// that requests the long term public key of the drand node running locally
func (d *Drand) PublicKey(ctx context.Context, in *drand.PublicKeyRequest) (*drand.PublicKeyResponse, error) {
	d.state.Lock()
	defer d.state.Unlock()
	keyPair, err := d.store.LoadKeyPair()
	if err != nil {
		return nil, err
	}
	protoKey, err := keyPair.Public.Key.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &drand.PublicKeyResponse{PubKey: protoKey}, nil
}

// PrivateKey is a functionality of Control Service defined in protobuf/control
// that requests the long term private key of the drand node running locally
func (d *Drand) PrivateKey(ctx context.Context, in *drand.PrivateKeyRequest) (*drand.PrivateKeyResponse, error) {
	d.state.Lock()
	defer d.state.Unlock()
	keyPair, err := d.store.LoadKeyPair()
	if err != nil {
		return nil, err
	}
	protoKey, err := keyPair.Key.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &drand.PrivateKeyResponse{PriKey: protoKey}, nil
}

// GroupFile replies with the distributed key in the response
func (d *Drand) GroupFile(ctx context.Context, in *drand.GroupRequest) (*drand.GroupPacket, error) {
	d.state.Lock()
	defer d.state.Unlock()
	if d.group == nil {
		return nil, errors.New("drand: no dkg group setup yet")
	}
	protoGroup := d.group.ToProto()
	return protoGroup, nil
}

// Shutdown stops the node
func (d *Drand) Shutdown(ctx context.Context, in *drand.ShutdownRequest) (*drand.ShutdownResponse, error) {
	d.Stop(ctx)
	return nil, nil
}

func extractGroup(i *drand.GroupInfo) (*key.Group, error) {
	var g = new(key.Group)
	switch x := i.Location.(type) {
	case *drand.GroupInfo_Path:
		// search group file via local filesystem path
		if err := key.Load(x.Path, g); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("control: can't allow new empty group")
	}
	if g.Threshold < vss.MinimumT(g.Len()) {
		return nil, errors.New("control: threshold of new group too low ")
	}
	return g, nil
}

func extractEntropy(i *drand.EntropyInfo) (io.Reader, bool) {
	if i == nil {
		return nil, false
	}
	r := entropy.NewScriptReader(i.Script)
	user := i.UserOnly
	return r, user
}

func (d *Drand) getPhaser(timeout uint32) *dkg.TimePhaser {
	tDuration := time.Duration(timeout) * time.Second
	if timeout == 0 {
		tDuration = DefaultDKGTimeout
	}
	return dkg.NewTimePhaserFunc(func(phase dkg.Phase) {
		d.opts.clock.Sleep(tDuration)
		d.log.Debug("phaser_finished", phase)
	})
}

func nodesContainAddr(nodes []*key.Node, addr string) bool {
	for _, n := range nodes {
		if n.Address() == addr {
			return true
		}
	}
	return false
}

// nodeUnion takes the union of two sets of nodes
func nodeUnion(a, b []*key.Node) []*key.Node {
	out := append([]*key.Node{}, a...)
	for _, n := range b {
		if !nodesContainAddr(a, n.Address()) {
			out = append(out, n)
		}
	}
	return out
}

type pushResult struct {
	address string
	err     error
}

// pushDKGInfoPacket sets a specific DKG info packet to spcified nodes, and returns a stream of responses.
func (d *Drand) pushDKGInfoPacket(ctx context.Context, nodes []*key.Node, packet *drand.DKGInfoPacket) chan pushResult {
	results := make(chan pushResult, len(nodes))

	for _, node := range nodes {
		if node.Address() == d.priv.Public.Address() {
			continue
		}
		go func(i *key.Identity) {
			err := d.privGateway.ProtocolClient.PushDKGInfo(ctx, i, packet)
			results <- pushResult{i.Address(), err}
		}(node.Identity)
	}

	return results
}

// pushDKGInfo sends the information to run the DKG to all specified nodes. The
// call is blocking until all nodes have replied or after one minute timeouts.
func (d *Drand) pushDKGInfo(outgoing, incoming []*key.Node, previousThreshold int, group *key.Group, secret []byte, timeout uint32) error {
	// sign the group to prove you are the leader
	signature, err := key.AuthScheme.Sign(d.priv.Key, group.Hash())
	if err != nil {
		d.log.Error("setup", "leader", "group_signature", err)
		return fmt.Errorf("drand: error signing group: %w", err)
	}
	packet := &drand.DKGInfoPacket{
		NewGroup:    group.ToProto(),
		SecretProof: secret,
		DkgTimeout:  timeout,
		Signature:   signature,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	newThreshold := group.Threshold
	if nodesContainAddr(outgoing, d.priv.Public.Address()) {
		previousThreshold--
	}
	if nodesContainAddr(incoming, d.priv.Public.Address()) {
		newThreshold--
	}
	to := nodeUnion(outgoing, incoming)

	results := d.pushDKGInfoPacket(ctx, to, packet)

	total := len(to) - 1
	for total > 0 {
		select {
		case ok := <-results:
			total--
			if ok.err != nil {
				d.log.Error("push_dkg", "failed to push", "to", ok.address, "err", ok.err)
				continue
			}
			d.log.Debug("push_dkg", "sending_group", "success_to", ok.address, "left", total)
			if nodesContainAddr(outgoing, ok.address) {
				previousThreshold--
			}
			if nodesContainAddr(incoming, ok.address) {
				newThreshold--
			}
		case <-d.opts.clock.After(time.Minute):
			if previousThreshold <= 0 && newThreshold <= 0 {
				d.log.Info("push_dkg", "sending_group", "status", "enough succeeded", "missed", total)
				return nil
			}
			d.log.Warn("push_dkg", "sending_group", "status", "timeout")
			return errors.New("push group timeout")
		}
	}
	if previousThreshold > 0 || newThreshold > 0 {
		d.log.Info("push_dkg", "sending_group", "status", "not enough succeeded", "prev", previousThreshold, "new", newThreshold)
		return errors.New("push group failure")
	}
	d.log.Info("push_dkg", "sending_group", "status", "all succeeded")
	return nil
}

func getNonce(g *key.Group) []byte {
	h := sha256.New()
	if g.TransitionTime != 0 {
		_ = binary.Write(h, binary.BigEndian, g.TransitionTime)
	} else {
		_ = binary.Write(h, binary.BigEndian, g.GenesisTime)
	}
	return h.Sum(nil)
}
