package domain

import (
	"fmt"
	"net/mail"
	"strings"
	"time"
)

type Candidate struct {
	// Candidate — профиль соискателя и его поисковые предпочтения.
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Email      string         `json:"email"`
	ResumeText string         `json:"resume_text"`
	Skills     []string       `json:"skills"`
	Criteria   SearchCriteria `json:"criteria"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type SearchCriteria struct {
	// SearchCriteria — фильтры, которые применяются до embeddings и LLM.
	Position       string   `json:"position"`
	Roles          []string `json:"roles"`
	Locations      []string `json:"locations"`
	WorkMode       string   `json:"work_mode"`
	EmploymentType string   `json:"employment_type"`
	MinSalary      int64    `json:"min_salary"`
	Currency       string   `json:"currency"`
	Remote         bool     `json:"remote"`
}

func (c Candidate) Validate() error {
	if strings.TrimSpace(c.ID) == "" || strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("%w: candidate id and name are required", ErrInvalidInput)
	}
	if _, err := mail.ParseAddress(c.Email); err != nil {
		return fmt.Errorf("%w: invalid candidate email", ErrInvalidInput)
	}
	if strings.TrimSpace(c.ResumeText) == "" && len(c.Skills) == 0 {
		return fmt.Errorf("%w: resume or skills are required", ErrInvalidInput)
	}
	if c.Criteria.MinSalary < 0 {
		return fmt.Errorf("%w: minimum salary cannot be negative", ErrInvalidInput)
	}
	if position := strings.TrimSpace(c.Criteria.Position); position != "" && len([]rune(position)) < 2 {
		return fmt.Errorf("%w: position must contain at least 2 characters", ErrInvalidInput)
	}
	if c.Criteria.WorkMode != "" && c.Criteria.WorkMode != "remote" &&
		c.Criteria.WorkMode != "hybrid" && c.Criteria.WorkMode != "any" {
		return fmt.Errorf("%w: unsupported work mode", ErrInvalidInput)
	}
	return nil
}
