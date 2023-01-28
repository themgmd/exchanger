package currency

import "time"

type InMemory interface {
	Put(key string, data float64, ttl time.Duration)
	Get(key string) (float64, bool)
	Delete(key string)
}
