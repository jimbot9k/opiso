package reverse

import (
	"container/list"
	"sync"
)

// Cache which tracks the order in which keys were used, and will evict the oldest key if the cache limit is hit when setting (LRU cache)
type Cache struct {
	capacity int
	mu       sync.Mutex
	items    map[string]*list.Element
	eviction *list.List
}

type entry struct {
	key   string
	value string
}

// Creates cache which tracks the order in which keys were used, and will evict the oldest key if the cache limit is hit when setting (LRU cache)
func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		eviction: list.New(),
	}
}

// Get value of key from cache
func (c *Cache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.eviction.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}
	return "", false
}

// Sets a cache entry, removing the last used entry if at limit. Returns true if entry was evicted, false if no entry was evicted.
func (c *Cache) Set(key string, value string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	evictionOccured := false

	if elem, ok := c.items[key]; ok {
		elem.Value.(*entry).value = value
		c.eviction.MoveToFront(elem)
		return evictionOccured
	}

	if c.eviction.Len() >= c.capacity {
		oldest := c.eviction.Back()
		if oldest != nil {
			c.eviction.Remove(oldest)
			delete(c.items, oldest.Value.(*entry).key)
			evictionOccured = true
		}
	}

	newEntry := &entry{key, value}
	elem := c.eviction.PushFront(newEntry)
	c.items[key] = elem
	return evictionOccured
}