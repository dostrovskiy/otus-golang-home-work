/*
Package hw04lrucache is an implementation of LRU cache.
*/
package hw04lrucache

import "sync"

/*
Key is a key for LRU cache. It is an alias for string type.
*/
type Key string

/*
Cache is an interface of LRU cache.
*/
type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

/*
NewCache returns a new LRU cache with given capacity.
*/
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if listItem, ok := c.items[key]; ok {
		c.queue.MoveToFront(listItem)
		return listItem.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if listItem, ok := c.items[key]; ok {
		c.queue.MoveToFront(listItem)
		listItem.Value.(*cacheItem).value = value
		return true
	}
	if c.queue.Len() >= c.capacity {
		backItem := c.queue.Back()
		delete(c.items, backItem.Value.(*cacheItem).key)
		c.queue.Remove(backItem)
	}
	listItem := c.queue.PushFront(&cacheItem{key: key, value: value})
	c.items[key] = listItem
	return false
}
