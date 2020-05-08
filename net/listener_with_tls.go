package net

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"github.com/drand/drand/protobuf/drand"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/nikkolasg/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// NewGRPCListenerForPublicAndProtocolWithTLS creates a new listener for the Public and Protocol APIs over GRPC with TLS.
func NewGRPCListenerForPublicAndProtocolWithTLS(bindingAddr string, certPath, keyPath string, s Service, opts ...grpc.ServerOption) (Listener, error) {
	lis, err := net.Listen("tcp", bindingAddr)
	if err != nil {
		return nil, err
	}

	x509KeyPair, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	grpcCreds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	opts = append(opts, grpc.Creds(grpcCreds))
	opts = append(opts, grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor))
	opts = append(opts, grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor))
	serverOpts := append(opts, grpc.Creds(grpcCreds))
	grpcServer := grpc.NewServer(serverOpts...)
	drand.RegisterPublicServer(grpcServer, s)
	drand.RegisterProtocolServer(grpcServer, s)

	httpServer := buildTLSServer(grpcServer, x509KeyPair)
	g := &tlsListener{
		Service:    s,
		server:     httpServer,
		grpcServer: grpcServer,
		l:          tls.NewListener(lis, httpServer.TLSConfig),
	}
	grpc_prometheus.Register(g.grpcServer)
	return g, nil
}

// NewRESTListenerForPublicWithTLS creates a new listener for the Public API over REST with TLS.
func NewRESTListenerForPublicWithTLS(bindingAddr string, certPath, keyPath string, s Service, opts ...grpc.ServerOption) (Listener, error) {
	lis, err := net.Listen("tcp", bindingAddr)
	if err != nil {
		return nil, err
	}

	x509KeyPair, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	grpcCreds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	opts = append(opts, grpc.Creds(grpcCreds))
	opts = append(opts, grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor))
	opts = append(opts, grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor))
	serverOpts := append(opts, grpc.Creds(grpcCreds))
	grpcServer := grpc.NewServer(serverOpts...)
	drand.RegisterPublicServer(grpcServer, s)
	drand.RegisterProtocolServer(grpcServer, s) //XXX

	o := runtime.WithMarshalerOption("*", defaultJSONMarshaller)
	gwMux := runtime.NewServeMux(o)
	proxy := &drandProxy{s}
	err = drand.RegisterPublicHandlerClient(context.Background(), gwMux, proxy)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle("/", gwMux)

	httpServer := buildTLSServer(grpcHandlerFunc(grpcServer, mux), x509KeyPair)
	g := &tlsListener{
		Service:    s,
		server:     httpServer,
		grpcServer: grpcServer,
		l:          tls.NewListener(lis, httpServer.TLSConfig),
	}
	grpc_prometheus.Register(g.grpcServer)
	return g, nil
}

func buildTLSServer(grpcServer http.Handler, x509KeyPair tls.Certificate) *http.Server {
	return &http.Server{
		Handler: grpcServer, // grpcHandlerFunc(grpcServer, mux),
		TLSConfig: &tls.Config{
			// From https://blog.cloudflare.com/exposing-go-on-the-internet/

			// Causes servers to use Go's default ciphersuite preferences,
			// which are tuned to avoid attacks. Does nothing on clients.
			PreferServerCipherSuites: true,

			// Only use curves which have assembly implementations
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519,
			},

			// Drand clients and servers are all modern software, and so we
			// can require TLS 1.2 and the best cipher suites.
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
			// End Cloudflare recommendations.

			Certificates: []tls.Certificate{x509KeyPair},
			NextProtos:   []string{"h2"},
		},
	}
}

type tlsListener struct {
	Service
	server     *http.Server
	grpcServer *grpc.Server
	l          net.Listener
}

func (g *tlsListener) Start() {
	if err := g.server.Serve(g.l); err != nil {
		fmt.Printf("grpc: tls listener start failed: %s", err)
	}
}

func (g *tlsListener) Stop() {
	// Graceful stop not supported with HTTP Server
	// https://github.com/grpc/grpc-go/issues/1384
	if err := g.server.Shutdown(context.TODO()); err != nil {
		slog.Debugf("grpc: tls listener shutdown failed: %s\n", err)
	}
	g.grpcServer.Stop()
	g.l.Close()
}
