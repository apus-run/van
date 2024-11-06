package zlog

import (
	"time"

	"go.uber.org/zap"
)

// TimeValue returns a Value for a time.Time.
// It discards the monotonic portion.
func TimeValue(v time.Time) any {
	return uint64(v.UnixNano())
}

// DurationValue returns a Value for a time.Duration.
func DurationValue(v time.Duration) uint64 {
	return uint64(v.Nanoseconds())
}

func Err(err error) Field {
	return zap.Error(err)
}

func String(key, v string) Field {
	return zap.String(key, v)
}

func Uint64(key string, v uint64) Field {
	return zap.Uint64(key, v)
}

func Int64(key string, v int64) Field {
	return zap.Int64(key, v)
}

func Float64(key string, v float64) Field {
	return zap.Float64(key, v)
}

func Bool(key string, b bool) Field {
	return zap.Bool(key, b)
}

func Int(key string, v int) Field {
	return Int64(key, int64(v))
}

func Any(key string, v any) Field {
	return zap.Any(key, v)
}

func Time(key string, v time.Time) Field {
	return zap.Time(key, v)
}

func Duration(key string, v time.Duration) Field {
	return zap.Duration(key, v)
}
