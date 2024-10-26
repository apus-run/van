package memory_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apus-run/van/storage"
	"github.com/apus-run/van/storage/memory"
)

func TestCache(t *testing.T) {
	ctx := context.Background()
	key := "john"
	val := "doe"
	t.Run("CacheGetMiss", func(t *testing.T) {
		if _, err := memory.New().Get(ctx, key); err == nil {
			t.Error("expected to get no value from cache")
		}
	})

	t.Run("CacheGetHit", func(t *testing.T) {
		c := memory.New()

		if err := c.Set(ctx, key, val, 0); err != nil {
			t.Error(err)
		}

		if a, err := c.Get(ctx, key); err != nil {
			t.Errorf("Expected a value, got err: %s", err)
		} else if a != val {
			t.Errorf("Expected '%v', got '%v'", val, a)
		}
	})

	t.Run("CacheGetExpired", func(t *testing.T) {
		c := memory.New()
		e := 20 * time.Millisecond

		if err := c.Set(ctx, key, val, e); err != nil {
			t.Error(err)
		}

		<-time.After(25 * time.Millisecond)
		r, err := c.Get(ctx, key)

		t.Logf("值为: %v", r)
		require.Error(t, err)
	})

	t.Run("CacheGetValid", func(t *testing.T) {
		c := memory.New()
		e := 25 * time.Millisecond

		if err := c.Set(ctx, key, val, e); err != nil {
			t.Error(err)
		}

		<-time.After(20 * time.Millisecond)
		r, err := c.Get(ctx, key)

		t.Logf("值为: %v", r)
		require.Error(t, err)
	})

	t.Run("CacheDeleteMiss", func(t *testing.T) {
		err := memory.New().Delete(ctx, key)
		require.NoError(t, err)
	})

	t.Run("CacheDeleteHit", func(t *testing.T) {
		c := memory.New()

		if err := c.Set(ctx, key, val, 0); err != nil {
			t.Error(err)
		}

		if err := c.Delete(ctx, key); err != nil {
			t.Errorf("Expected to delete an item, got err: %s", err)
		}

		if _, err := c.Get(ctx, key); err == nil {
			t.Errorf("Expected error")
		}
	})
}
func TestCacheWithOptions(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	key := "john"
	val := "doe"
	t.Run("CacheWithExpiration", func(t *testing.T) {
		c := memory.New(memory.WithGCInterval(20 * time.Millisecond))

		if err := c.Set(ctx, key, val, 0); err != nil {
			t.Error(err)
		}

		<-time.After(25 * time.Millisecond)
		r, err := c.Get(ctx, key)
		require.Equal(t, val, r)
		require.NoError(t, err)

	})

	t.Run("CacheWithItems", func(t *testing.T) {
		c := memory.New(memory.Data(map[string]memory.Item{key: {val, 0}}))

		if a, err := c.Get(ctx, key); err != nil {
			t.Errorf("Expected a value, got err: %s", err)
		} else if a != val {
			t.Errorf("Expected '%v', got '%v'", val, a)
		}
	})
}

func Test_Memory(t *testing.T) {
	t.Parallel()
	store := memory.New()
	var (
		key     = "john-internal"
		val any = "doe"
		exp     = 1 * time.Second
		ctx     = context.Background()
	)

	// Set key with value
	err := store.Set(ctx, key, val, 0)
	require.NoError(t, err)
	result, err := store.Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, val, result)

	// Get non-existing key
	result, err = store.Get(ctx, "empty")
	require.Error(t, err)
	require.Nil(t, result)

	// Set key with value and ttl
	err = store.Set(ctx, key, val, exp)
	require.NoError(t, err)
	time.Sleep(1100 * time.Millisecond)
	result, err = store.Get(ctx, key)
	require.Error(t, err)
	require.Nil(t, result)

	// Set key with value and no expiration
	err = store.Set(ctx, key, val, 0)
	require.NoError(t, err)
	result, err = store.Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, val, result)

	// Delete key
	err = store.Delete(ctx, key)
	require.NoError(t, err)
	result, err = store.Get(ctx, key)
	require.Error(t, err)
	require.Nil(t, result)

	// Reset all keys
	err = store.Set(ctx, "john-reset", val, 0)
	require.NoError(t, err)
	err = store.Set(ctx, "doe-reset", val, 0)
	require.NoError(t, err)
	err = store.Flush(ctx)
	require.NoError(t, err)

	// Check if all keys are deleted
	result, err = store.Get(ctx, "john-reset")
	require.Error(t, err)
	require.Nil(t, result)
	result, err = store.Get(ctx, "doe-reset")
	require.Error(t, err)
	require.Nil(t, result)
}

func Benchmark_Memory(b *testing.B) {
	ctx := context.Background()
	keyLength := 1000
	keys := make([]string, keyLength)
	for i := 0; i < keyLength; i++ {
		keys[i] = uuid.New().String()
	}
	value := "joe"

	ttl := 2 * time.Second
	b.Run("fiber_memory", func(b *testing.B) {
		d := memory.New()
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			for _, key := range keys {
				_ = d.Set(ctx, key, value, ttl)

			}
			for _, key := range keys {
				_, _ = d.Get(ctx, key)
			}
			for _, key := range keys {
				_ = d.Delete(ctx, key)

			}
		}
	})
}

