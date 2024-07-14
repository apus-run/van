package zlog

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const loggerKeyForGin = "GinXLogger"

// contextKey is how we find Loggers in a context.Context.
type contextKey struct{}

// WithContext returns a copy of context in which the log value is set.
func WithContext(ctx context.Context, keyvals ...any) context.Context {
	if l := FromContext(ctx); l != nil {
		return l.(*ZapLogger).WithContext(ctx, keyvals...)
	}
	return L().WithContext(ctx, keyvals...)
}

func (l *ZapLogger) WithContext(ctx context.Context, keyvals ...any) context.Context {
	with := func(l Logger) context.Context {
		return context.WithValue(ctx, contextKey{}, l)
	}

	// In order to be compatible with *gin.Context, the performance loss is acceptable.
	if c, ok := ctx.(*gin.Context); ok {
		with = func(l Logger) context.Context {
			c.Set(loggerKeyForGin, l)
			return c
		}
	}

	keylen := len(keyvals)
	if keylen == 0 || keylen%2 != 0 {
		return with(l)
	}

	data := make([]zap.Field, 0, (keylen/2)+1)
	for i := 0; i < keylen; i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	return with(l.With(data...))
}

// FromContext returns a logger with predefined values from a context.Context.
func FromContext(ctx context.Context, keyvals ...any) Logger {
	var key any = contextKey{}
	if _, ok := ctx.(*gin.Context); ok {
		key = loggerKeyForGin
	}

	var log Logger = L()
	if ctx != nil {
		if logger, ok := ctx.Value(key).(Logger); ok {
			log = logger
		}
	}

	keylen := len(keyvals)
	if keylen == 0 || keylen%2 != 0 {
		return log
	}

	data := make([]zap.Field, 0, (keylen/2)+1)
	for i := 0; i < keylen; i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	return log.With(data...)
}

// C represents for `FromContext` with empty keyvals.
// zlog.C(ctx).Info("Set function called")
func C(ctx context.Context) Logger {
	return FromContext(ctx).AddCallerSkip(-1)
}
