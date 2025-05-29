package zlog

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const LOGGER_KEY = "zapLogger"

var _ Logger = (*ZapLogger)(nil)

type ZapLogger struct {
	*zap.Logger
}

// NewZapLogger 只包装了 zap
func NewZapLogger(l *zap.Logger) Logger {
	return &ZapLogger{
		l,
	}
}

func ZapDevConfig() zap.Config {
	return zap.NewDevelopmentConfig()
}

func ZapProdConfig() zap.Config {
	return zap.NewProductionConfig()
}

// NewLogger 包装了 zap 和日志文件切割归档
func NewLogger(opts ...Option) *ZapLogger {
	options := Apply(opts...)

	// 编码器配置
	encoder := getEncoder(options)

	// 日志级别
	level := getLogLevel(options.LogLevel)

	core := zapcore.NewCore(encoder, options.Writer, level)

	return &ZapLogger{
		buildZapLogger(core, options.Mode),
	}
}

// buildZapLogger 封装核心构建逻辑
func buildZapLogger(core zapcore.Core, mode string) *zap.Logger {
	baseOpts := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1), // 跳过封装层
		zap.AddStacktrace(zap.ErrorLevel),
	}

	if mode != "prod" {
		return zap.New(core, append(baseOpts, zap.Development())...)
	}
	return zap.New(core, baseOpts...)
}

func getEncoder(options *Options) zapcore.Encoder {
	// 根据不同的日志格式创建不同的编码器
	switch options.Format {
	case FormatJSON:
		return zapcore.NewJSONEncoder(getJSONEncoderConfig())
	case FormatText:
		return zapcore.NewConsoleEncoder(getConsoleEncoderConfig())
	default:
		return zapcore.NewConsoleEncoder(getConsoleEncoderConfig())
	}
}

// getConsoleEncoderConfig 开发友好的控制台编码配置
func getConsoleEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 在日志文件中使用大写字母记录日志级别
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

// getJSONEncoderConfig 生产环境 JSON 编码配置
func getJSONEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     customTimeEncoder,             // 自定义时间格式
		EncodeDuration: zapcore.MillisDurationEncoder, // 毫秒精度
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// customTimeEncoder 自定义时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	// 指定时间序列化函数，将时间序列化为 `2006-01-02 15:04:05.000` 格式，更易读
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func getLogLevel(logLevel string) zapcore.Level {
	if logLevel == "" {
		return zapcore.InfoLevel // 默认级别
	}
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(logLevel)); err != nil {
		// 如果指定了非法的日志级别，则默认使用 error 级别
		return zap.ErrorLevel
	}

	return zapLevel
}

func (l *ZapLogger) Info(msg string, tags ...Field) {
	l.Logger.Info(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Error(msg string, tags ...Field) {
	l.Logger.Error(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Debug(msg string, tags ...Field) {
	l.Logger.Debug(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Warn(msg string, tags ...Field) {
	l.Logger.Warn(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Fatal(msg string, tags ...Field) {
	l.Logger.Fatal(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Panic(msg string, tags ...Field) {
	l.Logger.Panic(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Slow(msg string, tags ...Field) {
	l.Logger.Warn(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Stack(msg string) {
	l.Logger.Error(fmt.Sprint(msg), zap.Stack("stack"))
}

func (l *ZapLogger) Stat(msg string, tags ...Field) {
	l.Logger.Info(msg, l.toZapFields(tags)...)
}

func (l *ZapLogger) Debugf(format string, args ...any) {
	l.Logger.Sugar().Debugf(format, args...)
}

func (l *ZapLogger) Infof(format string, args ...any) {
	l.Logger.Sugar().Infof(format, args...)
}

func (l *ZapLogger) Warnf(format string, args ...any) {
	l.Logger.Sugar().Warnf(format, args...)
}

func (l *ZapLogger) Errorf(format string, args ...any) {
	l.Logger.Sugar().Errorf(format, args...)
}

func (l *ZapLogger) Fatalf(format string, args ...any) {
	l.Logger.Sugar().Fatalf(format, args...)
}

func (l *ZapLogger) Panicf(format string, args ...any) {
	l.Logger.Sugar().Panicf(format, args...)
}

func (l *ZapLogger) Print(args ...any) {
	l.Logger.Info(fmt.Sprint(args...))
}

func (l *ZapLogger) Printf(format string, args ...any) {
	l.Logger.Sugar().Infof(format, args...)
}

// Sync 刷新缓冲 (注: 文件同步可能报错)
func (l *ZapLogger) Sync() error {
	return l.Logger.Sync()
}

// Close 兼容关闭接口
func (l *ZapLogger) Close() error {
	return l.Sync()
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (l *ZapLogger) With(fields ...Field) Logger {
	if len(fields) == 0 {
		return l
	}

	lc := l.clone()
	lc.Logger = lc.Logger.With(l.toZapFields(fields)...)
	return lc
}

// AddCallerSkip increases the number of callers skipped by caller annotation
// (as enabled by the AddCaller option). When building wrappers around the
// Logger and SugaredLogger, supplying this Option prevents zap from always
// reporting the wrapper code as the caller.
func (l *ZapLogger) AddCallerSkip(skip int) Logger {
	lc := l.clone()
	lc.Logger = lc.Logger.WithOptions(zap.AddCallerSkip(skip))
	return lc
}

// clone 深度拷贝 zapLogger.
func (l *ZapLogger) clone() *ZapLogger {
	copied := *l
	return &copied
}

func (l *ZapLogger) toZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, arg := range fields {
		zapFields = append(zapFields, zap.Any(arg.Key, arg.String))
	}
	return zapFields
}

func (l *ZapLogger) W(ctx context.Context) Logger {
	contextExtractors := map[string]func(context.Context) string{
		"x-request-id": RequestIDFromContext, // 从 context 中提取请求 ID
		"x-trace-id":   TraceIDFromContext,   // 从 context 中提取跟踪 ID
		"x-span-id":    SpanIDFromContext,    // 从 context 中提取跨度 ID
	}

	return l.WC(ctx, contextExtractors)
}

// WC 解析传入的 context，尝试提取关注的键值，并添加到 zap.Logger 结构化日志中.
func (l *ZapLogger) WC(ctx context.Context, contextExtractors map[string]func(context.Context) string) Logger {
	lc := l.clone()

	// 遍历映射，从 context 中提取值并添加到日志中。
	for k, extractor := range contextExtractors {
		if v := extractor(ctx); v != "" {
			lc.Logger = lc.Logger.With(zap.String(k, v))
		}
	}

	return lc
}
