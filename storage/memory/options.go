package memory

import (
	"context"
	"time"

	"github.com/apus-run/van/storage/internal/timer"
)

// Option is config option.
type Option func(*Options)

type Options struct {
	// Context should contain all implementation specific options, using context.WithValue.
	Context context.Context

	Data map[string]Item

	// gcInterval 清理过期数据的时间间隔
	GCInterval time.Duration
}

// DefaultOptions .
func DefaultOptions() *Options {
	return &Options{
		Data:       make(map[string]Item),
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

// Data initializes the cache with preconfigured items.
func Data(items map[string]Item) Option {
	return func(o *Options) {
		o.Data = items
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
func (i Item) Value() any {
	return i.Val
}
