package service

import (
	"context"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/repository/postgres"
)

type VacancyService struct {
	repository postgres.VacancyRepository
}

// NewVacancyService — конструктор vacancy use case через интерфейс репозитория.
func NewVacancyService(repository postgres.VacancyRepository) *VacancyService {
	return &VacancyService{repository: repository}
}

func (s *VacancyService) Create(ctx context.Context, vacancy domain.Vacancy) error {
	// Create — принимает вакансию из scraper-а и передаёт её в storage.
	return s.repository.Save(ctx, vacancy)
}

func (s *VacancyService) Get(ctx context.Context, id string) (domain.Vacancy, error) {
	// Get — возвращает одну вакансию или domain.ErrNotFound.
	return s.repository.GetByID(ctx, id)
}

func (s *VacancyService) List(ctx context.Context, criteria domain.SearchCriteria) ([]domain.Vacancy, error) {
	// List — применяет быстрые фильтры поиска вакансий.
	return s.repository.List(ctx, criteria)
}
