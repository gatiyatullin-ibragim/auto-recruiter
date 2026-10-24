package service

import (
	"context"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/repository/postgres"
)

type ApplicationService struct {
	repository postgres.ApplicationRepository
}

// NewApplicationService — конструктор use case откликов.
func NewApplicationService(repository postgres.ApplicationRepository) *ApplicationService {
	return &ApplicationService{repository: repository}
}

func (s *ApplicationService) Create(ctx context.Context, application domain.Application) error {
	// Create — создаёт черновик отклика с начальным статусом.
	if application.Status == "" {
		application.Status = domain.ApplicationDraft
	}
	return s.repository.Save(ctx, application)
}

func (s *ApplicationService) ChangeStatus(ctx context.Context, id string, status domain.ApplicationStatus) error {
	return s.repository.UpdateStatus(ctx, id, status)
}
