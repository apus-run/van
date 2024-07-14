package http

import (
	"context"
	"errors"
	"github.com/apus-run/van/server"
	"net/http"
)

var _ server.Server = (*Server)(nil)

type Server struct {
	opts *Options
	srv  *http.Server
}

func NewServer(handler http.Handler, opts ...Option) *Server {
	options := Apply(opts...)
	srv := &Server{
		opts: options,
		srv: &http.Server{
			Addr:    options.Addr,
			Handler: handler,
		},
	}
	return srv
}

func (s *Server) Start() error {
	err := s.srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	// 创建 ctx 用于通知服务器 goroutine, 它有 10 秒时间完成当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), s.opts.ShutdownTimeout)
	defer cancel()
	return s.srv.Shutdown(ctx)
}
