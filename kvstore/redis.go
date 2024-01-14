package kvstore

import (
	"context"
	"time"
)

//TODO: implement this and move fetcher to its own app

type redisStore struct{}

func NewRedisStore(_, _ time.Duration) *redisStore {
	return &redisStore{}
}

// Write an item to the store, overriding the existing one
func (ims *redisStore) Write(_ context.Context, _ string, _ any, _ time.Duration) error {
	panic("redisStore.Write: unimplemented")
	return nil
}

// Read an item from the store, returns the item or nil if not found
func (ims *redisStore) Read(_ context.Context, _ string) any {
	panic("redisStore. Read: unimplemented")
	return nil
}
