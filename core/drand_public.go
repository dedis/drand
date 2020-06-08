package core

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/drand/drand/chain"
	"github.com/drand/drand/entropy"
	"github.com/drand/drand/key"
	"github.com/drand/drand/net"
	"github.com/drand/drand/protobuf/drand"
	"github.com/drand/kyber/encrypt/ecies"
)

// Setup is the public method to call during a DKG protocol.
func (d *Drand) FreshDKG(c context.Context, in *drand.DKGPacket) (*drand.Empty, error) {
	d.state.Lock()
	defer d.state.Unlock()
	if d.dkgDone {
		return nil, errors.New("drand: dkg finished already")
	}
	if d.dkgInfo == nil {
		return nil, errors.New("drand: no dkg running")
	}
	addr := net.RemoteAddress(c)
	if !d.dkgInfo.started {
		d.log.Info("init_dkg", "start", "signal_leader", addr, "group", hex.EncodeToString(d.dkgInfo.target.Hash()))
		d.dkgInfo.started = true
		go d.dkgInfo.phaser.Start()
	}
	d.dkgInfo.board.FreshDKG(c, in)
	return new(drand.Empty), nil
}

// Reshare is called when a resharing protocol is in progress
func (d *Drand) ReshareDKG(c context.Context, in *drand.ResharePacket) (*drand.Empty, error) {
	d.state.Lock()
	defer d.state.Unlock()

	if d.dkgInfo == nil {
		return nil, errors.New("drand: no dkg setup yet")
	}
	addr := net.RemoteAddress(c)
	if !d.dkgInfo.started {
		d.dkgInfo.started = true
		d.log.Info("init_reshare", "start", "signal_leader", addr, "group", hex.EncodeToString(d.dkgInfo.target.Hash()), "target_index", d.dkgInfo.target.Find(d.priv.Public).Index)
		go d.dkgInfo.phaser.Start()
	}

	d.dkgInfo.board.ReshareDKG(c, in)
	return new(drand.Empty), nil
}

// NewBeacon methods receives a beacon generation requests and answers
// with the partial signature from this drand node.
func (d *Drand) PartialBeacon(c context.Context, in *drand.PartialBeaconPacket) (*drand.Empty, error) {
	d.state.Lock()
	defer d.state.Unlock()
	if d.beacon == nil {
		return nil, errors.New("drand: beacon not setup yet")
	}
	return d.beacon.ProcessPartialBeacon(c, in)
}

// PublicRand returns a public random beacon according to the request. If the Round
// field is 0, then it returns the last one generated.
func (d *Drand) PublicRand(c context.Context, in *drand.PublicRandRequest) (*drand.PublicRandResponse, error) {
	var addr = net.RemoteAddress(c)
	d.state.Lock()
	defer d.state.Unlock()
	if d.beacon == nil {
		return nil, errors.New("drand: beacon generation not started yet")
	}
	var r *chain.Beacon
	var err error
	if in.GetRound() == 0 {
		r, err = d.beacon.Store().Last()
	} else {
		// fetch the correct entry or the next one if not found
		r, err = d.beacon.Store().Get(in.GetRound())
	}
	if err != nil || r == nil {
		d.log.Debug("public_rand", "unstored_beacon", "round", in.GetRound(), "from", addr)
		return nil, fmt.Errorf("can't retrieve beacon: %s %s", err, r)
	}
	d.log.Info("public_rand", addr, "round", r.Round, "reply", r.String())
	return beaconToProto(r), nil
}

