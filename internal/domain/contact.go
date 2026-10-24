package domain

import "time"

type ContactType string

const (
	ContactGeneral ContactType = "general"
	ContactHR      ContactType = "hr"
)

type Contact struct {
	Email      string      `json:"email"`
	Type       ContactType `json:"type"`
	SourceURL  string      `json:"source_url"`
	Confidence float64     `json:"confidence"`
	FoundAt    time.Time   `json:"found_at"`
}
