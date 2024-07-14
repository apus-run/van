package grpc

// Option is server option.
type Option func(*Options)

// Options is server options.
type Options struct {
	Addr string
}

// DefaultOptions is server default options.
func DefaultOptions() *Options {
	return &Options{
		Addr: ":8090",
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
