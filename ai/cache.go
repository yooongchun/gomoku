package ai

import (
	"container/list"
)

type Cache struct {
	capacity int
	cache    *list.List
	mapCache map[string]interface{}
}

func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		cache:    list.New(),
		mapCache: make(map[string]interface{}),
	}
}

func (c *Cache) Get(key string) interface{} {
	if val, ok := c.mapCache[key]; ok {
		return val
	}
	return nil
}

func (c *Cache) Put(key string, value interface{}) {
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
	_, ok := c.mapCache[key]
	return ok
}
