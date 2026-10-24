package search

import (
	"context"
	"testing"
	"time"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/repository/postgres"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/scraper/demo"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/service/aggregator"
)

func TestSearchDeliversHistoryToLateSubscriber(t *testing.T) {
	t.Parallel()
	store := postgres.NewVacancyStore()
	service := NewService(aggregator.NewService(demo.NewClient()), store)
	jobID := service.Start(context.Background(), domain.SearchCriteria{Position: "Go Developer"})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	events := service.Subscribe(ctx, jobID)
	var last Event
	for event := range events {
		last = event
		if event.Status == "done" {
			break
		}
	}
	if last.Status != "done" {
		t.Fatalf("last event = %+v, want done", last)
	}
	items, err := store.List(context.Background(), domain.SearchCriteria{})
	if err != nil || len(items) != 1 {
		t.Fatalf("stored vacancies = %d, error = %v", len(items), err)
	}
}
