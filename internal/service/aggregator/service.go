package aggregator

import (
	"context"
	"errors"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/scraper"
)

// Service — use case агрегации вакансий из нескольких независимых источников.
type Service struct{ clients []scraper.Client }

// NewService — конструктор с dependency injection scraper-клиентов.
func NewService(clients ...scraper.Client) *Service { return &Service{clients: clients} }

func (s *Service) Collect(ctx context.Context, query scraper.Query) ([]domain.Vacancy, error) {
	// Collect — вызывает источники и нормализует их ответы в domain.Vacancy.
	var result []domain.Vacancy
	var sourceErrors []error
	for _, client := range s.clients {
		items, err := client.Search(ctx, query)
		if err != nil {
			sourceErrors = append(sourceErrors, err)
			continue
		}
		for _, item := range items {
			result = append(result, item.Normalize())
		}
	}
	return result, errors.Join(sourceErrors...)
}
