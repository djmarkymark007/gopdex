package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	data     map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(t time.Duration) *Cache {
	res := Cache{interval: t, data: make(map[string]cacheEntry)}
	go res.reapLoop()
	return &res
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	c.data[key] = cacheEntry{createdAt: time.Now(), val: val}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	val, err := c.data[key]
	c.mu.Unlock()
	if !err {
		return nil, false
	}
	return val.val, true
}

func (c *Cache) reapLoop() {
	for {
		//time.ticker maybe
		time.Sleep(c.interval)
		t := time.Now()
		c.mu.Lock()
		for key, val := range c.data {
			if c.interval < t.Sub(val.createdAt) {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}
