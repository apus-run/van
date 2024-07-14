package store

import (
	"context"
	"errors"
	"time"
)

var ErrKeyNotFound = errors.New("key not found")

// Store is the interface for all available stores.
type Store interface {
	Get(ctx context.Context, key any) (any, error)
	GetWithTTL(ctx context.Context, key any) (any, time.Duration, error)
	Set(ctx context.Context, key any, value any) error
	SetWithTTL(ctx context.Context, key any, value any, ttl time.Duration) error
	Del(ctx context.Context, key any) error
	// Clear removes items that have an expired TTL.
	Clear(ctx context.Context) error
	Wait(ctx context.Context)
}
