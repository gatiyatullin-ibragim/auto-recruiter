package glassdoor

import (
	"context"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/scraper"
)

type Client struct{}

// Client — адаптер Glassdoor и место для retry/rate-limit политики.
func NewClient() *Client       { return &Client{} }
func (c *Client) Name() string { return "glassdoor" }
func (c *Client) Search(context.Context, scraper.Query) ([]domain.Vacancy, error) {
	return nil, scraper.ErrExternalUnavailable
}
