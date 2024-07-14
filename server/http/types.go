package http

import "time"

// Option is server option.
type Option func(*Options)

// Options is server options.
type Options struct {
	Addr string

	ShutdownTimeout time.Duration
}

// DefaultOptions is server default options.
func DefaultOptions() *Options {
	return &Options{
		Addr:            ":8080",
		ShutdownTimeout: 10 * time.Second,
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

func WithShutdownTimeout(t time.Duration) Option {
	return func(options *Options) {
		options.ShutdownTimeout = t
	}
}
