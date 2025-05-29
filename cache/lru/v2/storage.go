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

// 实现Storage接口的LRU缓存（泛型版本）
type Storage[K comparable, V any] struct {
	data       *lru.Cache[K, *Item[V]]
	mux        sync.RWMutex
	gcInterval time.Duration
	done       chan struct{}
}

// 创建带过期机制的LRU缓存（泛型）
func New[K comparable, V any](opts ...Option[K, V]) *Storage[K, V] {
	options := Apply(opts...)

	s := &Storage[K, V]{
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
func (s *Storage[K, V]) Set(ctx context.Context, key K, val V, exp time.Duration) error {
	var zeroK K
	if key == zeroK {
		return nil
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	var e int64
	if exp > 0 {
		e = timer.Timestamp() + int64(exp.Seconds())
	}

	item := &Item[V]{
		Val: val,
		Exp: e,
	}

	s.data.Add(key, item)
	return nil
}

// Get 获取缓存（带过期检查）
func (s *Storage[K, V]) Get(ctx context.Context, key K) (V, error) {
	var zeroV V
	var zeroK K
	if key == zeroK {
		return zeroV, errs.ErrKeyNotExist
	}

	s.mux.RLock()
	item, ok := s.data.Get(key)
	s.mux.RUnlock()

	if !ok {
		return zeroV, errs.ErrKeyNotExist
	}

	if item.Expired() {
		go s.Delete(context.Background(), key)
		return zeroV, errs.ErrItemExpired
	}
	return item.Value(), nil
}

// GetAny 获取封装后的值
func (s *Storage[K, V]) GetAny(ctx context.Context, key K) (val storage.Value) {
	var zeroK K
	if key == zeroK {
		val.Error = errs.ErrKeyNotExist
		return
	}

	s.mux.RLock()
	defer s.mux.RUnlock()

	value, err := s.Get(ctx, key)
	if err != nil {
		val.Error = err
	} else {
		val.Value = value
	}
	return
}

// Delete 删除指定键
func (s *Storage[K, V]) Delete(ctx context.Context, key K) error {
	var zeroK K
	if key == zeroK {
		return errs.ErrKeyNotExist
	}

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
func (s *Storage[K, V]) Deletes(ctx context.Context, keys ...K) (int64, error) {
	if len(keys) == 0 {
		return 0, nil
	}

	var count int64
	s.mux.Lock()
	defer s.mux.Unlock()

	for _, key := range keys {
		if _, ok := s.data.Peek(key); ok {
			s.data.Remove(key)
			count++
		}
	}
	return count, nil
}

// Flush 清空缓存
func (s *Storage[K, V]) Flush(ctx context.Context) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.data.Purge()
	return nil
}

// Keys 获取有效键列表
func (s *Storage[K, V]) Keys(ctx context.Context) []K {
	s.mux.RLock()
	defer s.mux.RUnlock()

	keys := s.data.Keys()
	result := make([]K, 0, len(keys))
	for _, k := range keys {
		if item, ok := s.data.Peek(k); ok && !item.Expired() {
			result = append(result, k)
		}
	}
	return result
}

// Contains 检查键是否存在且有效
func (s *Storage[K, V]) Contains(ctx context.Context, key K) bool {
	s.mux.RLock()
	defer s.mux.RUnlock()

	if item, ok := s.data.Peek(key); ok {
		return !item.Expired()
	}
	return false
}

// String 缓存状态描述
func (s *Storage[K, V]) String() string {
	return "lru"
}

func (s *Storage[K, V]) gc() {
	go func() {
		ticker := time.NewTicker(s.gcInterval)
		defer ticker.Stop()
		// 内存预分配
		expired := make([]K, 0, 100)

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
func (s *Storage[K, V]) Close() {
	close(s.done)
	s.Flush(context.Background())
}

// Conn 返回缓存数据
func (s *Storage[K, V]) Conn() map[K]Item[V] {
	s.mux.RLock()
	defer s.mux.RUnlock()

	data := make(map[K]Item[V])
	ts := timer.Timestamp()

	for _, key := range s.data.Keys() {
		if item, ok := s.data.Peek(key); ok && !item.IsExpired(ts) {
			data[key] = *item
		}
	}
	return data
}
