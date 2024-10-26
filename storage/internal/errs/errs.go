package errs

import "errors"

var (
	// ErrItemExpired is returned in Storage.Get when the item found in the cache
	// has expired.
	ErrItemExpired error = errors.New("item has expired")
	// ErrKeyNotExist is returned in Storage.Get and Storage.Delete when the
	// provided key could not be found in cache.
	ErrKeyNotExist error = errors.New("key not found in cache")

	ErrDeleteKeyFailed error = errors.New("delete key failed")
	ErrSetKeyFailed    error = errors.New("set key failed")
)
