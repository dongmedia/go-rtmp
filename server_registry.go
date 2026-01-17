package gortmp

import "sync"

type Registry struct {
	mu      sync.RWMutex
	streams map[string]*Stream
}

func NewRegistry() *Registry {
	return &Registry{streams: make(map[string]*Stream)}
}

func (r *Registry) Get(name string) *Stream {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.streams[name]
}

func (r *Registry) Upsert(name string, s *Stream) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.streams[name] = s
}
