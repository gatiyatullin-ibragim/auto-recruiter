package scraper

import (
	"context"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
)

var ErrExternalUnavailable = domain.ErrExternalUnavailable

type Query struct {
	// Query — унифицированный запрос к любому источнику вакансий.
	Keywords []string `json:"keywords"`
	Location string   `json:"location"`
	Remote   bool     `json:"remote"`
}

type Client interface {
	// Client — контракт адаптера внешнего job-board.
	Name() string
	Search(context.Context, Query) ([]domain.Vacancy, error)
}
