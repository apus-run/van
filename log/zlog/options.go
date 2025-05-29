package zlog

import (
	"io"
	"os"

	"go.uber.org/zap/zapcore"
)

// Option is config option.
type Option func(*Options)

type Options struct {
	// logger options
	Mode     string              // dev or prod
	Writer   zapcore.WriteSyncer // 日志输出
	LogLevel string              // debug, info, warn, error, panic, panic, fatal
	Format   Format              // text or json
}

// DefaultOptions .
func DefaultOptions() *Options {
	return &Options{
		Mode:     "dev",
		Writer:   zapcore.AddSync(os.Stdout),
		LogLevel: "info", // zapcore.InfoLevel.String(),
		Format:   FormatText,
	}
}

func Apply(opts ...Option) *Options {
	options := DefaultOptions()
	for _, o := range opts {
		o(options)
	}
	return options
}

// WithMode 运行模式
func WithMode(mode string) Option {
	return func(o *Options) {
		o.Mode = mode
	}
}

// WithLogLevel 日志级别
func WithLogLevel(level string) Option {
	return func(o *Options) {
		o.LogLevel = level
	}
}

// WithFormat 日志格式
func WithFormat(format Format) Option {
	return func(o *Options) {
		o.Format = format
	}
}

// WithWriter 日志输出
func WithWriter(writer io.Writer) Option {
	if writer == nil {
		writer = os.Stdout // 默认输出到标准输出
	}
	return func(o *Options) {
		o.Writer = zapcore.AddSync(writer)
	}
}

// WithWriters 日志输出
func WithWriters(writers ...io.Writer) Option {
	return func(o *Options) {
		wrs := make([]zapcore.WriteSyncer, 0, len(writers))
		for _, w := range writers {
			wrs = append(wrs, zapcore.AddSync(w))
		}
		o.Writer = zapcore.NewMultiWriteSyncer(wrs...)
	}
}

// WithOptions 设置所有配置
func WithOptions(fn func(options *Options)) Option {
	return func(options *Options) {
		fn(options)
	}
}
