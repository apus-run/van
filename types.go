package van

import (
	"context"
	"net/url"
	"os"

	"github.com/apus-run/van/registry"
	"github.com/apus-run/van/server"
)

// Option is an application option.
type Option func(o *options)

// options is an application options.
type options struct {
	// service id
	id string
	// service name
	name string
	// service version
	version string
	// metadata
	metadata map[string]string
	// server endpoints
	endpoints []*url.URL

	// registry
	registry registry.Registry
	// sevice servers
	servers []server.Server

	context context.Context
	signals []os.Signal

	// Before and After funcs
	beforeStart []func(context.Context) error
	beforeStop  []func(context.Context) error
	afterStart  []func(context.Context) error
	afterStop   []func(context.Context) error
}

// With with service id.
func WithID(id string) Option {
	return func(o *options) { o.id = id }
}

// WithName with service name.
func WithName(name string) Option {
	return func(o *options) { o.name = name }
}

// WithVersion with service version.
func WithVersion(version string) Option {
	return func(o *options) { o.version = version }
}

// WithMetadata with service metadata.
func WithMetadata(md map[string]string) Option {
	return func(o *options) { o.metadata = md }
}

// WithEndpoint with service endpoint.
func WithEndpoint(endpoints ...*url.URL) Option {
	return func(o *options) { o.endpoints = endpoints }
}

// WithContext with service context.
func WithContext(ctx context.Context) Option {
	return func(o *options) { o.context = ctx }
}

// WithServer with transport servers.
func WithServer(srv ...server.Server) Option {
	return func(o *options) { o.servers = srv }
}

// WithSignal with exit signals.
func WithSignal(sigs ...os.Signal) Option {
	return func(o *options) { o.signals = sigs }
}

// WithRegistry with service registry.
func WithRegistry(r registry.Registry) Option {
	return func(o *options) { o.registry = r }
}

// Before and Afters

// BeforeStart run funcs before app starts
func BeforeStart(fn func(context.Context) error) Option {
	return func(o *options) {
		o.beforeStart = append(o.beforeStart, fn)
	}
}

// BeforeStop run funcs before app stops
func BeforeStop(fn func(context.Context) error) Option {
	return func(o *options) {
		o.beforeStop = append(o.beforeStop, fn)
	}
}

// AfterStart run funcs after app starts
func AfterStart(fn func(context.Context) error) Option {
	return func(o *options) {
		o.afterStart = append(o.afterStart, fn)
	}
}

// AfterStop run funcs after app stops
func AfterStop(fn func(context.Context) error) Option {
	return func(o *options) {
		o.afterStop = append(o.afterStop, fn)
	}
}
