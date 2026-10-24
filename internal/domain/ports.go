package domain

import (
	"context"
	"time"
)

type CandidateRepository interface {
	GetByID(context.Context, string) (Candidate, error)
	Save(context.Context, Candidate) error
}

type VacancyRepository interface {
	GetByID(context.Context, string) (Vacancy, error)
	Save(context.Context, Vacancy) error
	List(context.Context, SearchCriteria) ([]Vacancy, error)
}

type ApplicationRepository interface {
	GetByID(context.Context, string) (Application, error)
	Save(context.Context, Application) error
	UpdateStatus(context.Context, string, ApplicationStatus) error
}

type Cache interface {
	Get(context.Context, string) (string, error)
	Set(context.Context, string, string, time.Duration) error
}
