package main

import (
	"sync"
	"time"
)

type CacheItem struct {
	Data      []int
	ExpiresAt time.Time
}

type Cache struct {
	data  map[string]CacheItem
	mutex sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]CacheItem),
	}
}

func (c *Cache) Get(key string) ([]int, bool) {
	c.mutex.RLock()
	item, ok := c.data[key]
	c.mutex.RUnlock()

	if ok && time.Now().Before(item.ExpiresAt) {
		return item.Data, true
	}

	if ok {
		// Item exists but is expired. Remove it.
		c.mutex.Lock()
		delete(c.data, key)
		c.mutex.Unlock()
	}

	return nil, false
}

func (c *Cache) Set(key string, value []int, duration time.Duration) {
	c.mutex.Lock()
	c.data[key] = CacheItem{
		Data:      value,
		ExpiresAt: time.Now().Add(duration),
	}
	c.mutex.Unlock()

	go c.purgeExpired() // Trigger purgeExpired asynchronously
}

func (c *Cache) purgeExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for key, item := range c.data {
		if time.Now().After(item.ExpiresAt) {
			delete(c.data, key)
		}
	}
}
