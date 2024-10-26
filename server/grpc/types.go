package grpc

import (
	"crypto/tls"
	"time"
)

// Option is server option.
type Option func(*Options)

// Options is server options.
type Options struct {
	Addr string

	// TLS config.
	TlsConfig *tls.Config

	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DefaultOptions is server default options.
func DefaultOptions() *Options {
	return &Options{
		Addr:         ":8090",
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

// Apply applies options.
func Apply(opts ...Option) *Options {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return options
}

func WithAddr(addr string) Option {
	return func(options *Options) {
		options.Addr = addr
	}
}

func WithIdleTimeout(t time.Duration) Option {
	return func(options *Options) {
		options.IdleTimeout = t
	}
}

func WithReadTimeout(t time.Duration) Option {
	return func(options *Options) {
		options.ReadTimeout = t
	}
}

func WithWriteTimeout(t time.Duration) Option {
	return func(options *Options) {
		options.WriteTimeout = t
	}
}

// TLSConfig with TLS config.
func TLSConfig(c *tls.Config) Option {
	return func(options *Options) {
		options.TlsConfig = c
	}
}
