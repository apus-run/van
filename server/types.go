package server

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"

	"google.golang.org/grpc"
)

// ServerOptions 定义了服务器的配置选项。
type ServerOptions struct {
	// Network 指定服务器监听的网络类型，如 "tcp" 或 "udp"。
	Network string
	// Address 指定服务器监听的地址。
	Address string

	Listener net.Listener

	// Handler 指定 HTTP 服务器的处理器。
	Handler http.Handler
	// Options 指定 gRPC 服务器的选项。
	Options []grpc.ServerOption

	// TLSConfig 指定 TLS 配置。
	TLSConfig *tls.Config

	Endpoint *url.URL
	Error    error
}

// ServerOption 是一个函数类型，用于设置 ServerOptions 的各个字段。
type ServerOption func(*ServerOptions)

// WithNetwork 设置服务器监听的网络类型。
func WithNetwork(network string) ServerOption {
	return func(options *ServerOptions) {
		options.Network = network
	}
}

// WithAddress 设置服务器监听的地址。
func WithAddress(address string) ServerOption {
	return func(options *ServerOptions) {
		options.Address = address
	}
}

// WithTlsConfig 设置服务器的 TLS 配置。
func WithTlsConfig(tlsConfig *tls.Config) ServerOption {
	return func(options *ServerOptions) {
		options.TLSConfig = tlsConfig
	}
}

// WithListener 设置服务器的监听器。
func WithListener(listener net.Listener) ServerOption {
	return func(options *ServerOptions) {
		options.Listener = listener
	}
}

// WithHandler 设置 HTTP 服务器的处理器。
func WithHandler(handler http.Handler) ServerOption {
	return func(options *ServerOptions) {
		options.Handler = handler
	}
}

// WithOptions 设置 gRPC 服务器的选项。
func WithOptions(gopts ...grpc.ServerOption) ServerOption {
	return func(options *ServerOptions) {
		options.Options = append(options.Options, gopts...)
	}
}
