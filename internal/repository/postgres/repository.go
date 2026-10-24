package postgres

import (
	"context"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
)

type CandidateRepository interface {
	// CandidateRepository — порт хранения профиля, не раскрывающий SQL.
	GetByID(context.Context, string) (domain.Candidate, error)
	Save(context.Context, domain.Candidate) error
}

type VacancyRepository interface {
	// VacancyRepository — порт хранения и поиска вакансий.
	GetByID(context.Context, string) (domain.Vacancy, error)
	Save(context.Context, domain.Vacancy) error
	List(context.Context, domain.SearchCriteria) ([]domain.Vacancy, error)
}

type ApplicationRepository interface {
	GetByID(context.Context, string) (domain.Application, error)
	Save(context.Context, domain.Application) error
	UpdateStatus(context.Context, string, domain.ApplicationStatus) error
}
