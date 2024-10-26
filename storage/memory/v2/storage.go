package memory

import (
	"context"
	"sync"
	"time"

	"github.com/apus-run/van/storage"
	"github.com/apus-run/van/storage/internal/errs"
	"github.com/apus-run/van/storage/internal/timer"
)

// Storage defines a concurrent safe in memory key-value data store.
type Storage[K comparable, V any] struct {
	data       map[K]Item[V]
	done       chan struct{}
	gcInterval time.Duration
	mux        sync.RWMutex
}

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

func (s *Storage[K, V]) Get(ctx context.Context, key K) (V, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	item, ok := s.data[key]
	if !ok {
		return item.Value(), errs.ErrKeyNotExist
	}
	if item.Expired() {
		return item.Value(), errs.ErrItemExpired
	}

	return item.Value(), nil
}

func (s *Storage[K, V]) GetAny(ctx context.Context, key K) (val storage.Value) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	var ok bool
	val.Value, ok = s.data[key]
	if !ok {
		val.Error = errs.ErrKeyNotExist
	}
	if val.Value.(*Item[V]).Expired() {
		val.Error = errs.ErrItemExpired
	}

	return
}

func (s *Storage[K, V]) Set(ctx context.Context, key K, val V, exp time.Duration) error {
	var e int64

	if exp > 0 {
		e = timer.Timestamp() + int64(exp.Seconds())
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	item := NewItem(val, e)
	s.data[key] = *item

	return nil
}

func (s *Storage[K, V]) Delete(ctx context.Context, key K) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	delete(s.data, key)

	return nil
}

func (s *Storage[K, V]) Deletes(ctx context.Context, keys ...K) (int64, error) {
	if len(keys) == 0 {
		return 0, nil
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	n := int64(0)
	for _, k := range keys {
		if _, ok := s.data[k]; ok {
			delete(s.data, k)
			n++
		}
	}

	return n, nil
}

func (s *Storage[K, V]) Len(ctx context.Context) int {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return len(s.data)
}

func (s *Storage[K, V]) Flush(ctx context.Context) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.data = make(map[K]Item[V])

	return nil
}

func (s *Storage[K, V]) Keys(ctx context.Context) []K {
	s.mux.RLock()
	defer s.mux.RUnlock()

	if len(s.data) == 0 {
		return nil
	}

	ts := timer.Timestamp()
	keys := make([]K, 0, len(s.data))

	for k, v := range s.data {
		// Filter out the expired keys
		if v.Exp == 0 || v.Exp > ts {
			keys = append(keys, k)
		}
	}

	if len(keys) == 0 {
		return nil
	}

	return keys
}

func (s *Storage[K, V]) Contains(ctx context.Context, key K) bool {
	s.mux.RLock()
	defer s.mux.RUnlock()

	v, ok := s.data[key]
	if ok {
		if v.Expired() {
			delete(s.data, key)

			return false
		}
	}

	return ok
}

// Close the memory storage.
func (s *Storage[K, V]) Close() error {
	s.done <- struct{}{}

	return nil
}

// Conn return database client.
func (s *Storage[K, V]) Conn() map[K]Item[V] {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.data
}

func (s *Storage[K, V]) String() string {
	return "memory"
}

func (s *Storage[K, V]) gc() {
	go func() {
		ticker := time.NewTicker(s.gcInterval)
		defer ticker.Stop()
		var expired []K

		for {
			select {
			case <-s.done:
				return
			case <-ticker.C:
				ts := timer.Timestamp()
				expired = expired[:0]

				// 锁定以读取数据
				s.mux.RLock()
				for id, v := range s.data {
					if v.Exp != 0 && v.Exp < ts {
						expired = append(expired, id)
					}
				}
				s.mux.RUnlock()

				// 锁定以删除过期项
				s.mux.Lock()
				for _, id := range expired {
					if v, ok := s.data[id]; ok && v.Exp <= ts {
						delete(s.data, id)
					}
				}
				s.mux.Unlock()
			}
		}
	}()
}
