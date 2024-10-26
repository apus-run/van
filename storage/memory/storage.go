package memory

import (
	"context"
	"sync"
	"time"

	"github.com/apus-run/van/storage"
	"github.com/apus-run/van/storage/internal/errs"
	"github.com/apus-run/van/storage/internal/timer"
)

var (
	_ storage.Storage = (*Storage)(nil)
)

// Storage defines a concurrent safe in memory key-value data store.
type Storage struct {
	data       map[string]Item
	done       chan struct{}
	gcInterval time.Duration
	mux        sync.RWMutex
}

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

func (s *Storage) Get(ctx context.Context, key string) (any, error) {
	if len(key) == 0 {
		return nil, errs.ErrKeyNotExist
	}
	s.mux.RLock()
	defer s.mux.RUnlock()

	item, ok := s.data[key]
	if !ok {
		return nil, errs.ErrKeyNotExist
	}
	if item.Expired() {
		return nil, errs.ErrItemExpired
	}

	return item.Value(), nil
}

func (s *Storage) GetAny(ctx context.Context, key string) (val storage.Value) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	if len(key) == 0 {
		val.Error = errs.ErrKeyNotExist
	}

	var ok bool
	val.Value, ok = s.data[key]
	if !ok {
		val.Error = errs.ErrKeyNotExist
	}
	if val.Value.(Item).Expired() {
		val.Error = errs.ErrItemExpired
	}

	return
}

func (s *Storage) Set(ctx context.Context, key string, val any, exp time.Duration) error {
	if len(key) == 0 || val == nil {
		return nil
	}

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

func (s *Storage) Delete(ctx context.Context, key string) error {
	if len(key) == 0 {
		return errs.ErrKeyNotExist
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	delete(s.data, key)

	return nil
}

func (s *Storage) Deletes(ctx context.Context, keys ...string) (int64, error) {
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

func (s *Storage) Flush(ctx context.Context) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.data = make(map[string]Item)

	return nil
}

func (s *Storage) Keys(ctx context.Context) []string {
	s.mux.RLock()
	defer s.mux.RUnlock()

	if len(s.data) == 0 {
		return nil
	}

	ts := timer.Timestamp()
	keys := make([]string, 0, len(s.data))

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

func (s *Storage) Contains(ctx context.Context, key string) bool {
	if len(key) == 0 {
		return false
	}

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
func (s *Storage) Close() error {
	s.done <- struct{}{}

	return nil
}

// Conn return database client.
func (s *Storage) Conn() map[string]Item {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.data
}

func (s *Storage) String() string {
	return "memory"
}

func (s *Storage) gc() {
	go func() {
		ticker := time.NewTicker(s.gcInterval)
		defer ticker.Stop()
		var expired []string

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
