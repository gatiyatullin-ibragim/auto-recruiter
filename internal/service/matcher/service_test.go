package matcher

import (
	"context"
	"testing"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
)

func TestServiceMatchFiltersAndSorts(t *testing.T) {
	t.Parallel()
	candidate := domain.Candidate{ID: "c1", Skills: []string{"go", "sql"}, Criteria: domain.SearchCriteria{Remote: true}}
	vacancies := []domain.Vacancy{
		{ID: "low", Remote: true, Skills: []string{"go", "sql", "aws"}},
		{ID: "high", Remote: true, Skills: []string{"go"}},
		{ID: "local", Location: "Almaty", Skills: []string{"go"}},
	}
	got := NewService().Match(context.Background(), candidate, vacancies)
	if len(got) != 2 || got[0].VacancyID != "high" || got[1].VacancyID != "low" {
		t.Fatalf("unexpected matches: %+v", got)
	}
}

func TestServiceMatchHonorsCancellation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	got := NewService().Match(ctx, domain.Candidate{ID: "c1"}, []domain.Vacancy{{ID: "v1"}})
	if len(got) != 0 {
		t.Fatalf("cancelled match returned results: %+v", got)
	}
}
