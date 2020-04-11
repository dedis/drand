package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/drand/drand/dkg"
	"github.com/drand/drand/entropy"
	"github.com/drand/drand/key"
	dnet "github.com/drand/drand/net"
	"github.com/drand/drand/protobuf/drand"
	control "github.com/drand/drand/protobuf/drand"
	vss "github.com/drand/kyber/share/vss/pedersen"
)

var syncTime = 500 * time.Millisecond

// InitDKG take a InitDKGPacket, extracts the informations needed and wait for the
// DKG protocol to finish. If the request specifies this node is a leader, it
// starts the DKG protocol.
func (d *Drand) InitDKG(c context.Context, in *control.InitDKGPacket) (*control.Empty, error) {
	d.log.Info("init_dkg", "begin")
	isLeader := in.GetInfo().GetLeader()
	if !isLeader {
		// different logic for leader than the rest
		return d.setupAutomaticDKG(c, in)
	}

	// setup the manager
	d.state.Lock()
	if d.manager != nil {
		return nil, errors.New("drand: setup dkg already in progress")
	}
	if d.dkgDone == true {
		d.state.Unlock()
		return nil, errors.New("drand: dkg phase already done. Can't run 2 init DKG")
	}

	manager, err := newSetupManager(d.log, d.opts.clock, d.priv.Public, in)
	if err != nil {
		return nil, fmt.Errorf("drand: invalid setup configuration: %s", err)
	}
	d.manager = manager
	d.state.Unlock()
	defer func() {
		d.state.Lock()
		// set back manager to nil afterwards to be able to run a new setup
		d.manager = nil
		d.state.Unlock()
	}()

	// wait to receive the keys & send them to the other nodes
	var group *key.Group
	select {
	case group = <-d.manager.WaitForGroupSent():
		var addr []string
		for _, k := range group.Identities() {
			addr = append(addr, k.Address())
		}
		d.log.Debug("init_dkg", "setup_phase", "keys_received", "["+strings.Join(addr, "-")+"]")
	case <-time.After(MaxWaitPrepareDKG):
		d.log.Debug("init_dkg", "time_out")
		manager.StopPreemptively()
		return nil, errors.New("time outs: no key received")
	}

	err = d.runDKG(true, group, in.GetInfo().GetTimeout(), in.GetEntropy())
	if err != nil {
		return nil, err
	}
	return &control.Empty{}, nil
}

func (d *Drand) runDKG(leader bool, group *key.Group, timeout string, entropy *control.EntropyInfo) error {
	reader, user := extractEntropy(entropy)
	dkgConfig := &dkg.Config{
		Suite:          key.KeyGroup.(dkg.Suite),
		NewNodes:       group,
		Key:            d.priv,
		Reader:         reader,
		UserReaderOnly: user,
		Clock:          d.opts.clock,
	}
	if err := setTimeout(dkgConfig, timeout); err != nil {
		return fmt.Errorf("drand: invalid timeout: %s", err)
	}
	d.state.Lock()
	d.nextConf = dkgConfig
	d.state.Unlock()

	if leader {
		d.log.Info("init_dkg", "start_dkg_leader")
		if err := d.StartDKG(dkgConfig); err != nil {
			return err
		}
	}

	d.log.Info("init_dkg", "wait_dkg_end")
	if err := d.WaitDKG(dkgConfig); err != nil {
		return fmt.Errorf("drand: err during DKG: %v", err)
	}
	// beacon will start at the genesis time specified
	d.StartBeacon(false)
	return nil
}

// This method sends the public key to the denoted leader address and then waits
// to receive the group file. After receiving it, it starts the DKG process in
// "waiting" mode, waiting for the leader to send the first packet.
func (d *Drand) setupAutomaticDKG(c context.Context, in *control.InitDKGPacket) (*control.Empty, error) {
	n, thr, dkgTimeout, err := validInitPacket(in)
	if err != nil {
		return nil, err
	}
	// determine the leader's address
	laddr := in.GetInfo().GetLeaderAddress()
	if _, _, err := net.SplitHostPort(laddr); err != nil {
		return nil, fmt.Errorf("invalid leader address: %s", err)
	}
	lpeer := dnet.CreatePeer(laddr, in.GetInfo().GetLeaderTls())
	d.state.Lock()
	if d.waitAutomatic {
		d.state.Unlock()
		return nil, errors.New("drand: already waiting for an automatic setup")
	}
	d.waitAutomatic = true
	d.state.Unlock()

	// send public key to leader
	key, _ := d.priv.Public.Key.MarshalBinary()
	id := &drand.Identity{
		Address: d.priv.Public.Address(),
		Key:     key,
		Tls:     d.priv.Public.IsTLS(),
	}
	prep := &drand.PrepareDKGPacket{
		Node:        id,
		Expected:    uint32(n),
		Threshold:   uint32(thr),
		DkgTimeout:  uint64(dkgTimeout.Seconds()),
		SecretProof: in.GetInfo().GetSecret(),
	}

	// we wait only a certain amount of time for the prepare phase
	nc, cancel := context.WithTimeout(c, MaxWaitPrepareDKG)
	defer cancel()

	// expect group
	groupPacket, err := d.gateway.ProtocolClient.PrepareDKGGroup(nc, lpeer, prep)
	if err != nil {
		return nil, fmt.Errorf("drand: err when receiving group: %s", err)
	}

	// verify things are all in order
	group, err := protoToGroup(groupPacket)
	if err != nil {
		return nil, fmt.Errorf("group from leader invalid: %s", err)
	}
	if group.GenesisTime < d.opts.clock.Now().Unix() {
		d.log.Error("genesis", "invalid", "given", group.GenesisTime, "now", d.opts.clock.Now().Unix())
		return nil, errors.New("control: group with genesis time in the past")
	}
	index, found := group.Index(d.priv.Public)
	if !found {
		return nil, errors.New("drand: public key not found in group")
	}
	d.state.Lock()
	d.index = index
	d.state.Unlock()

	// run the dkg !
	err = d.runDKG(false, group, in.GetInfo().GetTimeout(), in.GetEntropy())
	if err != nil {
		return nil, err
	}
	return new(drand.Empty), nil
}

