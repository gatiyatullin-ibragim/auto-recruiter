package aggregator

import (
	"context"
	"errors"
	"testing"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/scraper"
)

type fakeClient struct {
	name  string
	items []domain.Vacancy
	err   error
}

func (f fakeClient) Name() string { return f.name }
func (f fakeClient) Search(context.Context, scraper.Query) ([]domain.Vacancy, error) {
	return f.items, f.err
}

func TestCollectKeepsSuccessfulSources(t *testing.T) {
	t.Parallel()
	sourceErr := errors.New("source down")
	got, err := NewService(
		fakeClient{name: "failed", err: sourceErr},
		fakeClient{name: "ok", items: []domain.Vacancy{{Source: " HH ", Title: " Go ", Company: " Acme ", URL: "https://example.com"}}},
	).Collect(context.Background(), scraper.Query{})
	if !errors.Is(err, sourceErr) {
		t.Fatalf("Collect() error = %v, want source error", err)
	}
	if len(got) != 1 || got[0].Source != "hh" || got[0].Title != "Go" {
		t.Fatalf("successful source was not retained: %+v", got)
	}
}
