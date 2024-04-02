package kubenurse

import (
	"sync"
	"time"
)

type CacheEntry[K comparable] struct {
	val        K
	lastInsert time.Time
}

type TTLCache[K comparable] struct {
	m   map[K]*CacheEntry[K]
	TTL time.Duration
	mu  sync.Mutex
}

func (c *TTLCache[K]) Init(TTL time.Duration) {
	c.m = make(map[K]*CacheEntry[K])
	c.mu = sync.Mutex{}
	c.TTL = TTL
}

func (c *TTLCache[K]) Insert(k K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, ok := c.m[k]; ok {
		entry.lastInsert = time.Now()
	} else {
		entry := CacheEntry[K]{val: k, lastInsert: time.Now()}
		c.m[k] = &entry
	}
}

func (c *TTLCache[K]) ActiveEntries() int {
	c.cleanupExpired()
	return len(c.m)
}

func (c *TTLCache[K]) cleanupExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, entry := range c.m {
		if time.Since(entry.lastInsert) > c.TTL {
			delete(c.m, k)
		}
	}
}
