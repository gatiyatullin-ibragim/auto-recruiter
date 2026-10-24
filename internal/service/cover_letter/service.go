package coverletter

import (
	"context"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/ai/llm"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
)

type Generator interface {
	// Generator — контракт генерации письма, удобный для mock-тестов.
	Generate(context.Context, domain.Candidate, domain.Vacancy) (string, error)
}

// Service — use case, который формирует контекст и вызывает LLM.
type Service struct{ client llm.Client }

func NewService(client llm.Client) *Service { return &Service{client: client} }

func (s *Service) Generate(ctx context.Context, candidate domain.Candidate, vacancy domain.Vacancy) (string, error) {
	// Generate — создаёт письмо для конкретной пары кандидат-вакансия.
	if s.client == nil {
		return "", domain.ErrExternalUnavailable
	}
	prompt := "Write a concise professional cover letter for " + vacancy.Title + " at " + vacancy.Company +
		" for candidate " + candidate.Name + ". Candidate resume: " + candidate.ResumeText
	return s.client.Complete(ctx, prompt)
}
