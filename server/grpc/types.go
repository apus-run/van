package grpc

import (
	"github.com/apus-run/van/server"
)

// DefaultOptions is server default options.
func DefaultOptions() *server.ServerOptions {
	return &server.ServerOptions{
		Network: "tcp",
		Address: ":0",
	}
}

// Apply applies options.
func Apply(opts ...server.ServerOption) *server.ServerOptions {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return options
}
