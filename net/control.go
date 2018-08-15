package net

import (
	"context"
	"fmt"
	"log"

	"github.com/dedis/drand/protobuf/control"
	"github.com/dedis/drand/protobuf/crypto"
	"github.com/dedis/kyber"
	"github.com/nikkolasg/slog"
	"google.golang.org/grpc"
)

func RequestShare() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := control.NewControlClient(conn)
	response, err := c.Share(context.Background(), &control.ShareRequest{})
	if err != nil {
		log.Fatalf("Error when calling Share: %s", err)
	}
	log.Printf("Response: %s", response.Share)
}

type Server struct {
	S kyber.Scalar
}

func (s *Server) Share(ctx context.Context, in *control.ShareRequest, opts ...grpc.CallOption) (*control.ShareResponse, error) {
	share, err := crypto.KyberToProtoScalar(s.S)
	if err != nil {
		slog.Fatal("drand: could not load the private share")
	}
	fmt.Printf("Received request for share and returned the share : %s", share)
	return &control.ShareResponse{Share: share}, nil
}
