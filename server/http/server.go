package http

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"net/url"

	"github.com/apus-run/van/server"
	"github.com/apus-run/van/server/internal/endpoint"
	"github.com/apus-run/van/server/internal/host"
	"github.com/apus-run/van/server/internal/shutdown"
)

var _ server.Server = (*Server)(nil)

type Server struct {
	*http.Server
	options *server.ServerOptions
}

func NewServer(opts ...server.ServerOption) *Server {
	options := Apply(opts...)

	srv := &Server{
		options: options,
	}

	srv.Server = &http.Server{
		Handler:   srv,
		TLSConfig: options.TLSConfig,
	}

	return srv
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.listenAndEndpoint(); err != nil {
		return err
	}

	s.BaseContext = func(listener net.Listener) context.Context {
		return ctx
	}

	var err error
	if s.options.TLSConfig != nil {
		slog.Info("[HTTPS] server listen on", "address", s.options.Address)
		err = s.ServeTLS(s.options.Listener, "", "")
	} else {
		slog.Info("[HTTP] server listen on", "address", s.options.Address)
		err = s.Serve(s.options.Listener)
	}

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return shutdown.ShutdownWithContext(ctx, func(ctx context.Context) error {
		return s.Server.Shutdown(ctx)
	}, func() error {
		if err := s.Server.Close(); err != nil {
			return err
		}

		return nil
	})
}

// Endpoint return a real address to registry endpoint.
// examples:
//
//	https://127.0.0.1:8000
//	Legacy: http://127.0.0.1:8000?isSecure=false
func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, s.options.Error
	}
	return s.options.Endpoint, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.options.Handler.ServeHTTP(w, r)
}

// Health
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
		s.options.Endpoint = endpoint.NewEndpoint(endpoint.Scheme("http", s.options.TLSConfig != nil), addr)
	}
	return s.options.Error
}