// InitReshare receives information about the old and new group from which to
// operate the resharing protocol. It starts the resharing protocol if the
// received node is stated as a leader and is present in the old group.
// This function waits for the resharing DKG protocol to finish.
func (d *Drand) InitReshare(c context.Context, in *control.InitResharePacket) (*control.Empty, error) {
	var oldGroup, newGroup *key.Group
	var err error

	if newGroup, err = extractGroup(in.New); err != nil {
		return nil, err
	}

	d.state.Lock()
	if oldGroup, err = extractGroup(in.Old); err != nil {
		// try to get the current group
		if d.group == nil {
			d.state.Unlock()
			return nil, errors.New("drand: can't init-reshare if no old group provided")
		}
		d.log.With("module", "control").Debug("init_reshare", "old group equal current group")
		oldGroup = d.group
	}
	d.state.Unlock()

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

	oldIdx, oldPresent := oldGroup.Index(d.priv.Public)
	newIdx, newPresent := newGroup.Index(d.priv.Public)
	var dkgConf *dkg.Config
	err = func() error {
		d.state.Lock()
		defer d.state.Unlock()

		if oldPresent {
			if d.group == nil {
				return errors.New("control: present in old group but no dkg here")
			}
			// stateful verification checking if we are in the old group that we have
			// and the one we receive
			currHash, err := d.group.Hash()
			if err != nil {
				return err
			}
			oldHash, err := oldGroup.Hash()
			if err != nil {
				return err
			}
			if currHash != oldHash {
				return errors.New("control: given old group is not the same as one saved")
			}
		}

		// prepare dkg config to run the protocol
		dkgConf = &dkg.Config{
			OldNodes: oldGroup,
			NewNodes: newGroup,
			Key:      d.priv,
			Suite:    key.KeyGroup.(dkg.Suite),
			Clock:    d.opts.clock,
		}

		// gives the share to the dkg if we are a current node
		if oldPresent {
			if !d.dkgDone {
				return errors.New("control: can't reshare from old node when DKG not finished first")
			}
			if d.share == nil {
				return errors.New("control: can't reshare without a share")
			}
			dkgConf.Share = d.share
		}

		nextHash, err := newGroup.Hash()
		if err != nil {
			return err
		}

		if err := setTimeout(dkgConf, in.Timeout); err != nil {
			return fmt.Errorf("drand: invalid timeout: %s", err)
		}

		d.nextGroupHash = nextHash
		d.nextGroup = newGroup
		d.nextConf = dkgConf
		d.nextOldPresent = oldPresent
		return nil
	}()

	if err != nil {
		return nil, err
	}

	if oldPresent && in.GetIsLeader() {
		// only the root sends a pre-message to the other old
		// nodes and start the DKG
		d.startResharingAsLeader(dkgConf, oldIdx)
	}
	if err := d.WaitDKG(dkgConf); err != nil {
		return nil, err
	}
	d.log.Info("dkg_reshare", "finished")
	fmt.Printf("going to TRANSITION: oidx %d nidx %d oldpresent %v newpresent %v\n", oldIdx, newIdx, oldPresent, newPresent)
	go d.transition(oldGroup, oldPresent, newPresent)
	return &control.Empty{}, nil
}

func (d *Drand) startResharingAsLeader(dkgConf *dkg.Config, oidx int) {
	d.log.With("module", "control").Debug("leader_reshare", "start signalling")
	d.state.Lock()
	msg := &control.ResharePacket{GroupHash: d.nextGroupHash}
	// send resharing packet to signal start of the protocol to other old
	// nodes
	for i, p := range d.nextConf.OldNodes.Identities() {
		if i == oidx {
			continue
		}
		id := p
		// XXX find way to just have a small RPC timeout if one is down.
		//fmt.Printf("drand leader %s -> signal to %s\n", d.priv.Public.Addr, id.Addr)
		if _, err := d.gateway.ProtocolClient.ReshareDKG(context.TODO(), id, msg); err != nil {
			//if _, err := d.gateway.InternalClient.Reshare(id, msg, grpc.FailFast(true)); err != nil {
			d.log.With("module", "control").Error("leader_reshare", err)
		}
	}
	d.state.Unlock()
	d.log.With("module", "control").Debug("leader_reshare", "start DKG")
	d.StartDKG(dkgConf)
}

