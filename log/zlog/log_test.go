package zlog

import (
	"context"
	"io"
	"os"
	"testing"

	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

func getLogWriter() io.Writer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "logs.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Compress:   true,
	}
	return lumberJackLogger
}

func TestLog(t *testing.T) {
	wr := getLogWriter()
	logger := NewLogger(WithFormat(FormatJSON), WithWriters(os.Stdout, wr))
	defer logger.Close()

	logger.Info("This is an info message", zap.String("route", "/hello"), zap.Int64("port", 8090))
	logger.Infof("我是日志: %v, %v", zap.String("route", "/hello"), zap.Int64("port", 8090))
	logger.Error("This is an error message")
}

/*
// InitZapLogger 日志
func InitZapLogger(mode string) logger.Logger {
	var cfg zap.Config
	// 这里我们用一个小技巧，
	// 就是直接使用 zap 本身的配置结构体来处理
	if mode == "prod" {
		cfg = zap.NewProductionConfig()
		cfg.InitialFields = map[string]any{"version": "1.0.0"}
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	err := viper.UnmarshalKey("log", &cfg)
	if err != nil {
		panic(err)
	}
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}
*/

// 性能测试用例
func BenchmarkZapLoggerW(b *testing.B) {
	// 创建一个 zapLogger 实例（使用 zap.NewNop() 模拟 logger）
	logger := &ZapLogger{zap.NewNop()}

	// 创建一个包含上下文值的 context
	ctx := WithRequestID(context.Background(), "request-id-12345")
	ctx = WithUserID(ctx, "user-id-67890")

	// 重复调用 W 函数，测量性能
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = logger.W(ctx)
	}
}
