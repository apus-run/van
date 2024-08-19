package jwt

import (
	"context"
	"time"
)

// Storer token storage interface.
type Storer interface {
	// Set Store token data and specify expiration time.
	Set(ctx context.Context, accessToken string, val any, expiration time.Duration) error

	// Delete token data from storage.
	Delete(ctx context.Context, accessToken string) (bool, error)

	// Check if token exists.
	Check(ctx context.Context, accessToken string) (bool, error)
}
