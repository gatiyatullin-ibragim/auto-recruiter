package indeed

import (
	"context"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/scraper"
)

type Client struct{}

// Client — адаптер Indeed с преобразованием внешнего DTO в Vacancy.
func NewClient() *Client       { return &Client{} }
func (c *Client) Name() string { return "indeed" }
func (c *Client) Search(context.Context, scraper.Query) ([]domain.Vacancy, error) {
	return nil, scraper.ErrExternalUnavailable
}
