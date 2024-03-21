package cache

import (
	"sync"
	"time"
)

const NoTTL = 0

type Cache struct {
	mx      *sync.RWMutex
	storage map[string]float64
}

func New() *Cache {
	store := make(map[string]float64, 10)
	return &Cache{
		mx:      &sync.RWMutex{},
		storage: store,
	}
}

func (c *Cache) Put(key string, data float64, ttl time.Duration) {
	c.mx.Lock()
	c.storage[key] = data
	c.mx.Unlock()

	go c.clearData(key, ttl)
}

func (c *Cache) Get(key string) (data float64, ok bool) {
	c.mx.RLock()
	data, ok = c.storage[key]
	c.mx.RUnlock()
	return
}

func (c *Cache) Delete(key string) {
	c.mx.Lock()
	delete(c.storage, key)
	c.mx.Unlock()
}

func (c *Cache) clearData(key string, ttl time.Duration) {
	if ttl == NoTTL {
		return
	}

	select {
	case <-time.After(ttl):
		c.mx.Lock()
		delete(c.storage, key)
		c.mx.Unlock()
	}
}
