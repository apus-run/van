// go:build e2e

package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apus-run/van/storage"
	"github.com/apus-run/van/storage/internal/errs"
	cache "github.com/apus-run/van/storage/redis"
)

func TestCache_e2e_Set(t *testing.T) {
	rdb := newRedisClient()
	require.NoError(t, rdb.Ping(context.Background()).Err())

	testCases := []struct {
		name  string
		after func(ctx context.Context, t *testing.T)

		key        string
		val        string
		expiration time.Duration

		wantErr error
	}{
		{
			name: "set e2e value",
			after: func(ctx context.Context, t *testing.T) {
				result, err := rdb.Get(ctx, "name").Result()
				require.NoError(t, err)
				assert.Equal(t, "foo", result)

				_, err = rdb.Del(ctx, "name").Result()
				require.NoError(t, err)
			},
			key:        "name",
			val:        "foo",
			expiration: time.Minute,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
			defer cancelFunc()
			c := cache.New(rdb)

			err := c.Set(ctx, "name", "foo", time.Minute)
			assert.NoError(t, err)
			tc.after(ctx, t)
		})
	}
}

func TestCache_e2e_Get(t *testing.T) {
	rdb := newRedisClient()
	require.NoError(t, rdb.Ping(context.Background()).Err())

	testCases := []struct {
		name   string
		before func(ctx context.Context, t *testing.T)
		after  func(ctx context.Context, t *testing.T)

		key string

		wantVal string
		wantErr error
	}{
		{
			name: "get e2e value",
			before: func(ctx context.Context, t *testing.T) {
				require.NoError(t, rdb.Set(ctx, "name", "foo", time.Minute).Err())
			},
			after: func(ctx context.Context, t *testing.T) {
				require.NoError(t, rdb.Del(ctx, "name").Err())
			},
			key: "name",

			wantVal: "foo",
		},
		{
			name:    "get e2e error",
			key:     "name",
			before:  func(ctx context.Context, t *testing.T) {},
			after:   func(ctx context.Context, t *testing.T) {},
			wantErr: errs.ErrKeyNotExist,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
			defer cancelFunc()
			c := cache.New(rdb)

			tc.before(ctx, t)
			val := c.GetAny(ctx, tc.key)
			assert.Equal(t, tc.wantErr, val.Error)
			if val.Error != nil {
				return
			}
			assert.Equal(t, tc.wantVal, val.Value.(string))
			tc.after(ctx, t)
		})
	}
}

func TestCache_e2e_Delete(t *testing.T) {
	c, err := newCache()
	require.NoError(t, err)

	testCases := []struct {
		name   string
		before func(ctx context.Context, t *testing.T, cache storage.Storage)

		ctxFunc func() context.Context
		key     []string

		wantN   int64
		wantErr error
	}{
		{
			name: "delete single existed key",
			before: func(ctx context.Context, t *testing.T, cache storage.Storage) {
				require.NoError(t, c.Set(ctx, "name", "Alex", 0))
			},
			ctxFunc: func() context.Context {
				return context.Background()
			},
			key:   []string{"name"},
			wantN: 1,
		},
		{
			name:   "delete single does not existed key",
			before: func(ctx context.Context, t *testing.T, cache storage.Storage) {},
			ctxFunc: func() context.Context {
				return context.Background()
			},
			key: []string{"notExistedKey"},
		},
		{
			name: "delete multiple existed keys",
			before: func(ctx context.Context, t *testing.T, cache storage.Storage) {
				require.NoError(t, cache.Set(ctx, "name", "Alex", 0))
				require.NoError(t, cache.Set(ctx, "age", 18, 0))
			},
			ctxFunc: func() context.Context {
				return context.Background()
			},
			key:   []string{"name", "age"},
			wantN: 2,
		},
		{
			name:   "delete multiple do not existed keys",
			before: func(ctx context.Context, t *testing.T, cache storage.Storage) {},
			ctxFunc: func() context.Context {
				return context.Background()
			},
			key: []string{"name", "age"},
		},
		{
			name: "delete multiple keys, some do not existed keys",
			before: func(ctx context.Context, t *testing.T, cache storage.Storage) {
				require.NoError(t, cache.Set(ctx, "name", "Alex", 0))
				require.NoError(t, cache.Set(ctx, "age", 18, 0))
				require.NoError(t, cache.Set(ctx, "gender", "male", 0))
			},
			ctxFunc: func() context.Context {
				return context.Background()
			},
			key:   []string{"name", "age", "gender", "addr"},
			wantN: 3,
		},
		{
			name:   "timeout",
			before: func(ctx context.Context, t *testing.T, cache storage.Storage) {},
			ctxFunc: func() context.Context {
				timeout := time.Millisecond * 100
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()
				time.Sleep(timeout * 2)
				return ctx
			},
			key:     []string{"name", "age", "addr"},
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tc.ctxFunc()
			tc.before(ctx, t, c)
			n, err := c.Deletes(ctx, tc.key...)
			if err != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			assert.Equal(t, tc.wantN, n)
		})
	}
}

func newCache() (storage.Storage, error) {
	rdb := newRedisClient()
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return cache.New(rdb), nil
}

func newRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
