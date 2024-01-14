package kvstore

import (
	"context"
	"github.com/patrickmn/go-cache"
	"time"
)

type inMemoryStore struct {
	cache *cache.Cache
}

// NewInMemoryStore returns a new key value store with a given default expiration and
// cleanup interval
func NewInMemoryStore(defaultExpiration, cleanupInterval time.Duration) *inMemoryStore {
	return &inMemoryStore{
		cache: cache.New(defaultExpiration, cleanupInterval),
	}
}

// Write an item to the store, overriding the existing one
func (ims *inMemoryStore) Write(_ context.Context, key string, value any, expiry time.Duration) error {
	ims.cache.Set(key, value, expiry)
	return nil
}

// Read an item from the store, returns the item or nil if not found
func (ims *inMemoryStore) Read(_ context.Context, key string) any {
	val, _ := ims.cache.Get(key)
	return val
}
