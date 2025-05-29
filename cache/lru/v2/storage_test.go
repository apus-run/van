package lru_test

import (
	"context"
	"testing"
	"time"

	"github.com/apus-run/van/cache/lru/v2"
	"github.com/stretchr/testify/assert"
)

func TestStorageNew(t *testing.T) {
	cache := lru.New[string, string]()
	assert.NotNil(t, cache)
}

func TestStorage_SetAndGet(t *testing.T) {
	cache := lru.New[string, string]()
	ctx := context.Background()

	// Test setting and getting a value
	err := cache.Set(ctx, "key1", "value1", time.Minute)
	assert.NoError(t, err)

	val, err := cache.Get(ctx, "key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)

	// Test getting a non-existent key
	_, err = cache.Get(ctx, "key2")
	assert.Error(t, err)
}

func TestStorage_Expiration(t *testing.T) {
	cache := lru.New[string, string]()
	ctx := context.Background()

	// Set a value with a short expiration
	err := cache.Set(ctx, "key1", "value1", time.Second)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	// Ensure the value has expired
	_, err = cache.Get(ctx, "key1")
	assert.Error(t, err)
}

func TestStorage_Delete(t *testing.T) {
	cache := lru.New[string, string]()
	ctx := context.Background()

	// Set and delete a value
	err := cache.Set(ctx, "key1", "value1", time.Minute)
	assert.NoError(t, err)

	err = cache.Delete(ctx, "key1")
	assert.NoError(t, err)

	// Ensure the value is deleted
	_, err = cache.Get(ctx, "key1")
	assert.Error(t, err)
}

func TestStorage_Deletes(t *testing.T) {
	cache := lru.New[string, string]()
	ctx := context.Background()

	// Set multiple values
	err := cache.Set(ctx, "key1", "value1", time.Minute)
	assert.NoError(t, err)
	err = cache.Set(ctx, "key2", "value2", time.Minute)
	assert.NoError(t, err)

	// Delete multiple values
	count, err := cache.Deletes(ctx, "key1", "key2")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)

	// Ensure the values are deleted
	_, err = cache.Get(ctx, "key1")
	assert.Error(t, err)
	_, err = cache.Get(ctx, "key2")
	assert.Error(t, err)
}

func TestStorage_Flush(t *testing.T) {
	cache := lru.New[string, string]()
	ctx := context.Background()

	// Set multiple values
	err := cache.Set(ctx, "key1", "value1", time.Minute)
	assert.NoError(t, err)
	err = cache.Set(ctx, "key2", "value2", time.Minute)
	assert.NoError(t, err)

	// Flush the cache
	err = cache.Flush(ctx)
	assert.NoError(t, err)

	// Ensure all values are deleted
	_, err = cache.Get(ctx, "key1")
	assert.Error(t, err)
	_, err = cache.Get(ctx, "key2")
	assert.Error(t, err)
}

func TestStorage_Keys(t *testing.T) {
	cache := lru.New[string, string]()
	ctx := context.Background()

	// Set multiple values
	err := cache.Set(ctx, "key1", "value1", time.Minute)
	assert.NoError(t, err)
	err = cache.Set(ctx, "key2", "value2", time.Minute)
	assert.NoError(t, err)

	// Get all keys
	keys := cache.Keys(ctx)
	assert.ElementsMatch(t, []string{"key1", "key2"}, keys)
}

func TestStorage_Contains(t *testing.T) {
	cache := lru.New[string, string]()
	ctx := context.Background()

	// Set a value
	err := cache.Set(ctx, "key1", "value1", time.Minute)
	assert.NoError(t, err)

	// Check if the key exists
	assert.True(t, cache.Contains(ctx, "key1"))

	// Check a non-existent key
	assert.False(t, cache.Contains(ctx, "key2"))
}

func TestStorage_Close(t *testing.T) {
	cache := lru.New[string, string]()
	ctx := context.Background()

	// Set a value
	err := cache.Set(ctx, "key1", "value1", time.Minute)
	assert.NoError(t, err)

	// Close the cache
	cache.Close()

	// Ensure the cache is cleared
	_, err = cache.Get(ctx, "key1")
	assert.Error(t, err)
}
