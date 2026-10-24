package linkedin

import (
	"context"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/scraper"
)

type Client struct{}

// Client — адаптер LinkedIn, который скрывает внешний API от сервисов.
func NewClient() *Client       { return &Client{} }
func (c *Client) Name() string { return "linkedin" }
func (c *Client) Search(context.Context, scraper.Query) ([]domain.Vacancy, error) {
	return nil, scraper.ErrExternalUnavailable
}
