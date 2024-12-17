package http

import (
	"net/http"

	"github.com/apus-run/van/server"
)

// DefaultOptions is server default options.
func DefaultOptions() *server.ServerOptions {
	return &server.ServerOptions{
		Network: "tcp",
		Address: ":0",
		Handler: http.DefaultServeMux,
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
