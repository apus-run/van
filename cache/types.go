package cache

import (
	"context"
	"errors"
	"time"

	"github.com/apus-run/van/cache/internal/errs"
	"github.com/apus-run/van/pkg/value"
)

type Storage interface {
	Set(ctx context.Context, key string, val any, exp time.Duration) error
	Get(ctx context.Context, key string) (any, error)
	GetAny(ctx context.Context, key string) Value
	Delete(ctx context.Context, key string) error
	Deletes(ctx context.Context, keys ...string) (int64, error)
	Flush(ctx context.Context) error
	Keys(ctx context.Context) []string
	Contains(ctx context.Context, key string) bool
	String() string
}

// Value 代表一个从缓存中读取出来的值
type Value struct {
	value.AnyValue
}

func (v Value) KeyNotFound() bool {
	return errors.Is(v.Error, errs.ErrKeyNotExist)
}
