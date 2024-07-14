package grpc

import (
	"google.golang.org/grpc"
	"net"

	"github.com/apus-run/van/server"
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

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.opts.Addr)
	if err != nil {
		return err
	}
	return s.srv.Serve(lis)
}

func (s *Server) Stop() error {
	s.srv.GracefulStop()
	return nil
}
