package service

import (
	"context"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/repository/postgres"
)

type CandidateService struct {
	repository postgres.CandidateRepository
}

// NewCandidateService — связывает use case с абстракцией хранилища.
func NewCandidateService(repository postgres.CandidateRepository) *CandidateService {
	return &CandidateService{repository: repository}
}

func (s *CandidateService) Create(ctx context.Context, candidate domain.Candidate) error {
	// Create — валидирует и сохраняет новый профиль кандидата.
	if candidate.ID == "" {
		return domain.ErrInvalidInput
	}
	return s.repository.Save(ctx, candidate)
}

func (s *CandidateService) Get(ctx context.Context, id string) (domain.Candidate, error) {
	// Get — получает профиль по внутреннему идентификатору.
	return s.repository.GetByID(ctx, id)
}

func (s *CandidateService) Save(ctx context.Context, candidate domain.Candidate) error {
	return s.repository.Save(ctx, candidate)
}
