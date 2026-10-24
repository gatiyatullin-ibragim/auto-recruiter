package postgres

import (
	"context"
	"sync"
	"time"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
)

type ApplicationStore struct {
	mu    sync.RWMutex
	items map[string]domain.Application
}

func NewApplicationStore() *ApplicationStore {
	return &ApplicationStore{items: make(map[string]domain.Application)}
}

func (s *ApplicationStore) Save(_ context.Context, item domain.Application) error {
	if item.ID == "" || item.CandidateID == "" || item.VacancyID == "" {
		return domain.ErrInvalidInput
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.items[item.ID]; exists {
		return domain.ErrConflict
	}
	s.items[item.ID] = item
	return nil
}

func (s *ApplicationStore) GetByID(_ context.Context, id string) (domain.Application, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	if !ok {
		return domain.Application{}, domain.ErrNotFound
	}
	return item, nil
}

func (s *ApplicationStore) UpdateStatus(_ context.Context, id string, status domain.ApplicationStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[id]
	if !ok {
		return domain.ErrNotFound
	}
	if !item.CanTransitionTo(status) {
		return domain.ErrInvalidInput
	}
	item.Status = status
	item.UpdatedAt = time.Now()
	s.items[id] = item
	return nil
}
