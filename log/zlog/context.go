package zlog

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// contextLogKey is how we find Loggers in a context.Context.
type contextLogKey struct{}

// WithContext returns a copy of context in which the log value is set.
func WithContext(ctx context.Context, keyvals ...any) context.Context {
	if l := FromContext(ctx); l != nil {
		return l.(*ZapLogger).WithContext(ctx, keyvals...)
	}
	return L().WithContext(ctx, keyvals...)
}

func (l *ZapLogger) WithContext(ctx context.Context, keyvals ...any) context.Context {
	with := func(l Logger) context.Context {
		return context.WithValue(ctx, contextLogKey{}, l)
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
	var log Logger = L()
	if ctx != nil {
		if logger, ok := ctx.Value(contextLogKey{}).(Logger); ok {
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
func C(ctx context.Context) Logger {
	return FromContext(ctx).AddCallerSkip(-1)
}

// 定义用于上下文的键.
type (
	// requestIDKey 定义请求 ID 的上下文键.
	requestIDKey struct{}
	// userIDKey 定义用户 ID 的上下文键.
	userIDKey struct{}
	// userNameKey 定义用户名的上下文键.
	userNameKey struct{}

	// traceIDkey 定义跟踪 ID 的上下文键.
	traceIDKey struct{}
	// spanIDKey 定义跨度 ID 的上下文键.
	spanIDKey struct{}
)

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, requestID)
}

func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if requestID, ok := ctx.Value(requestIDKey{}).(string); ok {
		return requestID
	}
	return ""
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

func UserIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if userID, ok := ctx.Value(userIDKey{}).(string); ok {
		return userID
	}
	return ""
}

func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, userNameKey{}, username)
}

func UsernameFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if username, ok := ctx.Value(userNameKey{}).(string); ok {
		return username
	}
	return ""
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

func TraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if traceID, ok := ctx.Value(traceIDKey{}).(string); ok {
		return traceID
	}
	return ""
}

func WithSpanID(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, spanIDKey{}, spanID)
}

func SpanIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if spanID, ok := ctx.Value(spanIDKey{}).(string); ok {
		return spanID
	}
	return ""
}
