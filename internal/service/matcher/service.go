package matcher

import (
	"context"
	"sort"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
)

type Service struct{}

// Service — двухуровневый matcher: фильтры, embeddings и затем LLM.
func NewService() *Service { return &Service{} }

func (s *Service) Match(ctx context.Context, candidate domain.Candidate, vacancies []domain.Vacancy) []domain.Match {
	// Match — формирует результаты сравнения кандидата с вакансиями.
	result := make([]domain.Match, 0, len(vacancies))
	for _, vacancy := range vacancies {
		if err := ctx.Err(); err != nil {
			return result
		}
		if !vacancy.IsRemoteCompatible(candidate.Criteria.Remote) {
			continue
		}
		if !vacancy.MatchesCriteria(candidate.Criteria) {
			continue
		}
		score := skillScore(candidate.Skills, vacancy.Skills)
		result = append(result, domain.Match{
			UserID: candidate.ID, VacancyID: vacancy.ID,
			Score:   score,
			Reasons: []string{"passed basic criteria", "skill overlap score"},
		})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Score > result[j].Score })
	return result
}

func skillScore(candidateSkills, vacancySkills []string) float64 {
	if len(candidateSkills) == 0 || len(vacancySkills) == 0 {
		return 0
	}
	known := make(map[string]struct{}, len(candidateSkills))
	for _, skill := range candidateSkills {
		known[skill] = struct{}{}
	}
	matches := 0
	for _, skill := range vacancySkills {
		if _, ok := known[skill]; ok {
			matches++
		}
	}
	return float64(matches) / float64(len(vacancySkills))
}
