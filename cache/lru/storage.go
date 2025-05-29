package lru

import (
	"context"
	"sync"
	"time"

	storage "github.com/apus-run/van/cache"
	"github.com/apus-run/van/cache/internal/errs"
	"github.com/apus-run/van/cache/internal/timer"

	lru "github.com/hashicorp/golang-lru/v2"
)

var (
	_ storage.Storage = (*Storage)(nil)
)

// 实现Storage接口的LRU缓存
type Storage struct {
	data       *lru.Cache[string, *Item]
	mux        sync.RWMutex
	gcInterval time.Duration // 后台清理间隔
	done       chan struct{} // 停止信号
}

// 创建带过期机制的LRU缓存
func New(opts ...Option) *Storage {
	options := Apply(opts...)

	s := &Storage{
		gcInterval: options.GCInterval,
		data:       options.Data,
		done:       make(chan struct{}),
	}

	// Start garbage collector
	timer.StartTimeStampUpdater()
	s.gc()

	return s
}

// Set 设置缓存（带过期时间）
func (s *Storage) Set(ctx context.Context, key string, val any, exp time.Duration) error {
	if len(key) == 0 || val == nil {
		return nil
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	var e int64
	if exp > 0 {
		e = timer.Timestamp() + int64(exp.Seconds())
	}

	item := &Item{
		Val: val,
		Exp: e,
	}

	s.data.Add(key, item)
	return nil
}

// Get 获取缓存（带过期检查）
func (s *Storage) Get(ctx context.Context, key string) (any, error) {
	s.mux.RLock()
	item, ok := s.data.Get(key)
	s.mux.RUnlock()

	if !ok {
		return nil, errs.ErrKeyNotExist
	}

	if item.Expired() {
		go s.Delete(context.Background(), key) // 异步清理
		return nil, errs.ErrItemExpired
	}
	return item.Value(), nil
}

// GetAny 获取封装后的值
func (s *Storage) GetAny(ctx context.Context, key string) (val storage.Value) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	if len(key) == 0 {
		val.Error = errs.ErrKeyNotExist
	}

	value, err := s.Get(ctx, key)
	if err != nil {
		val.Error = errs.ErrKeyNotExist
	} else {
		val.Value = value
	}

	return
}

// Delete 删除指定键
func (s *Storage) Delete(ctx context.Context, key string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, ok := s.data.Peek(key); !ok {
		return errs.ErrKeyNotExist
	}

	if s.data.Remove(key) {
		return nil
	}
	return errs.ErrKeyNotExist
}

// Deletes 批量删除
func (s *Storage) Deletes(ctx context.Context, keys ...string) (int64, error) {
	if len(keys) == 0 {
		return 0, nil
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	var n int64
	for _, key := range keys {
		if s.data.Remove(key) {
			n++
		}
	}
	return n, nil
}

// Flush 清空缓存
func (s *Storage) Flush(ctx context.Context) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.data.Purge()
	return nil
}

// Keys 获取有效键列表
func (s *Storage) Keys(ctx context.Context) []string {
	s.mux.RLock()
	defer s.mux.RUnlock()

	keys := s.data.Keys()
	result := make([]string, 0, len(keys))
	for _, k := range keys {
		if item, ok := s.data.Peek(k); ok && !item.Expired() {
			result = append(result, k)
		}
	}
	return result
}

// Contains 检查键是否存在且有效
func (s *Storage) Contains(ctx context.Context, key string) bool {
	s.mux.RLock()
	defer s.mux.RUnlock()

	if item, ok := s.data.Peek(key); ok {
		return !item.Expired()
	}
	return false
}

// String 缓存状态描述
func (s *Storage) String() string {
	return "lru"
}

func (s *Storage) gc() {
	go func() {
		ticker := time.NewTicker(s.gcInterval)
		defer ticker.Stop()
		// 内存预分配
		expired := make([]string, 0, 100)

		for {
			select {
			case <-s.done:
				return
			case <-ticker.C:
				ts := timer.Timestamp()
				expired = expired[:0]

				// 锁定以读取数据
				s.mux.RLock()
				keys := s.data.Keys()
				for _, key := range keys {
					if item, ok := s.data.Peek(key); ok && item.IsExpired(ts) {
						expired = append(expired, key)
					}
				}
				s.mux.RUnlock()

				// 锁定以删除过期项
				if len(expired) > 0 {
					s.mux.Lock()
					for _, key := range expired {
						if item, ok := s.data.Peek(key); ok && item.IsExpired(ts) {
							s.data.Remove(key)
						}
					}
					s.mux.Unlock()
				}
			}
		}
	}()
}

// 关闭缓存释放资源
func (s *Storage) Close() {
	close(s.done)
	s.Flush(context.Background())
}

// Conn 返回缓存数据
func (s *Storage) Conn() map[string]Item {
	s.mux.RLock()
	defer s.mux.RUnlock()

	data := make(map[string]Item)
	ts := timer.Timestamp()

	for _, key := range s.data.Keys() {
		if item, ok := s.data.Peek(key); ok && !item.IsExpired(ts) {
			data[key] = *item
		}
	}
	return data
}
