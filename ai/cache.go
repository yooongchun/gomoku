package ai

import (
	"container/list"
)

type Cache struct {
	capacity int
	cache    *list.List
	mapCache map[uint64]interface{}
}

func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		cache:    list.New(),
		mapCache: make(map[uint64]interface{}),
	}
}

func (c *Cache) Get(key uint64) interface{} {
	if val, ok := c.mapCache[key]; ok {
		return val
	}
	return nil
}

func (c *Cache) Put(key uint64, value interface{}) {
	if c.cache.Len() >= c.capacity {
		oldest := c.cache.Back()
		c.cache.Remove(oldest)
		delete(c.mapCache, oldest.Value.(uint64))
	}

	if _, ok := c.mapCache[key]; !ok {
		c.cache.PushFront(key)
	}
	c.mapCache[key] = value
}

func (c *Cache) Has(key uint64) bool {
	_, ok := c.mapCache[key]
	return ok
}