func (d *Drand) SyncChain(req *drand.SyncRequest, stream drand.Protocol_SyncChainServer) error {
	d.state.Lock()
	beacon := d.beacon
	d.state.Unlock()
	if beacon != nil {
		beacon.SyncChain(req, stream)
	}
	return nil
}

// DistKey returns the distributed key corresponding to the current group
func (d *Drand) DistKey(context.Context, *drand.DistKeyRequest) (*drand.DistKeyResponse, error) {
	pt, err := d.store.LoadDistPublic()
	if err != nil {
		return nil, errors.New("drand: could not load dist. key")
	}
	buff, err := pt.Key().MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &drand.DistKeyResponse{
		Key: buff,
	}, nil
}

// PingPong simply responds with an empty packet, proving that this drand node
// is up and alive.
func (d *Drand) PingPong(c context.Context, in *control.Ping) (*control.Pong, error) {
	return &control.Pong{}, nil
}

// Share is a functionality of Control Service defined in protobuf/control that requests the private share of the drand node running locally
func (d *Drand) Share(ctx context.Context, in *control.ShareRequest) (*control.ShareResponse, error) {
	share, err := d.store.LoadShare()
	if err != nil {
		return nil, err
	}
	id := uint32(share.Share.I)
	buff, err := share.Share.V.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &control.ShareResponse{Index: id, Share: buff}, nil
}

// PublicKey is a functionality of Control Service defined in protobuf/control that requests the long term public key of the drand node running locally
func (d *Drand) PublicKey(ctx context.Context, in *control.PublicKeyRequest) (*control.PublicKeyResponse, error) {
	d.state.Lock()
	defer d.state.Unlock()
	key, err := d.store.LoadKeyPair()
	if err != nil {
		return nil, err
	}
	protoKey, err := key.Public.Key.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &control.PublicKeyResponse{PubKey: protoKey}, nil
}

// PrivateKey is a functionality of Control Service defined in protobuf/control that requests the long term private key of the drand node running locally
func (d *Drand) PrivateKey(ctx context.Context, in *control.PrivateKeyRequest) (*control.PrivateKeyResponse, error) {
	d.state.Lock()
	defer d.state.Unlock()
	key, err := d.store.LoadKeyPair()
	if err != nil {
		return nil, err
	}
	protoKey, err := key.Key.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &control.PrivateKeyResponse{PriKey: protoKey}, nil
}

// CollectiveKey replies with the distributed key in the response
func (d *Drand) CollectiveKey(ctx context.Context, in *control.CokeyRequest) (*control.CokeyResponse, error) {
	d.state.Lock()
	defer d.state.Unlock()

	key, err := d.store.LoadDistPublic()
	if err != nil {
		return nil, err
	}
	protoKey, err := key.Key().MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &control.CokeyResponse{CoKey: protoKey}, nil
}

// GroupFile replies with the distributed key in the response
func (d *Drand) GroupFile(ctx context.Context, in *control.GroupTOMLRequest) (*control.GroupTOMLResponse, error) {
	d.state.Lock()
	defer d.state.Unlock()
	if d.group == nil {
		return nil, errors.New("drand: no dkg group setup yet")
	}
	gtoml := d.group.TOML()
	var buff bytes.Buffer
	err := toml.NewEncoder(&buff).Encode(gtoml)
	if err != nil {
		return nil, fmt.Errorf("drand: error encoding group to TOML: %s", err)
	}
	groupStr := buff.String()
	return &drand.GroupTOMLResponse{GroupToml: groupStr}, nil
}

func (d *Drand) Shutdown(ctx context.Context, in *control.ShutdownRequest) (*control.ShutdownResponse, error) {
	d.Stop()
	return nil, nil
}

func extractGroup(i *control.GroupInfo) (*key.Group, error) {
	var g = &key.Group{}
	switch x := i.Location.(type) {
	case *control.GroupInfo_Path:
		// search group file via local filesystem path
		if err := key.Load(x.Path, g); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("control: can't allow new empty group")
	}
	// run a few checks on the proposed group
	if g.Len() < 4 {
		return nil, errors.New("control: can't accept group with fewer than 4 members")
	}
	if g.Threshold < vss.MinimumT(g.Len()) {
		return nil, errors.New("control: threshold of new group too low ")
	}

	return g, nil
}

func extractEntropy(i *control.EntropyInfo) (io.Reader, bool) {
	if i == nil {
		return nil, false
	}
	r := entropy.NewScriptReader(i.Script)
	user := i.UserOnly
	return r, user
}

func setTimeout(c *dkg.Config, timeoutStr string) error {
	// try parsing the timeout
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		if timeoutStr != "" {
			return fmt.Errorf("invalid timeout: %s", err)
		}
		timeout, _ = time.ParseDuration(DefaultDKGTimeout)
	}
	c.Timeout = timeout
	return nil
}
