package demo

import (
	"context"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/scraper"
)

type Client struct{}

func NewClient() *Client       { return &Client{} }
func (c *Client) Name() string { return "demo" }

func (c *Client) Search(ctx context.Context, query scraper.Query) ([]domain.Vacancy, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	title := "Go Backend Developer"
	if len(query.Keywords) > 0 && query.Keywords[0] != "" {
		title = query.Keywords[0]
	}
	return []domain.Vacancy{{
		ID: "demo-go-backend", Source: c.Name(), ExternalID: "demo-go-backend",
		Title: title, Company: "Demo Company", Description: "Демонстрационная вакансия для локального поиска.",
		Location: "Remote", URL: "https://example.com/jobs/demo-go-backend",
		Remote: true, Skills: []string{"go", "sql"},
	}}, nil
}
