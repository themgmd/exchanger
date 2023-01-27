package inmemory

import (
	"sync"
	"time"
)

type storage map[string]float64

type InMemory struct {
	mx sync.RWMutex
	storage
}

func New() *InMemory {
	store := make(storage, 10)
	return &InMemory{storage: store}
}

func (inm *InMemory) Put(key string, data float64, ttl time.Duration) {
	inm.mx.Lock()
	inm.storage[key] = data
	inm.mx.Unlock()
	go inm.clearData(key, ttl)
}

func (inm *InMemory) Get(key string) (float64, bool) {
	inm.mx.RLock()
	defer inm.mx.RUnlock()
	data, ok := inm.storage[key]
	return data, ok
}

func (inm *InMemory) Delete(key string) {
	inm.mx.Lock()
	delete(inm.storage, key)
	inm.mx.Unlock()
}

func (inm *InMemory) clearData(key string, ttl time.Duration) {
	select {
	case <-time.After(ttl):
		delete(inm.storage, key)
	}
}
