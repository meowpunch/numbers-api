package cache

import (
	"sync"
	"time"
)

type Item struct {
	Data      []int
	ExpiresAt time.Time
}

type InMemoryCache struct {
	data          map[string]Item
	mutex         sync.RWMutex
	purgeInterval time.Duration
}

func NewInMemoryCache(purgeInterval time.Duration) *InMemoryCache {
	c := &InMemoryCache{
		data:          make(map[string]Item),
		purgeInterval: purgeInterval,
	}
	go c.purgeExpired() // Regularly purge expired items
	return c
}

func (c *InMemoryCache) Get(key string) ([]int, bool) {
	c.mutex.RLock()
	item, ok := c.data[key]
	c.mutex.RUnlock()

	if ok && time.Now().Before(item.ExpiresAt) {
		return item.Data, true
	}

	if ok {
		c.mutex.Lock()
		delete(c.data, key)
		c.mutex.Unlock()
	}

	return nil, false
}

func (c *InMemoryCache) Set(key string, value []int, duration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = Item{
		Data:      value,
		ExpiresAt: time.Now().Add(duration),
	}
}

func (c *InMemoryCache) purgeExpired() {
	ticker := time.NewTicker(c.purgeInterval) // Adjust as needed
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		for key, item := range c.data {
			if time.Now().After(item.ExpiresAt) {
				delete(c.data, key)
			}
		}
		c.mutex.Unlock()
	}
}
