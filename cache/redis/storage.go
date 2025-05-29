package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	storage "github.com/apus-run/van/cache"
	"github.com/apus-run/van/cache/internal/errs"
)

var (
	_ storage.Storage = (*Storage)(nil)
)

type Storage struct {
	client redis.Cmdable
}

func New(client redis.Cmdable) *Storage {
	return &Storage{
		client: client,
	}
}

func (s *Storage) Set(ctx context.Context, key string, val any, exp time.Duration) error {
	return s.client.Set(ctx, key, val, exp).Err()
}

func (s *Storage) Get(ctx context.Context, key string) (any, error) {
	val, err := s.client.Get(ctx, key).Result()
	if err != nil && errors.Is(err, redis.Nil) {
		return nil, errs.ErrKeyNotExist
	}
	return val, err
}

func (s *Storage) GetAny(ctx context.Context, key string) (val storage.Value) {
	val.Value, val.Error = s.client.Get(ctx, key).Result()
	if val.Error != nil && errors.Is(val.Error, redis.Nil) {
		val.Error = errs.ErrKeyNotExist
	}
	return
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

func (s *Storage) Deletes(ctx context.Context, keys ...string) (int64, error) {
	return s.client.Del(ctx, keys...).Result()
}

func (s *Storage) Flush(ctx context.Context) error {
	return s.client.FlushDBAsync(ctx).Err()
}
func (s *Storage) Keys(ctx context.Context) []string {
	return s.client.Keys(ctx, "*").Val()
}

func (s *Storage) Contains(ctx context.Context, key string) bool {
	return s.client.Exists(ctx, key).Val() > 0
}

func (s *Storage) String() string {
	return "redis"
}
