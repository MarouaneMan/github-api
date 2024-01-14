package kvstore

import (
	"context"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInMemoryStore(t *testing.T) {

	ctx := context.Background()
	store := NewInMemoryStore(50*time.Millisecond, 10*time.Millisecond)
	key := "testKey"
	value := "testValue"

	t.Run("Write", func(t *testing.T) {
		err := store.Write(ctx, key, value, cache.DefaultExpiration)
		assert.NoError(t, err, "should not error on write")
	})

	t.Run("Read", func(t *testing.T) {
		readValue := store.Read(ctx, key)
		assert.Equal(t, value, readValue, "read value should match written value")
	})

	t.Run("Expiration", func(t *testing.T) {
		time.Sleep(60 * time.Millisecond) // Wait for the item to expire
		readValue := store.Read(ctx, key)
		assert.Nil(t, readValue, "read value should be nil after expiration")
	})

	t.Run("ReadNonExistent", func(t *testing.T) {
		readValue := store.Read(ctx, "non-existent-key")
		assert.Nil(t, readValue, "read value should be nil for non-existent key")
	})
}
