package ai

import (
	"container/list"
)

type Cache struct {
	capacity int
	cache    *list.List
	mapCache map[string]uint64
}

func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		cache:    list.New(),
		mapCache: make(map[string]uint64),
	}
}

func (c *Cache) Get(key string) uint64 {
	if !config.EnableCache {
		return 0
	}
	if val, ok := c.mapCache[key]; ok {
		return val
	}
	return 0
}

func (c *Cache) Put(key string, value uint64) {
	if !config.EnableCache {
		return
	}
	if c.cache.Len() >= c.capacity {
		oldest := c.cache.Back()
		c.cache.Remove(oldest)
		delete(c.mapCache, oldest.Value.(string))
	}

	if _, ok := c.mapCache[key]; !ok {
		c.cache.PushFront(key)
	}
	c.mapCache[key] = value
}

func (c *Cache) Has(key string) bool {
	if !config.EnableCache {
		return false
	}
	_, ok := c.mapCache[key]
	return ok
}
