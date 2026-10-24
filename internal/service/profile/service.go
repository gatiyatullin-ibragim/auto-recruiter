package profile

import (
	"context"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
	"strings"
)

type Service struct{}

// Service — use case нормализации и подготовки профиля кандидата.
func NewService() *Service { return &Service{} }

func (s *Service) Normalize(_ context.Context, candidate domain.Candidate) (domain.Candidate, error) {
	// Normalize — очищает базовые поля профиля перед сохранением.
	candidate.Name = strings.TrimSpace(candidate.Name)
	candidate.Email = strings.ToLower(strings.TrimSpace(candidate.Email))
	candidate.ResumeText = strings.TrimSpace(candidate.ResumeText)
	for i := range candidate.Skills {
		candidate.Skills[i] = strings.ToLower(strings.TrimSpace(candidate.Skills[i]))
	}
	if err := candidate.Validate(); err != nil {
		return domain.Candidate{}, err
	}
	return candidate, nil
}
