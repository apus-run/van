package van

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/apus-run/van/server/registry"
)

// Service is an application components lifecycle manager.
type Service struct {
	options  *options
	ctx      context.Context
	cancel   context.CancelFunc
	mu       sync.Mutex
	instance *registry.ServiceInstance
}

// New create an application lifecycle manager.
func New(opts ...Option) *Service {
	o := &options{
		context: context.Background(),
		signals: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}
	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}
	for _, opt := range opts {
		opt(o)
	}

	ctx, cancel := context.WithCancel(o.context)
	return &Service{
		ctx:     ctx,
		cancel:  cancel,
		options: o,
	}
}

// ID returns app instance id.
func (s *Service) ID() string { return s.options.id }

// Name returns service name.
func (s *Service) Name() string { return s.options.name }

// Version returns app version.
func (s *Service) Version() string { return s.options.version }

// Metadata returns service metadata.
func (s *Service) Metadata() map[string]string { return s.options.metadata }

// Endpoint returns endpoints.
func (s *Service) Endpoint() []string {
	if s.instance != nil {
		return s.instance.Endpoints
	}
	return nil
}

// Run executes all OnStart hooks registered with the application's Lifecycle.
func (s *Service) Run() error {
	instance, err := s.registryService()
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.instance = instance
	s.mu.Unlock()
	c := NewContext(s.ctx, s)
	eg, ctx := errgroup.WithContext(c)
	wg := sync.WaitGroup{}

	for _, fn := range s.options.beforeStart {
		if err = fn(c); err != nil {
			return err
		}
	}
	for _, srv := range s.options.servers {
		server := srv
		eg.Go(func() error {
			<-ctx.Done() // wait for stop signal
			stopCtx, cancel := context.WithTimeout(NewContext(s.options.context, s), 10*time.Second)
			defer cancel()
			return server.Stop(stopCtx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done() // here is to ensure server start has begun running before register, so defer is not needed
			return server.Start(NewContext(s.options.context, s))
		})
	}
	wg.Wait()
	if s.options.registry != nil {
		rctx, rcancel := context.WithTimeout(ctx, 10*time.Second)
		defer rcancel()
		if err = s.options.registry.Register(rctx, instance); err != nil {
			return err
		}
	}
	for _, fn := range s.options.afterStart {
		if err = fn(c); err != nil {
			return err
		}
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, s.options.signals...)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ch:
			return s.Stop()
		}
	})
	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	err = nil
	for _, fn := range s.options.afterStop {
		err = fn(c)
	}
	return err
}

// Stop gracefully stops the application.
func (s *Service) Stop() (err error) {
	sctx := NewContext(s.ctx, s)
	for _, fn := range s.options.beforeStop {
		err = fn(sctx)
	}

	s.mu.Lock()
	instance := s.instance
	s.mu.Unlock()
	if s.options.registry != nil && instance != nil {
		ctx, cancel := context.WithTimeout(NewContext(s.ctx, s), 10*time.Second)
		defer cancel()
		if err = s.options.registry.Deregister(ctx, instance); err != nil {
			return err
		}
	}
	if s.cancel != nil {
		s.cancel()
	}
	return err
}

func (s *Service) registryService() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 0, len(s.options.endpoints))
	for _, e := range s.options.endpoints {
		endpoints = append(endpoints, e.String())
	}
	if len(endpoints) == 0 {
		for _, srv := range s.options.servers {
			e, err := srv.Endpoint()
			if err != nil {
				return nil, err
			}
			endpoints = append(endpoints, e.String())
		}
	}
	return &registry.ServiceInstance{
		ID:        s.options.id,
		Name:      s.options.name,
		Version:   s.options.version,
		Metadata:  s.options.metadata,
		Endpoints: endpoints,
	}, nil
}

// serviceKey is a context key used to store the service instance into its base context.
type serviceKey struct{}

func NewContext(ctx context.Context, l *Service) context.Context {
	return context.WithValue(ctx, serviceKey{}, l)
}

func FromContext(ctx context.Context) (*Service, bool) {
	if l, ok := ctx.Value(serviceKey{}).(*Service); ok {
		return l, true
	}
	return nil, false
}
