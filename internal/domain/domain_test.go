package domain

import (
	"errors"
	"math"
	"testing"
)

func TestCandidateValidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		candidate Candidate
		wantErr   bool
	}{
		{"valid", Candidate{ID: "c1", Name: "Ada", Email: "ada@example.com", Skills: []string{"go"}}, false},
		{"missing identity", Candidate{Name: "Ada", Email: "ada@example.com", Skills: []string{"go"}}, true},
		{"invalid email", Candidate{ID: "c1", Name: "Ada", Email: "invalid", Skills: []string{"go"}}, true},
		{"empty profile", Candidate{ID: "c1", Name: "Ada", Email: "ada@example.com"}, true},
		{"negative salary", Candidate{ID: "c1", Name: "Ada", Email: "ada@example.com", Skills: []string{"go"}, Criteria: SearchCriteria{MinSalary: -1}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.candidate.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, want error: %v", err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, ErrInvalidInput) {
				t.Fatalf("Validate() error = %v, want ErrInvalidInput", err)
			}
		})
	}
}

func TestVacancyValidateRequiresAbsoluteHTTPURL(t *testing.T) {
	t.Parallel()
	base := Vacancy{Source: "hh", ExternalID: "1", Title: "Go", Company: "Acme", URL: "https://example.com/jobs/1"}
	if err := base.Validate(); err != nil {
		t.Fatalf("valid vacancy rejected: %v", err)
	}
	for _, url := range []string{"example.com/jobs/1", "/jobs/1", "ftp://example.com/jobs/1"} {
		v := base
		v.URL = url
		if err := v.Validate(); err == nil {
			t.Fatalf("Validate(%q) accepted invalid URL", url)
		}
	}
}

func TestVacancyNormalizeAndCriteria(t *testing.T) {
	t.Parallel()
	v := (Vacancy{
		Source: " HH ", Title: " Go ", Company: " Acme ", Location: "Remote",
		URL: "https://example.com", Skills: []string{"Go", " go ", "", "SQL"},
	}).Normalize()
	if v.Source != "hh" || v.Title != "Go" || len(v.Skills) != 2 || v.Skills[0] != "go" {
		t.Fatalf("unexpected normalized vacancy: %+v", v)
	}
	if !v.MatchesCriteria(SearchCriteria{Remote: true}) {
		t.Fatal("remote vacancy should match remote criteria")
	}
	if (Vacancy{Location: "Almaty"}).MatchesCriteria(SearchCriteria{Remote: true}) {
		t.Fatal("non-remote vacancy should not match remote criteria")
	}
}

func TestMatchAccepted(t *testing.T) {
	t.Parallel()
	tests := []struct {
		score, threshold float64
		want             bool
	}{{0, 0, true}, {0.9, 0.8, true}, {0.79, 0.8, false}, {math.NaN(), 0.8, false}}
	for _, tt := range tests {
		if got := (Match{Score: tt.score}).Accepted(tt.threshold); got != tt.want {
			t.Errorf("Accepted(%v) = %v, want %v", tt.threshold, got, tt.want)
		}
	}
}

func TestApplicationTransitions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		from, to ApplicationStatus
		want     bool
	}{{ApplicationDraft, ApplicationGenerated, true}, {ApplicationGenerated, ApplicationSent, true}, {ApplicationGenerated, ApplicationRejected, true}, {ApplicationSent, ApplicationDraft, false}, {ApplicationDraft, ApplicationSent, false}}
	for _, tt := range tests {
		if got := (Application{Status: tt.from}).CanTransitionTo(tt.to); got != tt.want {
			t.Errorf("%s -> %s = %v, want %v", tt.from, tt.to, got, tt.want)
		}
	}
}
