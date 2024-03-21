package data

import "sync"

type Json map[string]any

type Sanitizer interface {
	Sanitize(data Json, keys ...string) Json
}

type sanitizer struct {
	mx *sync.Mutex
}

func NewSanitizer() Sanitizer {
	return &sanitizer{
		mx: &sync.Mutex{},
	}
}

func (s sanitizer) Sanitize(data Json, protected ...string) Json {
	for idx := range protected {
		s.mx.Lock()
		delete(data, protected[idx])
		s.mx.Unlock()
	}

	return data
}
