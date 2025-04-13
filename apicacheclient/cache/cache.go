package cache

import "C"
import (
	"container/list"
	"sync"
	"time"
)

type cacheItem struct {
	key       string
	value     interface{}
	expiresAt time.Time
}

type LRUCache struct {
	capacity int
	ttl      time.Duration
	store    map[string]*list.Element
	list     *list.List
	lock     sync.Mutex
}

func NewLRUCache(MaxCapacity int, ttl time.Duration) *LRUCache {
	return &LRUCache{
		capacity: MaxCapacity,
		ttl:      ttl,
		store:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if elem, found := c.store[key]; found {
		item := elem.Value.(*cacheItem)
		if time.Now().After(item.expiresAt) {
			c.list.Remove(elem)
			delete(c.store, key)
			return nil, false
		}
		c.list.MoveToFront(elem)
		return item.value, true
	}
	return nil, false
}

func (c *LRUCache) Set(key string, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if elem, found := c.store[key]; found {
		c.list.MoveToFront(elem)
		elem.Value.(*cacheItem).value = value
		elem.Value.(*cacheItem).expiresAt = time.Now().Add(c.ttl)
		return
	}

	if c.list.Len() >= c.capacity {
		tail := c.list.Back()
		if tail != nil {
			item := tail.Value.(*cacheItem)
			delete(c.store, item.key)
			c.list.Remove(tail)
		}
	}

	ci := &cacheItem{
		key:       key,
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
	elem := c.list.PushFront(ci)
	c.store[key] = elem
}
