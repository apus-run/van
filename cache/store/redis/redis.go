package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/apus-run/van/cache/store"
)

const (
	// RedisType represents the storage type as a string value.
	RedisType = "redis"
)

// RedisStore is a store for Redis.
type RedisStore struct {
	client redis.Cmdable
}

// NewRedis creates a new store to Redis instance(s).
func NewRedis(client redis.Cmdable) *RedisStore {
	return &RedisStore{
		client: client,
	}
}

// Get returns data stored from a given key.
func (s *RedisStore) Get(ctx context.Context, key any) (any, error) {
	obj, err := s.client.Get(ctx, key.(string)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, store.ErrKeyNotFound
	}
	return obj, err
}

// GetWithTTL returns data stored from a given key and its corresponding TTL.
func (s *RedisStore) GetWithTTL(ctx context.Context, key any) (any, time.Duration, error) {
	obj, err := s.client.Get(ctx, key.(string)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, 0, store.ErrKeyNotFound
	}
	if err != nil {
		return nil, 0, err
	}

	ttl, err := s.client.TTL(ctx, key.(string)).Result()
	if err != nil {
		return nil, 0, err
	}

	return obj, ttl, err
}

// Set defines data in Redis for given key identifier.
func (s *RedisStore) Set(ctx context.Context, key any, value any) error {
	err := s.client.Set(ctx, key.(string), value, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// Set defines data in Redis for given key identifier.
func (s *RedisStore) SetWithTTL(ctx context.Context, key any, value any, ttl time.Duration) error {
	err := s.client.Set(ctx, key.(string), value, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

// Del removes data from Redis for given key identifier.
func (s *RedisStore) Del(ctx context.Context, key any) error {
	_, err := s.client.Del(ctx, key.(string)).Result()
	return err
}

// Clear resets all data in the store.
func (s *RedisStore) Clear(ctx context.Context) error {
	if err := s.client.FlushAll(ctx).Err(); err != nil {
		return err
	}

	return nil
}

func (s *RedisStore) Wait(ctx context.Context) {
}
