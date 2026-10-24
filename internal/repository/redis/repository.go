package redis

import (
	"context"
	"sync"
	"time"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
)

type Repository interface {
	Get(context.Context, string) (string, error)
	Set(context.Context, string, string, time.Duration) error
}

type entry struct {
	value     string
	expiresAt time.Time
}

type Store struct {
	mu     sync.RWMutex
	values map[string]entry
}

func NewStore() *Store {
	return &Store{values: make(map[string]entry)}
}

func (s *Store) Get(_ context.Context, key string) (string, error) {
	s.mu.RLock()
	item, ok := s.values[key]
	s.mu.RUnlock()
	if !ok || (!item.expiresAt.IsZero() && time.Now().After(item.expiresAt)) {
		return "", domain.ErrCacheMiss
	}
	return item.value, nil
}

func (s *Store) Set(_ context.Context, key, value string, ttl time.Duration) error {
	if key == "" {
		return domain.ErrInvalidInput
	}
	item := entry{value: value}
	if ttl > 0 {
		item.expiresAt = time.Now().Add(ttl)
	}
	s.mu.Lock()
	s.values[key] = item
	s.mu.Unlock()
	return nil
}
