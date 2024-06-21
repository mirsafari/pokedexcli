package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	value     []byte
}

type Cache struct {
	mu       sync.Mutex
	entry    map[string]cacheEntry
	interval time.Duration
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		mu:       sync.Mutex{},
		entry:    make(map[string]cacheEntry),
		interval: interval,
	}

	go cache.reapLoop()

	return &cache
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := cacheEntry{
		createdAt: time.Now(),
		value:     value,
	}

	c.entry[key] = entry

	return
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, ok := c.entry[key]; ok {
		return entry.value, true
	}

	return nil, false
}

func (c *Cache) reapLoop() {

	for {
		time.Sleep(c.interval)

		cutoff := time.Now().Add(-c.interval)

		for key, entry := range c.entry {

			if entry.createdAt.Before(cutoff) {
				c.mu.Lock()
				delete(c.entry, key)
				c.mu.Unlock()
			}

		}
	}
}
