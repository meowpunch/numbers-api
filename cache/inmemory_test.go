package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCacheSetAndGet(t *testing.T) {
	cache := NewInMemoryCache(10 * time.Millisecond)

	// Test setting and getting values
	cache.Set("key1", []int{1, 2, 3}, 10*time.Millisecond)
	val, ok := cache.Get("key1")

	require.True(t, ok, "Expected key1 to exist")
	assert.Equal(t, []int{1, 2, 3}, val, "Expected to retrieve [1,2,3]")
}

func TestExpirationBeforePurge(t *testing.T) {
	cache := NewInMemoryCache(30 * time.Millisecond)

	cache.Set("key1", []int{1, 2, 3}, 10*time.Millisecond)
	cache.Set("key2", []int{4, 5, 6}, 20*time.Millisecond)

	time.Sleep(10 * time.Millisecond)

	_, ok1 := cache.Get("key1")
	_, ok2 := cache.Get("key2")

	assert.False(t, ok1, "Expected key1 to have expired after purging.")
	assert.True(t, ok2, "Expected key2 to not have expired after purging.")
}

func TestExpirationAfterPurge(t *testing.T) {
	cache := NewInMemoryCache(10 * time.Millisecond)

	cache.Set("key1", []int{1, 2, 3}, 10*time.Millisecond)
	cache.Set("key2", []int{4, 5, 6}, 20*time.Millisecond)

	time.Sleep(10 * time.Millisecond)

	_, ok1 := cache.Get("key1")
	_, ok2 := cache.Get("key2")

	assert.False(t, ok1, "Expected key1 to have expired after purging.")
	assert.True(t, ok2, "Expected key2 to not have expired after purging.")
}