func Test_Storage_Memory_Set(t *testing.T) {

	t.Parallel()
	var (
		testStore        = memory.New()
		key       string = "john"
		val       any    = "hello"
		ctx              = context.Background()
	)

	err := testStore.Set(ctx, key, val, 0)
	require.NoError(t, err)

	keys := testStore.Keys(ctx)
	require.Len(t, keys, 1)
}

func Test_Storage_Memory_Del(t *testing.T) {

	t.Parallel()
	var (
		testStore        = memory.New()
		key       string = "john"
		ctx              = context.Background()
	)

	err := testStore.Delete(ctx, key)
	require.NoError(t, err)
}

func Test_Storage_Memory_Set_Override(t *testing.T) {
	t.Parallel()
	var (
		testStore        = memory.New()
		key       string = "john"
		val       any    = "hello"
		ctx              = context.Background()
	)

	err := testStore.Set(ctx, key, val, 0)
	require.NoError(t, err)

	err = testStore.Set(ctx, key, val, 0)
	require.NoError(t, err)

	keys := testStore.Keys(ctx)

	require.Len(t, keys, 1)
}

func Test_Storage_Memory_Get(t *testing.T) {
	t.Parallel()
	var (
		testStore        = memory.New()
		key       string = "john"
		val       any    = "hello"
		ctx              = context.Background()
	)

	err := testStore.Set(ctx, key, val, 0)
	require.NoError(t, err)

	result, err := testStore.Get(ctx, key)
	t.Logf("值为: %v", result)
	require.NoError(t, err)
	require.Equal(t, val, result)

	keys := testStore.Keys(ctx)
	require.NoError(t, err)
	require.Len(t, keys, 1)
}

func Test_Storage_Memory_Set_Expiration(t *testing.T) {
	t.Parallel()
	var (
		testStore     = memory.New(memory.WithGCInterval(300 * time.Millisecond))
		key           = "john"
		val       any = "hello"
		exp           = 1 * time.Second
		ctx           = context.Background()
	)

	err := testStore.Set(ctx, key, val, exp)
	require.NoError(t, err)

	// interval + expire + buffer
	time.Sleep(1500 * time.Millisecond)

	result, err := testStore.Get(ctx, key)
	t.Logf("错误为: %v", err)
	t.Logf("值为: %v", result)
	require.Error(t, err)
	require.Nil(t, result)

	keys := testStore.Keys(ctx)
	t.Logf("值为: %v", len(keys))
	require.Nil(t, keys)
}

func Test_Storage_Memory_Set_Long_Expiration_with_Keys(t *testing.T) {
	t.Parallel()
	var (
		testStore     = memory.New()
		key           = "john"
		val       any = "hello"
		exp           = 3 * time.Second
		ctx           = context.Background()
	)

	keys := testStore.Keys(ctx)
	require.Nil(t, keys)

	err := testStore.Set(ctx, key, val, exp)
	require.NoError(t, err)

	time.Sleep(1100 * time.Millisecond)

	keys = testStore.Keys(ctx)
	require.Len(t, keys, 1)

	time.Sleep(4000 * time.Millisecond)
	result, err := testStore.Get(ctx, key)
	t.Logf("错误为: %v", err)
	t.Logf("值为: %v", result)
	require.Error(t, err)
	require.Nil(t, result)

	keys = testStore.Keys(ctx)
	t.Logf("值为: %v", len(keys))
	require.Nil(t, keys)
}

func Test_Storage_Memory_Get_NotExist(t *testing.T) {
	t.Parallel()
	var (
		testStore = memory.New()
		ctx       = context.Background()
	)
	result, err := testStore.Get(ctx, "notexist")
	t.Logf("错误为: %v", err)
	t.Logf("值为: %v", result)
	require.Error(t, err)
	require.Nil(t, result)

	keys := testStore.Keys(ctx)
	require.Nil(t, keys)
}

func Test_Storage_Memory_Delete(t *testing.T) {
	t.Parallel()
	var (
		testStore     = memory.New()
		key           = "john"
		val       any = "hello"
		ctx           = context.Background()
	)

	err := testStore.Set(ctx, key, val, 0)
	require.NoError(t, err)

	keys := testStore.Keys(ctx)
	require.NoError(t, err)
	require.Len(t, keys, 1)

	err = testStore.Delete(ctx, key)
	require.NoError(t, err)

	result, err := testStore.Get(ctx, key)
	t.Logf("错误为: %v", err)
	t.Logf("值为: %v", result)
	require.Error(t, err)
	require.Nil(t, result)

	keys = testStore.Keys(ctx)
	require.Nil(t, keys)
}

