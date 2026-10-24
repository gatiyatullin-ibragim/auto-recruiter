package domain

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Vacancy struct {
	// Vacancy — доменная модель вакансии, независимая от LinkedIn, Indeed
	// и базы данных. Она хранит данные, с которыми работают сервисы системы.
	ID             string    `json:"id"`
	Source         string    `json:"source"`
	ExternalID     string    `json:"external_id"`
	Title          string    `json:"title"`
	Company        string    `json:"company"`
	Description    string    `json:"description"`
	Location       string    `json:"location"`
	URL            string    `json:"url"`
	EmploymentType string    `json:"employment_type"`
	Remote         bool      `json:"remote"`
	SalaryMin      int64     `json:"salary_min"`
	SalaryMax      int64     `json:"salary_max"`
	Currency       string    `json:"currency"`
	Skills         []string  `json:"skills"`
	PublishedAt    time.Time `json:"published_at"`
	CollectedAt    time.Time `json:"collected_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (v Vacancy) Validate() error {
	// Validate — проверяет бизнес-инварианты до сохранения вакансии.
	switch {
	case strings.TrimSpace(v.Title) == "":
		return fmt.Errorf("%w: vacancy title is required", ErrInvalidInput)
	case strings.TrimSpace(v.Source) == "":
		return fmt.Errorf("%w: vacancy source is required", ErrInvalidInput)
	case strings.TrimSpace(v.Company) == "":
		return fmt.Errorf("%w: vacancy company is required", ErrInvalidInput)
	case strings.TrimSpace(v.URL) == "":
		return fmt.Errorf("%w: vacancy URL is required", ErrInvalidInput)
	case v.SalaryMin < 0 || v.SalaryMax < 0:
		return fmt.Errorf("%w: vacancy salary cannot be negative", ErrInvalidInput)
	case v.SalaryMax > 0 && v.SalaryMin > v.SalaryMax:
		return fmt.Errorf("%w: vacancy salary minimum cannot exceed maximum", ErrInvalidInput)
	default:
		parsed, err := url.ParseRequestURI(v.URL)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" ||
			(parsed.Scheme != "http" && parsed.Scheme != "https") {
			return fmt.Errorf("%w: vacancy URL is invalid", ErrInvalidInput)
		}
		return nil
	}
}

func (v Vacancy) Normalize() Vacancy {
	// Normalize — приводит текстовые поля и навыки к единому виду.
	v.Title = strings.TrimSpace(v.Title)
	v.Company = strings.TrimSpace(v.Company)
	v.Location = strings.TrimSpace(v.Location)
	v.Source = strings.ToLower(strings.TrimSpace(v.Source))
	v.URL = strings.TrimSpace(v.URL)
	v.Currency = strings.ToUpper(strings.TrimSpace(v.Currency))
	seen := make(map[string]struct{}, len(v.Skills))
	skills := make([]string, 0, len(v.Skills))
	for _, skill := range v.Skills {
		skill = strings.ToLower(strings.TrimSpace(skill))
		if skill != "" {
			if _, ok := seen[skill]; !ok {
				seen[skill] = struct{}{}
				skills = append(skills, skill)
			}
		}
	}
	v.Skills = skills
	return v
}

func (v Vacancy) IsRemoteCompatible(wantsRemote bool) bool {
	// IsRemoteCompatible — быстрый фильтр до дорогого semantic matching.
	if !wantsRemote {
		return true
	}
	location := strings.ToLower(strings.TrimSpace(v.Location))
	return v.Remote || location == "remote" || location == "anywhere"
}

func (v Vacancy) MatchesCriteria(criteria SearchCriteria) bool {
	if !v.IsRemoteCompatible(criteria.Remote) {
		return false
	}
	if criteria.MinSalary > 0 && v.SalaryMax > 0 && v.SalaryMax < criteria.MinSalary {
		return false
	}
	if criteria.EmploymentType != "" && !strings.EqualFold(criteria.EmploymentType, v.EmploymentType) {
		return false
	}
	if len(criteria.Locations) > 0 {
		found := false
		for _, location := range criteria.Locations {
			if strings.EqualFold(strings.TrimSpace(location), v.Location) {
				found = true
				break
			}
		}
		if !found && !v.Remote {
			return false
		}
	}
	return true
}
