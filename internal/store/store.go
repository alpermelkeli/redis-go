package store

import (
	"sync"
	"time"
)

type StoreElement struct {
	expiresAt time.Time
	value     string
}

func (s *StoreElement) isExpired() bool {
	//This means it doesn't have ttl
	if s.expiresAt.IsZero() {
		return false
	}
	return s.expiresAt.Before(time.Now())
}

type Store struct {
	mu   sync.RWMutex
	data map[string]StoreElement
	stop chan struct{}
}

func New() *Store {
	s := &Store{
		data: make(map[string]StoreElement),
		stop: make(chan struct{}),
	}
	go s.cleanup()
	return s
}

func (s *Store) cleanup() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.mu.Lock()
			for key, elem := range s.data {
				if elem.isExpired() {
					delete(s.data, key)
				}
			}
			s.mu.Unlock()
		case <-s.stop:
			return
		}
	}
}

func (s *Store) Close() {
	close(s.stop)
}

// TTL was given in seconds.
func (s *Store) SetWithTTL(key, value string, ttl int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val := StoreElement{expiresAt: time.Now().Add(time.Duration(ttl) * time.Second), value: value}
	s.data[key] = val
}

func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = StoreElement{value: value}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.data[key]
	if !ok || v.isExpired() {
		return "", false
	}
	return v.value, true
}

func (s *Store) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.data[key]
	delete(s.data, key)
	return ok
}
