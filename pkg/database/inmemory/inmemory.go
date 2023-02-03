package inmemory

import (
	"sync"
	"time"
)

type InMemory interface {
	Put(key string, data float64, ttl time.Duration)
	Get(key string) (float64, bool)
	Delete(key string)
}

type storage map[string]float64

type inMemory struct {
	mx sync.RWMutex
	storage
}

func New() InMemory {
	store := make(storage, 10)
	return &inMemory{storage: store}
}

func (inm *inMemory) Put(key string, data float64, ttl time.Duration) {
	inm.mx.Lock()
	inm.storage[key] = data
	inm.mx.Unlock()
	go inm.clearData(key, ttl)
}

func (inm *inMemory) Get(key string) (data float64, ok bool) {
	inm.mx.RLock()
	data, ok = inm.storage[key]
	inm.mx.RUnlock()
	return
}

func (inm *inMemory) Delete(key string) {
	inm.mx.Lock()
	delete(inm.storage, key)
	inm.mx.Unlock()
}

func (inm *inMemory) clearData(key string, ttl time.Duration) {
	select {
	case <-time.After(ttl):
		delete(inm.storage, key)
	}
}
