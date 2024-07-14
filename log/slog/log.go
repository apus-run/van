// https://tonybai.com/2023/09/01/slog-a-new-choice-for-logging-in-go/

package slog

import (
	"fmt"

	"go/build"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	errorsx "github.com/apus-run/van/log/slog/errors"
	"gopkg.in/natefinch/lumberjack.v2"
)

const LOGGER_KEY = "slogLogger"
const AttrErrorKey = "error"

var factory = &LoggerFactory{
	loggers: make(map[string]*slog.Logger),
}

type LoggerFactory struct {
	mu      sync.Mutex
	loggers map[string]*slog.Logger
}

var defaultHandler *Handler

func NewLogger(opts ...Option) *slog.Logger {
	options := Apply(opts...)

	factory.mu.Lock()
	if logger, ok := factory.loggers[options.LogFilename]; ok {
		factory.mu.Unlock()
		return logger
	}
	defer factory.mu.Unlock()

	// 日志文件切割归档
	writerSyncer := getLogWriter(options)

	// 日志级别
	level := getLogLevel(options.LogLevel)

	var handler slog.Handler
	handlerOptions := &slog.HandlerOptions{
		Level:       level,
		AddSource:   true,
		ReplaceAttr: ReplaceAttr,
	}
	if len(options.LogFilename) == 0 && strings.ToLower(options.Encoding) == "console" {
		handler = slog.NewTextHandler(os.Stdout, handlerOptions)
	} else {
		handler = slog.NewJSONHandler(writerSyncer, handlerOptions)
	}

	if options.LogGroup != "" {
		handler = handler.WithGroup(options.LogGroup)
	}
	if len(options.LogAttrs) > 0 {
		handler = handler.WithAttrs(options.LogAttrs)
	}

	defaultHandler = NewHandler(handler).(*Handler)
	logger := slog.New(defaultHandler)
	// 此处设置默认日志, 最好手动设置
	// slog.SetDefault(l)

	factory.loggers[options.LogFilename] = logger

	logger.Info("the log module has been initialized successfully.", slog.Any("options", options))

	return logger
}

func getLogWriter(opts *Options) io.WriteCloser {
	return &lumberjack.Logger{
		Filename:   opts.LogFilename,
		MaxSize:    opts.MaxSize, // megabytes
		MaxBackups: opts.MaxBackups,
		MaxAge:     opts.MaxAge, //days
		Compress:   opts.Compress,
	}
}

func getLogLevel(logLevel string) slog.Level {
	level := new(slog.Level)
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		return slog.LevelError
	}

	return *level
}

func ApplyHandlerOption(opt HandlerOption) {
	defaultHandler.Apply(opt)
}

// ReplaceAttr handle log key-value pair
func ReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.TimeKey:
		return slog.String(a.Key, a.Value.Time().Format(time.RFC3339))
	case slog.LevelKey:
		return slog.String(a.Key, strings.ToLower(a.Value.String()))
	case slog.SourceKey:
		if v, ok := a.Value.Any().(*slog.Source); ok {
			a.Value = slog.StringValue(fmt.Sprintf("%s:%d", getBriefSource(v.File), v.Line))
		}
		return a
	case AttrErrorKey:
		v, ok := a.Value.Any().(interface {
			StackTrace() errorsx.StackTrace
		})
		if ok {
			st := v.StackTrace()
			return slog.Any(a.Key, slog.GroupValue(
				slog.String("msg", a.Value.String()),
				slog.Any("stack", errorsx.StackTrace(st)),
			))
		}
		return a
	}
	return a
}

func projectPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		return filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	}
	return ""
}

func getBriefSource(source string) string {
	gp := filepath.ToSlash(build.Default.GOPATH)
	if strings.HasPrefix(source, gp) {
		return strings.TrimPrefix(source, gp+"/pkg/mod/")
	}
	pp := filepath.ToSlash(projectPath())
	return strings.TrimPrefix(source, pp+"/")
}
