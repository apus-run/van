package server

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"

	"google.golang.org/grpc"
)

type ServerOption func(*ServerOptions)

type ServerOptions struct {
	// server listen network tcp, udp
	Network string
	// server listen address
	Address string

	Listener net.Listener

	// handler for http server
	Handler http.Handler
	// Options for gRPC server
	Options []grpc.ServerOption

	// TLS config.
	TLSConfig *tls.Config

	Endpoint *url.URL
	Error    error
}

func WithNetwork(network string) ServerOption {
	return func(options *ServerOptions) {
		options.Network = network
	}
}

func WithAddress(address string) ServerOption {
	return func(options *ServerOptions) {
		options.Address = address
	}
}

func WithTlsConfig(tlsConfig *tls.Config) ServerOption {
	return func(options *ServerOptions) {
		options.TLSConfig = tlsConfig
	}
}

func WithListener(listener net.Listener) ServerOption {
	return func(options *ServerOptions) {
		options.Listener = listener
	}
}

// WithHandler http server handler
func WithHandler(handler http.Handler) ServerOption {
	return func(options *ServerOptions) {
		options.Handler = handler
	}
}

// WithOptions gRPC server option
func WithOptions(gopts ...grpc.ServerOption) ServerOption {
	return func(options *ServerOptions) {
		options.Options = append(options.Options, gopts...)
	}
}
