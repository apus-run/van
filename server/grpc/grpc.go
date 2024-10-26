package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/apus-run/van/server"
)

var _ server.Server = (*Server)(nil)

type Server struct {
	*grpc.Server
	ctx     context.Context
	options *Options
}

func NewServer(grpcServer *grpc.Server, opts ...Option) *Server {
	options := Apply(opts...)
	srv := &Server{
		options: options,
		ctx:     context.Background(),
		Server:  grpcServer,
	}

	return srv
}

func (s *Server) Start(ctx context.Context) error {
	s.ctx = ctx

	lis, err := net.Listen("tcp", s.options.Addr)
	if err != nil {
		return err
	}
	return s.Serve(lis)
}

func (s *Server) Stop(ctx context.Context) error {
	s.GracefulStop()
	return nil
}
