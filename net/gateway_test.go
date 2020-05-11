package net

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	run "runtime"
	"testing"
	"time"

	"github.com/drand/drand/protobuf/drand"
	"github.com/kabukky/httpscerts"
	"github.com/stretchr/testify/require"
)

type testPeer struct {
	addr string
	t    bool
}

func (t *testPeer) Address() string {
	return t.addr
}

func (t *testPeer) IsTLS() bool {
	return t.t
}

type testRandomnessServer struct {
	*EmptyServer
	round uint64
}

func (t *testRandomnessServer) NewBeacon(context.Context, *drand.BeaconPacket) (*drand.Empty, error) {
	return new(drand.Empty), errors.New("no beacon")
}
func (t *testRandomnessServer) PublicRand(context.Context, *drand.PublicRandRequest) (*drand.PublicRandResponse, error) {
	return &drand.PublicRandResponse{Round: t.round}, nil
}
func (t *testRandomnessServer) PrivateRand(context.Context, *drand.PrivateRandRequest) (*drand.PrivateRandResponse, error) {
	return &drand.PrivateRandResponse{}, nil
}
func (t *testRandomnessServer) Group(context.Context, *drand.GroupRequest) (*drand.GroupPacket, error) {
	return nil, nil
}
func (t *testRandomnessServer) DistKey(context.Context, *drand.DistKeyRequest) (*drand.DistKeyResponse, error) {
	return nil, nil
}
func (t *testRandomnessServer) Home(context.Context, *drand.HomeRequest) (*drand.HomeResponse, error) {
	return nil, nil
}

func TestListener(t *testing.T) {
	ctx := context.Background()
	hostAddr := "127.0.0.1"
	addrGRPC := hostAddr + ":4000"
	peerGRPC := &testPeer{addrGRPC, false}
	addrREST := hostAddr + ":4001"
	peerREST := &testPeer{addrREST, false}
	randServer := &testRandomnessServer{round: 42}

	lisGRPC := NewGRPCListenerForPublicAndProtocol(addrGRPC, randServer)
	lisREST := NewRESTListenerForPublic(addrREST, randServer)
	go lisGRPC.Start()
	defer lisGRPC.Stop()
	go lisREST.Start()
	defer lisREST.Stop()
	time.Sleep(100 * time.Millisecond)

	// GRPC
	client := NewGrpcClient()
	resp, err := client.PublicRand(ctx, peerGRPC, &drand.PublicRandRequest{})
	require.NoError(t, err)
	expected := &drand.PublicRandResponse{Round: randServer.round}
	require.Equal(t, expected.GetRound(), resp.GetRound())

	//XXX: Why is this producing error logs.
	err = client.PartialBeacon(ctx, peerGRPC, &drand.PartialBeaconPacket{})
	require.Error(t, err)

	// REST
	rest := NewRestClient()
	resp, err = rest.PublicRand(ctx, peerREST, &drand.PublicRandRequest{})
	require.NoError(t, err)
	expected = &drand.PublicRandResponse{Round: randServer.round}
	require.Equal(t, expected.GetRound(), resp.GetRound())
}

// ref https://bbengfort.github.io/programmer/2017/03/03/secure-grpc.html
func TestListenerTLS(t *testing.T) {
	if run.GOOS == "windows" {
		fmt.Println("Skipping TestClientTLS as operating on Windows")
		t.Skip("crypto/x509: system root pool is not available on Windows")
	}
	hostAddr := "127.0.0.1"
	addrGRPC := hostAddr + ":4000"
	peerGRPC := &testPeer{addrGRPC, true}
	addrREST := hostAddr + ":4001"
	peerREST := &testPeer{addrREST, true}

	tmpDir := path.Join(os.TempDir(), "drand-net")
	require.NoError(t, os.MkdirAll(tmpDir, 0766))
	defer os.RemoveAll(tmpDir)
	certPath := path.Join(tmpDir, "server.crt")
	keyPath := path.Join(tmpDir, "server.key")
	if httpscerts.Check(certPath, keyPath) != nil {
		require.NoError(t, httpscerts.Generate(certPath, keyPath, hostAddr))
	}

	randServer := &testRandomnessServer{round: 42}

	lisGRPC, err := NewGRPCListenerForPublicAndProtocolWithTLS(addrGRPC, certPath, keyPath, randServer)
	require.NoError(t, err)
	lisREST, err := NewRESTListenerForPublicWithTLS(addrREST, certPath, keyPath, randServer)
	require.NoError(t, err)

	go lisGRPC.Start()
	defer lisGRPC.Stop()
	go lisREST.Start()
	defer lisREST.Stop()

	time.Sleep(100 * time.Millisecond)

	require.Equal(t, peerGRPC.Address(), addrGRPC)
	require.Equal(t, peerREST.Address(), addrREST)
	certManager := NewCertManager()
	certManager.Add(certPath)

	ctx := context.Background()
	// test GRPC variant
	client := NewGrpcClientFromCertManager(certManager)
	resp, err := client.PublicRand(ctx, peerGRPC, &drand.PublicRandRequest{})
	require.Nil(t, err)
	expected := &drand.PublicRandResponse{Round: randServer.round}
	require.Equal(t, expected.GetRound(), resp.GetRound())

	// test REST variant
	rest := NewRestClientFromCertManager(certManager)
	resp, err = rest.PublicRand(ctx, peerREST, &drand.PublicRandRequest{})
	require.NoError(t, err)
	expected = &drand.PublicRandResponse{Round: randServer.round}
	require.Equal(t, expected.GetRound(), resp.GetRound())
}
