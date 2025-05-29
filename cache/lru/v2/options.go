package lru

import (
	"context"
	"time"

	"github.com/apus-run/van/cache/internal/timer"
	lru "github.com/hashicorp/golang-lru/v2"
)

// Option 泛型配置项
type Option[K comparable, V any] func(*Options[K, V])

type Options[K comparable, V any] struct {
	Context    context.Context
	Data       *lru.Cache[K, *Item[V]]
	GCInterval time.Duration
	Size       int
}

// DefaultOptions 默认泛型配置
func DefaultOptions[K comparable, V any]() *Options[K, V] {
	cache, _ := lru.New[K, *Item[V]](1000)

	return &Options[K, V]{
		Data:       cache,
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

// WithGCInterval 泛型版本
func WithGCInterval[K comparable, V any](d time.Duration) Option[K, V] {
	return func(s *Options[K, V]) {
		s.GCInterval = d
	}
}

// WithSize 泛型版本
func WithSize[K comparable, V any](size int) Option[K, V] {
	return func(s *Options[K, V]) {
		s.Size = size
		if s.Data == nil {
			cache, _ := lru.New[K, *Item[V]](size)
			s.Data = cache
		}
	}
}

// Data 泛型数据初始化
func Data[K comparable, V any](items map[K]Item[V]) Option[K, V] {
	return func(o *Options[K, V]) {
		cache, _ := lru.New[K, *Item[V]](len(items))
		for k, v := range items {
			cache.Add(k, &v)
		}
		o.Data = cache
	}
}

// WithContext 保持原样
func WithContext[K comparable, V any](ctx context.Context) Option[K, V] {
	return func(o *Options[K, V]) {
		o.Context = ctx
	}
}

// WithOptions 泛型版本
func WithOptions[K comparable, V any](fn func(options *Options[K, V])) Option[K, V] {
	return func(options *Options[K, V]) {
		fn(options)
	}
}

// Item 泛型存储项
type Item[V any] struct {
	Val V
	Exp int64
}

func NewItem[V any](val V, exp int64) *Item[V] {
	return &Item[V]{Val: val, Exp: exp}
}

// Expired 保持逻辑不变
func (i Item[V]) Expired() bool {
	return i.Exp != 0 && i.Exp <= timer.Timestamp()
}

func (i Item[V]) IsExpired(currentTS int64) bool {
	return i.Exp != 0 && i.Exp < currentTS
}

func (i Item[V]) Value() V {
	return i.Val
}
