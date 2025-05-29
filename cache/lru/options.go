package lru

import (
	"context"
	"time"

	"github.com/apus-run/van/cache/internal/timer"
	lru "github.com/hashicorp/golang-lru/v2"
)

// Option is config option.
type Option func(*Options)

type Options struct {
	// Context should contain all implementation specific options, using context.WithValue.
	Context context.Context

	Data *lru.Cache[string, *Item]

	// gcInterval 清理过期数据的时间间隔
	GCInterval time.Duration

	Size int
}

// DefaultOptions .
func DefaultOptions() *Options {
	cache, _ := lru.New[string, *Item](1000)

	return &Options{
		Data:       cache,
		GCInterval: 10 * time.Second,
	}
}

func Apply(opts ...Option) *Options {
	options := DefaultOptions()
	for _, o := range opts {
		o(options)
	}
	return options
}

func WithGCInterval(d time.Duration) Option {
	return func(s *Options) {
		s.GCInterval = d
	}
}

func WithSize(size int) Option {
	return func(s *Options) {
		s.Size = size
	}
}

// Data initializes the cache with preconfigured items.
func Data(items map[string]Item) Option {
	return func(o *Options) {
		cache, _ := lru.New[string, *Item](len(items))
		for k, v := range items {
			cache.Add(k, &v)
		}
		o.Data = cache
	}
}

// WithContext sets the cache context, for any extra configuration.
func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// WithOptions 设置所有配置
func WithOptions(fn func(options *Options)) Option {
	return func(options *Options) {
		fn(options)
	}
}

type Item struct {
	Val any
	Exp int64
}

func NewItem(val any, exp int64) *Item {
	return &Item{Val: val, Exp: exp}
}

// Expired returns true if the item has expired.
func (i Item) Expired() bool {
	return i.Exp != 0 && i.Exp <= timer.Timestamp()
}

func (i Item) IsExpired(currentTS int64) bool {
	return i.Exp != 0 && i.Exp < currentTS
}

func (i Item) Value() any {
	return i.Val
}
