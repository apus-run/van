package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/apus-run/van/server"
)

var _ server.Server = (*Server)(nil)

type Server struct {
	*http.Server
	options *Options
}

func NewServer(handler http.Handler, opts ...Option) *Server {
	options := Apply(opts...)

	srv := &Server{
		options: options,
	}

	// 初始化 http.Server
	if options.TlsConfig != nil {
		srv.Server = &http.Server{
			Addr:         options.Addr,
			Handler:      handler,
			TLSConfig:    options.TlsConfig,
			IdleTimeout:  options.IdleTimeout,
			ReadTimeout:  options.ReadTimeout,
			WriteTimeout: options.WriteTimeout,
		}
	} else {
		srv.Server = &http.Server{
			Addr:         options.Addr,
			Handler:      handler,
			IdleTimeout:  options.IdleTimeout,
			ReadTimeout:  options.ReadTimeout,
			WriteTimeout: options.WriteTimeout,
		}
	}
	return srv
}

func (s *Server) Start(ctx context.Context) error {
	var err error
	if s.options.TlsConfig != nil {
		err = s.ListenAndServeTLS("", "")
	} else {
		err = s.ListenAndServe()
	}
	if errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.Shutdown(ctx)
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.Handler.ServeHTTP(res, req)
}
