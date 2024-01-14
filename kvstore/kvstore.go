package kvstore

import (
	"context"
	"time"
)

const (
	// DefaultExpiration uses the store's default expiration time.
	DefaultExpiration time.Duration = 0

	// NoExpiration indicates that the item should never expire.
	NoExpiration time.Duration = -1
)

type Reader interface {
	// Read fetches an item by key. Returns nil if not found.
	Read(ctx context.Context, key string) any
}

type Writer interface {
	// Write sets an item with a specified expiry. Use DefaultExpiration for the store's default
	// expiration time, or NoExpiration for no expiration.
	Write(ctx context.Context, key string, value any, expiry time.Duration) error
}
