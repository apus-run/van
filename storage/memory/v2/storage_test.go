package lru_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	memory "github.com/apus-run/van/storage/lru"
)

func Test_Memory(t *testing.T) {
	t.Parallel()
	store := memory.New[string, string]()
	var (
		key = "john"
		val = "doe"
		exp = 1 * time.Second
		ctx = context.Background()
	)
	require.NotNil(t, store)

	// Set key with value
	err := store.Set(ctx, key, val, exp)
	require.NoError(t, err)
	//
	//// Get key
	//require.NoError(t, err)
	//result, err := store.Get(ctx, key)
	//require.NoError(t, err)
	//require.Equal(t, val, result)
	//
	//// Get non-existing key
	//result, err = store.Get(ctx, "empty")
	//require.Error(t, err)
	//require.Nil(t, result)
	//
	//// Set key with value and ttl
	//err = store.Set(ctx, key, val, exp)
	//require.NoError(t, err)
	//time.Sleep(1100 * time.Millisecond)
	//result, err = store.Get(ctx, key)
	//require.Error(t, err)
	//require.Nil(t, result)
	//
	//// Set key with value and no expiration
	//err = store.Set(ctx, key, val, 0)
	//require.NoError(t, err)
	//result, err = store.Get(ctx, key)
	//require.NoError(t, err)
	//require.Equal(t, val, result)
	//
	//// Delete key
	//err = store.Delete(ctx, key)
	//require.NoError(t, err)
	//result, err = store.Get(ctx, key)
	//require.Error(t, err)
	//require.Nil(t, result)
	//
	//// Reset all keys
	//err = store.Set(ctx, "john-reset", val, 0)
	//require.NoError(t, err)
	//err = store.Set(ctx, "doe-reset", val, 0)
	//require.NoError(t, err)
	//err = store.Flush(ctx)
	//require.NoError(t, err)
	//
	//// Check if all keys are deleted
	//result, err = store.Get(ctx, "john-reset")
	//require.Error(t, err)
	//require.Nil(t, result)
	//result, err = store.Get(ctx, "doe-reset")
	//require.Error(t, err)
	//require.Nil(t, result)
}

func Benchmark_Memory(b *testing.B) {
	//ctx := context.Background()
	//keyLength := 1000
	//keys := make([]string, keyLength)
	//for i := 0; i < keyLength; i++ {
	//	keys[i] = uuid.New().String()
	//}
	//value := "joe"
	//
	//ttl := 2 * time.Second
	//b.Run("fiber_memory", func(b *testing.B) {
	//	d := memory.New[string, string]()
	//	b.ReportAllocs()
	//	b.ResetTimer()
	//	for n := 0; n < b.N; n++ {
	//		for _, key := range keys {
	//			_ = d.Set(ctx, key, value, ttl)
	//
	//		}
	//		for _, key := range keys {
	//			_, _ = d.Get(ctx, key)
	//		}
	//		for _, key := range keys {
	//			_ = d.Delete(ctx, key)
	//
	//		}
	//	}
	//})
}
