package lru_test

import (
	"context"
	"testing"
	"time"

	"github.com/apus-run/van/cache/internal/errs"
	"github.com/apus-run/van/cache/lru"
	"github.com/stretchr/testify/assert"
)

func TestStorage_New(t *testing.T) {
	store := lru.New()
	defer store.Close()

	assert.NotNil(t, store)
}

func TestStorage_SetWithEmptyKey(t *testing.T) {
	store := lru.New()
	defer store.Close()

	ctx := context.Background()
	value := "test-value"
	expiration := time.Minute

	// Test Set with an empty key
	err := store.Set(ctx, "", value, expiration)
	assert.NoError(t, err)

	// Ensure the key does not exist
	result, err := store.Get(ctx, "")
	assert.ErrorIs(t, err, errs.ErrKeyNotExist)
	assert.Nil(t, result)
}

func TestStorage_SetWithNilValue(t *testing.T) {
	store := lru.New()
	defer store.Close()

	ctx := context.Background()
	key := "test-key"
	expiration := time.Minute

	// Test Set with a nil value
	err := store.Set(ctx, key, nil, expiration)
	assert.NoError(t, err)

	// Ensure the key does not exist
	result, err := store.Get(ctx, key)
	assert.ErrorIs(t, err, errs.ErrKeyNotExist)
	assert.Nil(t, result)
}

func TestStorage_Deletes(t *testing.T) {
	store := lru.New()
	defer store.Close()

	ctx := context.Background()
	key1 := "key1"
	key2 := "key2"
	value := "value"

	// Set multiple keys
	err := store.Set(ctx, key1, value, 0)
	assert.NoError(t, err)
	err = store.Set(ctx, key2, value, 0)
	assert.NoError(t, err)

	// Delete multiple keys
	deletedCount, err := store.Deletes(ctx, key1, key2)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), deletedCount)

	// Ensure keys are deleted
	_, err = store.Get(ctx, key1)
	assert.ErrorIs(t, err, errs.ErrKeyNotExist)
	_, err = store.Get(ctx, key2)
	assert.ErrorIs(t, err, errs.ErrKeyNotExist)
}

func TestStorage_Conn(t *testing.T) {
	store := lru.New()
	defer store.Close()

	ctx := context.Background()
	key := "test-key"
	value := "test-value"

	// Set a value
	err := store.Set(ctx, key, value, 0)
	assert.NoError(t, err)

	// Get the internal connection
	conn := store.Conn()
	assert.Contains(t, conn, key)
	assert.Equal(t, value, conn[key].Value())
}

func TestStorage_Close(t *testing.T) {
	store := lru.New()

	ctx := context.Background()
	key := "test-key"
	value := "test-value"

	// Set a value
	err := store.Set(ctx, key, value, 0)
	assert.NoError(t, err)

	// Close the storage
	store.Close()

	// Ensure the storage is flushed
	result, err := store.Get(ctx, key)
	assert.ErrorIs(t, err, errs.ErrKeyNotExist)
	assert.Nil(t, result)
}
func TestStorage_SetAndGet(t *testing.T) {
	store := lru.New()
	defer store.Close()

	ctx := context.Background()
	key := "test-key"
	value := "test-value"
	expiration := time.Minute

	// Test Set
	err := store.Set(ctx, key, value, expiration)
	assert.NoError(t, err)

	// Test Get
	result, err := store.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, value, result)
}

func TestStorage_GetExpired(t *testing.T) {
	store := lru.New()
	defer store.Close()

	ctx := context.Background()
	key := "test-key"
	value := "test-value"
	expiration := time.Second

	// Set with a short expiration
	err := store.Set(ctx, key, value, expiration)
	assert.NoError(t, err)

	// Wait for expiration
	time.Sleep(2 * time.Second)

	// Test Get after expiration
	result, err := store.Get(ctx, key)
	assert.ErrorIs(t, err, errs.ErrItemExpired)
	t.Logf("result: %v, err: %v", result, err)
	assert.Nil(t, result)
}

func TestStorage_Delete(t *testing.T) {
	store := lru.New()
	defer store.Close()

	ctx := context.Background()
	key := "test-key"
	value := "test-value"

	// Set a value
	err := store.Set(ctx, key, value, 0)
	assert.NoError(t, err)

	// Delete the value
	err = store.Delete(ctx, key)
	assert.NoError(t, err)

	// Try to Get the deleted value
	result, err := store.Get(ctx, key)
	assert.ErrorIs(t, err, errs.ErrKeyNotExist)
	assert.Nil(t, result)
}

func TestStorage_Keys(t *testing.T) {
	store := lru.New()
	defer store.Close()

	ctx := context.Background()
	key1 := "key1"
	key2 := "key2"
	value := "value"

	// Set multiple keys
	err := store.Set(ctx, key1, value, 0)
	assert.NoError(t, err)
	err = store.Set(ctx, key2, value, 0)
	assert.NoError(t, err)

	// Get all keys
	keys := store.Keys(ctx)
	assert.Contains(t, keys, key1)
	assert.Contains(t, keys, key2)
}

func TestStorage_Flush(t *testing.T) {
	store := lru.New()
	defer store.Close()

	ctx := context.Background()
	key := "test-key"
	value := "test-value"

	// Set a value
	err := store.Set(ctx, key, value, 0)
	assert.NoError(t, err)

	// Flush the storage
	err = store.Flush(ctx)
	assert.NoError(t, err)

	// Try to Get the flushed value
	result, err := store.Get(ctx, key)
	assert.ErrorIs(t, err, errs.ErrKeyNotExist)
	assert.Nil(t, result)
}

func TestStorage_Contains(t *testing.T) {
	store := lru.New()
	defer store.Close()

	ctx := context.Background()
	key := "test-key"
	value := "test-value"

	// Set a value
	err := store.Set(ctx, key, value, 0)
	assert.NoError(t, err)

	// Check if the key exists
	exists := store.Contains(ctx, key)
	assert.True(t, exists)

	// Delete the key
	err = store.Delete(ctx, key)
	assert.NoError(t, err)

	// Check if the key exists after deletion
	exists = store.Contains(ctx, key)
	assert.False(t, exists)
}
