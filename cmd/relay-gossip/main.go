package main

import (
	"context"
	"fmt"
	"os"

	"github.com/drand/drand/client/grpc"
	"github.com/drand/drand/cmd/client/lib"
	dlog "github.com/drand/drand/log"
	"github.com/drand/drand/lp2p"
	"github.com/drand/drand/metrics"
	"github.com/drand/drand/metrics/pprof"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	peer "github.com/libp2p/go-libp2p-core/peer"
	cli "github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
)

// Automatically set through -ldflags
// Example: go install -ldflags "-X main.version=`git describe --tags` -X main.buildDate=`date -u +%d/%m/%Y@%H:%M:%S` -X main.gitCommit=`git rev-parse HEAD`"
var (
	version   = "master"
	gitCommit = "none"
	buildDate = "unknown"
)

var log = dlog.DefaultLogger

func main() {
	app := &cli.App{
		Name:     "drand-relay-gossip",
		Version:  version,
		Usage:    "pubsub relay for drand randomness beacon",
		Commands: []*cli.Command{runCmd, clientCmd, idCmd},
	}
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("drand gossip relay %v (date %v, commit %v)\n", version, buildDate, gitCommit)
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}
}

var idFlag = &cli.StringFlag{
	Name:  "identity",
	Usage: "path to a file containing a libp2p identity (base64 encoded)",
	Value: "identity.key",
}

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "starts a drand gossip relay process",
	Flags: []cli.Flag{
		lib.GRPCConnectFlag,
		&cli.StringSliceFlag{
			Name:  "http-connect",
			Usage: "URL(s) of drand HTTP API(s) to relay",
		},
		&cli.StringFlag{
			Name:  "store",
			Usage: "datastore directory",
			Value: "./datastore",
		},
		lib.CertFlag,
		lib.InsecureFlag,
		&cli.StringFlag{
			Name:  "listen",
			Usage: "listening address for libp2p",
			Value: "/ip4/0.0.0.0/tcp/44544",
		},
		&cli.StringFlag{
			Name:  "metrics",
			Usage: "local host:port to bind a metrics servlet (optional)",
		},
		lib.HashFlag,
		lib.RelayFlag,
		idFlag,
	},

	Action: func(cctx *cli.Context) error {
		if cctx.IsSet("metrics") {
			metricsListener := metrics.Start(cctx.String("metrics"), pprof.WithProfile(), nil)
			defer metricsListener.Close()
			metrics.PrivateMetrics.Register(grpc_prometheus.DefaultClientMetrics)
		}

		chainHash := cctx.String(lib.HashFlag.Name)
		if chainHash == "" {
			return xerrors.Errorf("missing required chain-hash parameter")
		}
		grpclient, err := grpc.New(cctx.String(lib.GRPCConnectFlag.Name), cctx.String(lib.CertFlag.Name), cctx.Bool(lib.InsecureFlag.Name))
		if err != nil {
			return err
		}
		cfg := &lp2p.GossipRelayConfig{
			ChainHash:    chainHash,
			PeerWith:     cctx.StringSlice(lib.RelayFlag.Name),
			Addr:         cctx.String("listen"),
			DataDir:      cctx.String("store"),
			IdentityPath: cctx.String(idFlag.Name),
			Client:       grpclient,
		}
		if _, err := lp2p.NewGossipRelayNode(log, cfg); err != nil {
			return err
		}
		<-(chan int)(nil)
		return nil
	},
}

var clientCmd = &cli.Command{
	Name:  "client",
	Flags: lib.ClientFlags,
	Action: func(cctx *cli.Context) error {
		c, err := lib.Create(cctx, false)
		if err != nil {
			return xerrors.Errorf("constructing client: %w", err)
		}

		for rand := range c.Watch(context.Background()) {
			log.Info("client", "got randomness", "round", rand.Round(), "signature", rand.Signature()[:16])
		}

		return nil
	},
}

var idCmd = &cli.Command{
	Name:  "peerid",
	Usage: "prints the libp2p peer ID or creates one if it does not exist",
	Flags: []cli.Flag{idFlag},
	Action: func(cctx *cli.Context) error {
		priv, err := lp2p.LoadOrCreatePrivKey(cctx.String(idFlag.Name), log)
		if err != nil {
			return xerrors.Errorf("loading p2p key: %w", err)
		}
		peerID, err := peer.IDFromPrivateKey(priv)
		if err != nil {
			return xerrors.Errorf("computing peerid: %w", err)
		}
		fmt.Printf("%s\n", peerID)
		return nil
	},
}
