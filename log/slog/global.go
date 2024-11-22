package slog

import (
	"context"
	"sync"
)

// globalLogger is designed as a global logger in current process.
var global = &loggerAppliance{}

// loggerAppliance is the proxy of `Logger` to
// make logger change will affect all sub-logger.
type loggerAppliance struct {
	lock sync.Mutex
	Logger
}

func init() {
	logger := NewLogger()

	global.SetLogger(logger)
}

func (a *loggerAppliance) SetLogger(in Logger) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.Logger = in
}

// SetLogger should be called before any other log call.
// And it is NOT THREAD SAFE.
func SetLogger(logger Logger) {
	global.SetLogger(logger)
}

// GetLogger returns global logger appliance as logger in current process.
func GetLogger() Logger {
	return global.Logger
}

// L 是 GetLogger 简写
func L() Logger {
	return global.Logger
}

// Default 是 GetLogger
func Default() Logger {
	return global.Logger
}

func Info(msg string, args ...Attr) {
	L().Info(msg, args...)
}

func Infof(format string, args ...any) {
	L().Infof(format, args...)
}

func InfoContext(ctx context.Context, msg string, args ...Attr) {
	L().InfoContext(ctx, msg, args...)
}

func Error(msg string, args ...Attr) {
	L().Error(msg, args...)
}

func Errorf(format string, args ...any) {
	L().Errorf(format, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...Attr) {
	L().ErrorContext(ctx, msg, args...)
}

func Debug(msg string, args ...Attr) {
	L().Debug(msg, args...)
}

func Debugf(format string, args ...any) {
	L().Debugf(format, args...)
}

func DebugContext(ctx context.Context, msg string, args ...Attr) {
	L().DebugContext(ctx, msg, args...)
}

func Warn(msg string, args ...Attr) {
	L().Warn(msg, args...)
}

func Warnf(format string, args ...any) {
	L().Warnf(format, args...)
}

func WarnContext(ctx context.Context, msg string, args ...Attr) {
	L().WarnContext(ctx, msg, args...)
}

func Fatal(msg string, args ...Attr) {
	L().Fatal(msg, args...)
}

func Fatalf(format string, args ...any) {
	L().Fatalf(format, args...)
}

func FatalContext(ctx context.Context, msg string, args ...Attr) {
	L().FatalContext(ctx, msg, args...)
}

func Panic(msg string, args ...Attr) {
	L().Panic(msg, args...)
}

func Panicf(format string, args ...any) {
	L().Panicf(format, args...)
}

func PanicContext(ctx context.Context, msg string, args ...Attr) {
	L().PanicContext(ctx, msg, args...)
}

func With(args ...any) *SlogLogger {
	return L().With(args...)
}

func WithGroup(name string) *SlogLogger {
	return L().WithGroup(name)
}
