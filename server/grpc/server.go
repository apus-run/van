package grpc

import (
	"context"
	"log/slog"
	"net"
	"net/url"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/apus-run/van/server"
	"github.com/apus-run/van/server/internal/endpoint"
	"github.com/apus-run/van/server/internal/host"
	"github.com/apus-run/van/server/internal/shutdown"
)

var _ server.Server = (*Server)(nil)

type Server struct {
	*grpc.Server

	baseCtx context.Context
	options *server.ServerOptions
}

func NewServer(opts ...server.ServerOption) *Server {
	options := Apply(opts...)

	srv := &Server{
		options: options,
	}

	grpcOpts := []grpc.ServerOption{}

	if options.TLSConfig != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(options.TLSConfig)))
	}

	if len(options.Options) > 0 {
		grpcOpts = append(grpcOpts, options.Options...)
	}

	srv.Server = grpc.NewServer(grpcOpts...)

	return srv
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.listenAndEndpoint(); err != nil {
		return s.options.Error
	}
	s.baseCtx = ctx

	slog.Info("[gRPC] server listen on", "address", s.options.Address)

	return s.Serve(s.options.Listener)
}

func (s *Server) Stop(ctx context.Context) error {
	return shutdown.ShutdownWithContext(ctx, func(_ context.Context) error {
		s.Server.GracefulStop()
		return nil
	}, func() error {
		s.Server.Stop()

		return nil
	})
}
func (s *Server) Health() bool {
	if s.options.Listener == nil {
		return false
	}
	conn, err := s.options.Listener.Accept()
	if err != nil {
		return false
	}
	e := conn.Close()
	return e == nil
}

// Endpoint return a real address to registry endpoint.
// examples:
//
//	grpc://127.0.0.1:9000?isSecure=false
func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, s.options.Error
	}
	return s.options.Endpoint, nil
}

func (s *Server) listenAndEndpoint() error {
	if s.options.Listener == nil {
		lis, err := net.Listen(s.options.Network, s.options.Address)
		if err != nil {
			s.options.Error = err
			return err
		}
		s.options.Listener = lis
	}
	if s.options.Endpoint == nil {
		addr, err := host.Extract(s.options.Address, s.options.Listener)
		if err != nil {
			s.options.Error = err
			return err
		}
		s.options.Endpoint = endpoint.NewEndpoint(endpoint.Scheme("grpc", s.options.TLSConfig != nil), addr)
	}
	return s.options.Error
}