func TestCache_Deletes(t *testing.T) {
	c := memory.New()

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
				require.NoError(t, cache.Set(ctx, "name", "Alex", 0))
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

func Test_Storage_Memory_Reset(t *testing.T) {
	t.Parallel()
	var (
		testStore     = memory.New()
		val       any = "hello"
		ctx           = context.Background()
	)

	err := testStore.Set(ctx, "john1", val, 0)
	require.NoError(t, err)

	err = testStore.Set(ctx, "john2", val, 0)
	require.NoError(t, err)

	keys := testStore.Keys(ctx)
	t.Logf("值为: %v", keys)
	require.Len(t, keys, 2)

	isBoll := testStore.Contains(ctx, keys[0])
	require.True(t, isBoll)

	err = testStore.Flush(ctx)
	require.NoError(t, err)

	result, err := testStore.Get(ctx, "john1")
	require.Error(t, err)
	require.Nil(t, result)

	result, err = testStore.Get(ctx, "john2")
	require.Error(t, err)
	require.Nil(t, result)

	keys = testStore.Keys(ctx)
	require.Nil(t, keys)
}

func Test_Storage_Memory_Close(t *testing.T) {
	t.Parallel()

	var (
		testStore = memory.New()
	)

	err := testStore.Close()
	t.Logf("错误为: %v", err)
	require.NoError(t, err)

}

func Test_Storage_Memory_Conn(t *testing.T) {
	t.Parallel()
	testStore := memory.New()
	require.NotNil(t, testStore.Conn())
}

// Benchmarks for Set operation
func Benchmark_Memory_Set(b *testing.B) {
	var (
		testStore = memory.New()
		ctx       = context.Background()
	)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = testStore.Set(ctx, "john", "doe", 0) //nolint: errcheck // error not needed for benchmark
	}
}

func Benchmark_Memory_Set_Parallel(b *testing.B) {
	var (
		testStore = memory.New()
		ctx       = context.Background()
	)
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = testStore.Set(ctx, "john", "doe", 0) //nolint: errcheck // error not needed for benchmark
		}
	})
}

func Benchmark_Memory_Set_Asserted(b *testing.B) {
	var (
		testStore = memory.New()
		ctx       = context.Background()
	)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := testStore.Set(ctx, "john", "doe", 0)
		require.NoError(b, err)
	}
}

func Benchmark_Memory_Set_Asserted_Parallel(b *testing.B) {
	var (
		testStore = memory.New()
		ctx       = context.Background()
	)
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := testStore.Set(ctx, "john", "doe", 0)
			require.NoError(b, err)
		}
	})
}

// Benchmarks for Get operation
func Benchmark_Memory_Get(b *testing.B) {
	var (
		testStore = memory.New()
		ctx       = context.Background()
	)
	err := testStore.Set(ctx, "john", "doe", 0)
	require.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = testStore.Get(ctx, "john") //nolint: errcheck // error not needed for benchmark
	}
}

func Benchmark_Memory_Get_Parallel(b *testing.B) {
	var (
		testStore = memory.New()
		ctx       = context.Background()
	)
	err := testStore.Set(ctx, "john", "doe", 0)
	require.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = testStore.Get(ctx, "john") //nolint: errcheck // error not needed for benchmark
		}
	})
}

func Benchmark_Memory_Get_Asserted(b *testing.B) {
	var (
		testStore = memory.New()
		ctx       = context.Background()
	)
	err := testStore.Set(ctx, "john", "doe", 0)
	require.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := testStore.Get(ctx, "john")
		require.NoError(b, err)
	}
}

func Benchmark_Memory_Get_Asserted_Parallel(b *testing.B) {
	var (
		testStore = memory.New()
		ctx       = context.Background()
	)
	err := testStore.Set(ctx, "john", "doe", 0)
	require.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := testStore.Get(ctx, "john")
			require.NoError(b, err)
		}
	})
}

// Benchmarks for SetAndDelete operation
func Benchmark_Memory_SetAndDelete(b *testing.B) {
	var (
		testStore = memory.New()
		ctx       = context.Background()
	)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = testStore.Set(ctx, "john", "doe", 0) //nolint: errcheck // error not needed for benchmark
		_ = testStore.Delete(ctx, "john")        //nolint: errcheck // error not needed for benchmark
	}
}

func Benchmark_Memory_SetAndDelete_Parallel(b *testing.B) {
	var (
		testStore = memory.New()
		ctx       = context.Background()
	)
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = testStore.Set(ctx, "john", "doe", 0) //nolint: errcheck // error not needed for benchmark
			_ = testStore.Delete(ctx, "john")        //nolint: errcheck // error not needed for benchmark
		}
	})
}

func Benchmark_Memory_SetAndDelete_Asserted(b *testing.B) {
	var (
		testStore = memory.New()
		ctx       = context.Background()
	)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := testStore.Set(ctx, "john", "doe", 0)
		require.NoError(b, err)

		err = testStore.Delete(ctx, "john")
		require.NoError(b, err)
	}
}

func Benchmark_Memory_SetAndDelete_Asserted_Parallel(b *testing.B) {
	var (
		testStore = memory.New()
		ctx       = context.Background()
	)
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := testStore.Set(ctx, "john", "doe", 3*time.Second)
			require.NoError(b, err)

			err = testStore.Delete(ctx, "john")
			require.NoError(b, err)
		}
	})
}
