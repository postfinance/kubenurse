package kubenurse

import (
	"sync"
	"time"
)

// This TTLCache is a suboptimal, first-shot implementation of TTL cache,
// that is, entries expire after a certain duration.
// It should ideally be implemented with a Red-Black Binary tree to keep track
// of the entries to expire.

type TTLCache[K comparable] struct {
	m   map[K]*CacheEntry[K]
	TTL time.Duration
	mu  sync.Mutex
}

type CacheEntry[K comparable] struct {
	val        K
	lastInsert time.Time
}

func (c *TTLCache[K]) Init(ttl time.Duration) {
	c.m = make(map[K]*CacheEntry[K])
	c.mu = sync.Mutex{}
	c.TTL = ttl
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
