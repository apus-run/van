package grpc

import (
	"context"
	"github.com/apus-run/van/server"
	"google.golang.org/grpc"
	"net"
)

var _ server.Server = (*Server)(nil)

type Server struct {
	opts *Options
	srv  *grpc.Server
}

func NewServer(grpcServer *grpc.Server, opts ...Option) *Server {
	options := Apply(opts...)
	srv := &Server{
		opts: options,
		srv:  grpcServer,
	}
	return srv
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.opts.Addr)
	if err != nil {
		return err
	}
	return s.srv.Serve(lis)
}

func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.opts.ShutdownTimeout)
	defer cancel()

	s.srv.GracefulStop()
	return nil
}
