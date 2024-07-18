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

func (s *Server) Start(ctx context.Context) error {
	err := s.srv.ListenAndServe()
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	// 创建 ctx 用于通知服务器 goroutine, 它有 10 秒时间完成当前正在处理的请求
	ctx, cancel := context.WithTimeout(ctx, s.opts.ShutdownTimeout)
	defer cancel()
	return s.srv.Shutdown(ctx)
}
