package memory

import (
	"context"
	"time"

	"github.com/apus-run/van/cache/internal/timer"
)

// Option is config option.
type Option[K comparable, V any] func(*Options[K, V])

type Options[K comparable, V any] struct {
	// Context should contain all implementation specific options, using context.WithValue.
	Context context.Context

	Data map[K]Item[V]

	// gcInterval 清理过期数据的时间间隔
	GCInterval time.Duration
}

// DefaultOptions .
func DefaultOptions[K comparable, V any]() *Options[K, V] {
	return &Options[K, V]{
		Data:       make(map[K]Item[V]),
		GCInterval: 10 * time.Second,
	}
}

func Apply[K comparable, V any](opts ...Option[K, V]) *Options[K, V] {
	options := DefaultOptions[K, V]()
	for _, o := range opts {
		o(options)
	}
	return options
}

func WithGCInterval[K comparable, V any](d time.Duration) Option[K, V] {
	return func(s *Options[K, V]) {
		s.GCInterval = d
	}
}

// Data initializes the cache with preconfigured items.
func Data[K comparable, V any](items map[K]Item[V]) Option[K, V] {
	return func(o *Options[K, V]) {
		o.Data = items
	}
}

// WithContext sets the cache context, for any extra configuration.
func WithContext[K comparable, V any](ctx context.Context) Option[K, V] {
	return func(o *Options[K, V]) {
		o.Context = ctx
	}
}

// WithOptions 设置所有配置
func WithOptions[K comparable, V any](fn func(options *Options[K, V])) Option[K, V] {
	return func(options *Options[K, V]) {
		fn(options)
	}
}

type Item[V any] struct {
	Val V
	Exp int64
}

func NewItem[V any](val V, exp int64) *Item[V] {
	return &Item[V]{Val: val, Exp: exp}
}

func (i *Item[V]) Expired() bool {
	return i.Exp != 0 && i.Exp <= timer.Timestamp()
}

func (i *Item[V]) Value() V {
	return i.Val
}
