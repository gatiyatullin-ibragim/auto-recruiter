package postgres

import (
	"context"
	"sort"
	"sync"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
)

type VacancyStore struct {
	mu    sync.RWMutex
	items map[string]domain.Vacancy
}

func NewVacancyStore() *VacancyStore {
	return &VacancyStore{items: make(map[string]domain.Vacancy)}
}

func (s *VacancyStore) GetByID(_ context.Context, id string) (domain.Vacancy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	if !ok {
		return domain.Vacancy{}, domain.ErrNotFound
	}
	return item, nil
}

func (s *VacancyStore) Save(_ context.Context, item domain.Vacancy) error {
	item = item.Normalize()
	if err := item.Validate(); err != nil {
		return err
	}
	if item.ID == "" {
		return domain.ErrInvalidInput
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, existing := range s.items {
		if id != item.ID && item.ExternalID != "" &&
			existing.Source == item.Source && existing.ExternalID == item.ExternalID {
			item.ID = id
			break
		}
	}
	s.items[item.ID] = item
	return nil
}

func (s *VacancyStore) List(_ context.Context, criteria domain.SearchCriteria) ([]domain.Vacancy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.Vacancy, 0, len(s.items))
	for _, item := range s.items {
		if item.MatchesCriteria(criteria) {
			result = append(result, item)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].PublishedAt.After(result[j].PublishedAt)
	})
	return result, nil
}
