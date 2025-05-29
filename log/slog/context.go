package slog

import (
	"context"
)

type contextLogKey struct{}

func NewContext(ctx context.Context, l *SlogLogger) context.Context {
	return context.WithValue(ctx, contextLogKey{}, l)
}

func WithContext(ctx context.Context, l *SlogLogger) context.Context {
	if _, ok := ctx.Value(contextLogKey{}).(*SlogLogger); ok {
		return ctx
	}
	return context.WithValue(ctx, contextLogKey{}, l)
}

func FromContext(ctx context.Context) *SlogLogger {
	if l, ok := ctx.Value(contextLogKey{}).(*SlogLogger); ok {
		return l
	}
	return nil
}

// C represents for `FromContext` with empty keyvals.
// slog.C(ctx).Info("Set function called")
func C(ctx context.Context) *SlogLogger {
	return FromContext(ctx)
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
