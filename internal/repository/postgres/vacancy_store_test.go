package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
)

func TestVacancyStoreUpsertsBySourceAndExternalID(t *testing.T) {
	t.Parallel()
	store := NewVacancyStore()
	first := domain.Vacancy{ID: "first", Source: "HH", ExternalID: "42", Title: "Old", Company: "Acme", URL: "https://example.com/old"}
	second := first
	second.ID = "second"
	second.Title = "New"
	if err := store.Save(context.Background(), first); err != nil {
		t.Fatal(err)
	}
	if err := store.Save(context.Background(), second); err != nil {
		t.Fatal(err)
	}
	got, err := store.GetByID(context.Background(), "first")
	if err != nil || got.Title != "New" {
		t.Fatalf("upsert result = %+v, %v", got, err)
	}
	if _, err := store.GetByID(context.Background(), "second"); err == nil {
		t.Fatal("upsert created duplicate ID")
	}
}

func TestVacancyStoreListsNewestFirst(t *testing.T) {
	t.Parallel()
	store := NewVacancyStore()
	for _, item := range []domain.Vacancy{
		{ID: "old", Source: "x", Title: "Old", Company: "A", URL: "https://example.com/old", PublishedAt: time.Unix(1, 0)},
		{ID: "new", Source: "x", Title: "New", Company: "B", URL: "https://example.com/new", PublishedAt: time.Unix(2, 0)},
	} {
		if err := store.Save(context.Background(), item); err != nil {
			t.Fatal(err)
		}
	}
	got, err := store.List(context.Background(), domain.SearchCriteria{})
	if err != nil || len(got) != 2 || got[0].ID != "new" {
		t.Fatalf("List() = %+v, %v", got, err)
	}
}
