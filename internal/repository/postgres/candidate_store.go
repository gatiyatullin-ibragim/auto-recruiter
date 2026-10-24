package postgres

import (
	"context"
	"sync"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
)

type CandidateStore struct {
	mu    sync.RWMutex
	items map[string]domain.Candidate
}

func NewCandidateStore() *CandidateStore {
	return &CandidateStore{items: make(map[string]domain.Candidate)}
}

func (s *CandidateStore) GetByID(_ context.Context, id string) (domain.Candidate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	if !ok {
		return domain.Candidate{}, domain.ErrNotFound
	}
	return item, nil
}

func (s *CandidateStore) Save(_ context.Context, item domain.Candidate) error {
	if item.ID == "" {
		return domain.ErrInvalidInput
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[item.ID] = item
	return nil
}