func (d *Drand) PublicRandStream(req *drand.PublicRandRequest, stream drand.Public_PublicRandStreamServer) error {
	var b *chain.Handler
	d.state.Lock()
	if d.beacon == nil {
		return errors.New("beacon has not started on this node yet")
	}
	b = d.beacon
	d.state.Unlock()
	lastb, err := b.Store().Last()
	if err != nil {
		return err
	}
	addr := net.RemoteAddress(stream.Context())
	done := make(chan error, 1)
	d.log.Debug("request", "stream", "from", addr, "round", req.GetRound())
	if req.GetRound() != 0 && req.GetRound() <= lastb.Round {
		// we need to stream from store first
		var err error
		b.Store().Cursor(func(c chain.Cursor) {
			for bb := c.Seek(req.GetRound()); bb != nil; bb = c.Next() {
				if err = stream.Send(beaconToProto(bb)); err != nil {
					d.log.Debug("stream", err)
					return
				}
			}
		})
		if err != nil {
			return err
		}
	}
	// then we can stream from any new rounds
	// register a callback for the duration of this stream
	d.callbacks.AddCallback(addr, func(b *chain.Beacon) {
		err := stream.Send(&drand.PublicRandResponse{
			Round:             b.Round,
			Signature:         b.Signature,
			PreviousSignature: b.PreviousSig,
			Randomness:        b.Randomness(),
		})
		// if connection has a problem, we drop the callback
		if err != nil {
			d.callbacks.DelCallback(addr)
			done <- err
		}
	})
	return <-done
}

// PrivateRand returns an ECIES encrypted random blob of 32 bytes from /dev/urandom
func (d *Drand) PrivateRand(c context.Context, priv *drand.PrivateRandRequest) (*drand.PrivateRandResponse, error) {
	if !d.opts.enablePrivate {
		return nil, errors.New("private randomness is disabled")
	}
	msg, err := ecies.Decrypt(key.KeyGroup, d.priv.Key, priv.GetRequest(), EciesHash)
	if err != nil {
		d.log.With("module", "public").Error("private", "invalid ECIES", "err", err.Error())
		return nil, errors.New("invalid ECIES request")
	}

	clientKey := key.KeyGroup.Point()
	if err := clientKey.UnmarshalBinary(msg); err != nil {
		return nil, errors.New("invalid client key")
	}
	randomness, err := entropy.GetRandom(nil, 32)
	if err != nil {
		return nil, fmt.Errorf("error gathering randomness: %s", err)
	} else if len(randomness) != 32 {
		return nil, fmt.Errorf("error gathering randomness: expected 32 bytes, got %d", len(randomness))
	}

	obj, err := ecies.Encrypt(key.KeyGroup, clientKey, randomness[:], EciesHash)
	return &drand.PrivateRandResponse{Response: obj}, err
}

// Home ...
func (d *Drand) Home(c context.Context, in *drand.HomeRequest) (*drand.HomeResponse, error) {
	d.log.With("module", "public").Info("home", net.RemoteAddress(c))
	return &drand.HomeResponse{
		Status: fmt.Sprintf("drand up and running on %s",
			d.priv.Public.Address()),
	}, nil
}

// ChainInfo replies with the chain information this node participates to
func (d *Drand) ChainInfo(ctx context.Context, in *drand.ChainInfoRequest) (*drand.ChainInfoPacket, error) {
	d.state.Lock()
	defer d.state.Unlock()
	if d.group == nil {
		return nil, errors.New("drand: no dkg group setup yet")
	}
	return chain.NewChainInfo(d.group).ToProto(), nil
}

func (d *Drand) SignalDKGParticipant(ctx context.Context, p *drand.SignalDKGPacket) (*drand.Empty, error) {
	d.state.Lock()
	defer d.state.Unlock()
	if d.manager == nil {
		return nil, errors.New("no manager")
	}
	addr := net.RemoteAddress(ctx)
	// manager will verify if information are correct
	err := d.manager.ReceivedKey(addr, p)
	if err != nil {
		return nil, err
	}
	return new(drand.Empty), nil
}

func (d *Drand) PushDKGInfo(ctx context.Context, in *drand.DKGInfoPacket) (*drand.Empty, error) {
	d.state.Lock()
	defer d.state.Unlock()
	if d.receiver == nil {
		return nil, errors.New("no receiver setup")
	}
	d.log.Info("push_group", "received_new")
	// the control routine will receive this info and start the dkg at the right
	// time - if that is the right secret.
	err := d.receiver.PushDKGInfo(in)
	if err != nil {
		return nil, err
	}
	return new(drand.Empty), nil
}

// SyncChain is a inter-node protocol that replies to a syncing request from a
// given round
func (d *Drand) SyncChain(req *drand.SyncRequest, stream drand.Protocol_SyncChainServer) error {
	d.state.Lock()
	beacon := d.beacon
	d.state.Unlock()
	if beacon != nil {
		beacon.SyncChain(req, stream)
	}
	return nil
}
