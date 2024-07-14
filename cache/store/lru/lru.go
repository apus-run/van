package lru

import (
	"context"
	"fmt"
	"time"

	"github.com/apus-run/van/cache/store"
)

const (
	// RistrettoType represents the storage type as a string value.
	RistrettoType = "golang-lru"
	// RistrettoTagPattern represents the tag pattern to be used as a key in specified storage.
	RistrettoTagPattern = "gocache_tag_%s"
)

// GolangLruClientInterface represents a github.com/hashicorp/golang-lru/v2 client.
type GolangLruClientInterface interface {
	// Adds a value to the cache, returns true if an eviction occurred and
	// updates the "recently used"-ness of the key.
	Add(key, value any) bool

	// Returns key's value from the cache and
	// updates the "recently used"-ness of the key. #value, isFound
	Get(key any) (value any, ok bool)

	// Checks if a key exists in cache without updating the recent-ness.
	Contains(key any) (ok bool)

	// Returns key's value without updating the "recently used"-ness of the key.
	Peek(key any) (value any, ok bool)

	// Removes a key from the cache.
	Remove(key any) bool

	// Returns a slice of the keys in the cache, from oldest to newest.
	Keys() []any

	// Returns the number of items in the cache.
	Len() int

	// Clears all cache entries.
	Purge()
}

// GolangLruStore is a store for GoCache (memory) library.
type GolangLruStore struct {
	client GolangLruClientInterface
}

// NewGoCache creates a new store to GoCache (memory) library instance.
func NewGoCache(client GolangLruClientInterface) *GolangLruStore {
	return &GolangLruStore{
		client: client,
	}
}

// Get returns data stored from a given key.
func (s *GolangLruStore) Get(_ context.Context, key any) (any, error) {
	var err error

	value, ok := s.client.Get(key)
	if !ok {
		err = store.ErrKeyNotFound
	}

	return value, err
}

// GetWithTTL returns data stored from a given key and its corresponding TTL. expiration 无效 由lru 统一控制过期时间
func (s *GolangLruStore) GetWithTTL(ctx context.Context, key any) (any, time.Duration, error) {
	value, err := s.Get(ctx, key)
	return value, 0, err
}

// Set defines data in Ristretto memory cache for given key identifier.
func (s *GolangLruStore) Set(_ context.Context, key any, value any) error {
	if ok := s.client.Add(key, value); !ok {
		return fmt.Errorf("an error has occurred while setting value '%v' on key '%v'", value, key)
	}

	return nil
}

// SetWithTTL ttl 无效 由lru 统一控制过期时间
func (s *GolangLruStore) SetWithTTL(ctx context.Context, key any, value any, ttl time.Duration) error {
	if ok := s.client.Add(key, value); !ok {
		return fmt.Errorf("an error has occurred while setting value '%v' on key '%v'", value, key)
	}

	return nil
}

// Delete removes data in Ristretto memory cache for given key identifier.
func (s *GolangLruStore) Del(_ context.Context, key any) error {
	ok := s.client.Remove(key)
	if !ok {
		return fmt.Errorf("an error has occurred while deleting key '%v'", key)
	}
	return nil
}

// Clear resets all data in the store.
func (s *GolangLruStore) Clear(_ context.Context) error {
	s.client.Purge()
	return nil
}

func (s *GolangLruStore) Wait(_ context.Context) {

}
