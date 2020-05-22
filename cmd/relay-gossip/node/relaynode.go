package node

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/drand/drand/cmd/relay-gossip/lp2p"
	"github.com/drand/drand/protobuf/drand"
	"github.com/gogo/protobuf/proto"
	bds "github.com/ipfs/go-ds-badger2"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/prometheus/common/log"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// GossipRelayConfig configures a gossip relay node.
type GossipRelayConfig struct {
	Network         string
	PeerWith        []string
	Addr            string
	DataDir         string
	IdentityPath    string
	CertPath        string
	Insecure        bool
	DrandPublicGRPC string
}

// GossipRelayNode is a gossip relay runtime.
type GossipRelayNode struct {
	bootstrap []ma.Multiaddr
	ds        *bds.Datastore
	priv      crypto.PrivKey
	h         host.Host
	ps        *pubsub.PubSub
	t         *pubsub.Topic
	opts      []grpc.DialOption
	addr      []ma.Multiaddr
	done      chan struct{}
}

// NewGossipRelayNode starts a new gossip relay node.
func NewGossipRelayNode(cfg *GossipRelayConfig) (*GossipRelayNode, error) {
	bootstrap, err := ParseMultiaddrSlice(cfg.PeerWith)
	if err != nil {
		return nil, xerrors.Errorf("parsing peer-with: %w", err)
	}

	ds, err := bds.NewDatastore(cfg.DataDir, nil)
	if err != nil {
		return nil, xerrors.Errorf("opening datastore: %w", err)
	}

	priv, err := lp2p.LoadOrCreatePrivKey(cfg.IdentityPath)
	if err != nil {
		return nil, xerrors.Errorf("loading p2p key: %w", err)
	}

	h, ps, err := lp2p.ConstructHost(ds, priv, cfg.Addr, bootstrap)
	if err != nil {
		return nil, xerrors.Errorf("constructing host: %w", err)
	}

	addr, err := h.Network().InterfaceListenAddresses()
	if err != nil {
		return nil, xerrors.Errorf("getting InterfaceListenAddresses: %w", err)
	}

	for _, a := range addr {
		log.Infof("%s/p2p/%s\n", a, h.ID())
	}

	t, err := ps.Join(lp2p.PubSubTopic(cfg.Network))
	if err != nil {
		return nil, xerrors.Errorf("joining topic: %w", err)
	}

	opts := []grpc.DialOption{}
	if cfg.CertPath != "" {
		creds, err := credentials.NewClientTLSFromFile(cfg.CertPath, "")
		if err != nil {
			return nil, xerrors.Errorf("loading cert file: %w", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else if cfg.Insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}
	g := &GossipRelayNode{
		bootstrap: bootstrap,
		ds:        ds,
		priv:      priv,
		h:         h,
		ps:        ps,
		t:         t,
		opts:      opts,
		addr:      addr,
		done:      make(chan struct{}),
	}
	go g.start(cfg.DrandPublicGRPC)
	return g, nil
}

func (g *GossipRelayNode) Addr() []ma.Multiaddr {
	return g.addr
}

// Shutdown stops the relay node.
func (g *GossipRelayNode) Shutdown() {
	close(g.done)
}

func ParseMultiaddrSlice(peer []string) ([]ma.Multiaddr, error) {
	out := make([]ma.Multiaddr, len(peer))
	for i, peer := range peer {
		m, err := ma.NewMultiaddr(peer)
		if err != nil {
			return nil, xerrors.Errorf("parsing multiaddr\"%s\": %w", peer, err)
		}
		out[i] = m
	}
	return out, nil
}

func (g *GossipRelayNode) start(drandPublicGRPC string) {
	for {
		select {
		case <-g.done:
			return
		default:
		}
		conn, err := grpc.Dial(drandPublicGRPC, g.opts...)
		if err != nil {
			log.Warnf("error connecting to grpc: %+v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		client := drand.NewPublicClient(conn)
		err = g.workRelay(client)
		if err != nil {
			log.Warnf("error relaying: %+v", err)
			time.Sleep(5 * time.Second)
		}
	}
}

func (g *GossipRelayNode) workRelay(client drand.PublicClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	curr, err := client.PublicRand(ctx, &drand.PublicRandRequest{Round: 0})
	if err != nil {
		return xerrors.Errorf("getting initial round failed: %w", err)
	}
	log.Infof("got latest rand: %d", curr.Round)

	// context.Background() on purpose as this applies to whole, long lived stream
	stream, err := client.PublicRandStream(context.Background(), &drand.PublicRandRequest{Round: curr.Round})
	if err != nil {
		return xerrors.Errorf("getting rand stream: %w", err)
	}

	for {
		select {
		case <-g.done:
			return xerrors.Errorf("relay shutdown")
		default:
		}
		rand, err := stream.Recv()
		if err != nil {
			return xerrors.Errorf("receving on stream: %w", err)
		}

		randB, err := proto.Marshal(rand)
		if err != nil {
			return xerrors.Errorf("marshaling: %w", err)
		}

		err = g.t.Publish(context.TODO(), randB)
		if err != nil {
			return xerrors.Errorf("publishing on pubsub: %w", err)
		}
		log.Infof("Published randomness on pubsub, round: %d", rand.Round)
	}
}
