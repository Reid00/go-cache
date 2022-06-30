package lru

import (
	"container/list"
)

type (
	Cache struct {
		maxBytes int64
		nbyte    int64
		ll       *list.List
		cache    map[string]*list.Element
		// optional and executed when an entry is purged.
		OnEvicted func(key string, value Value)
	}

	// Value use Len to count how many bytes it takes
	Value interface {
		Len() int
	}

	entry struct {
		key   string
		value Value
	}
)

// New return a cache pointer
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Add add a value to the cache
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbyte += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
		c.nbyte += int64(len(key)) + int64(value.Len())
	}

	// if maxBytes == 0, do not remove any oldest ele
	for c.maxBytes != 0 && c.maxBytes < c.nbyte {
		c.RemoveOldest()
	}
}

// Get look up a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		// move the active element to the front
		c.ll.MoveToFront(ele)
		// parse the ele to Value type
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest remove the oldest element
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbyte -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
