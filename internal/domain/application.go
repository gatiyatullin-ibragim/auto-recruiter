package domain

import "time"

type ApplicationStatus string

const (
	// Статусы описывают жизненный цикл отклика кандидата.
	ApplicationDraft     ApplicationStatus = "draft"
	ApplicationGenerated ApplicationStatus = "generated"
	ApplicationQueued    ApplicationStatus = "queued"
	ApplicationSent      ApplicationStatus = "sent"
	ApplicationRejected  ApplicationStatus = "rejected"
	ApplicationFailed    ApplicationStatus = "failed"
	ApplicationManual    ApplicationStatus = "manual_required"
)

type Application struct {
	// Application связывает кандидата, вакансию и сопроводительное письмо.
	ID          string            `json:"id"`
	CandidateID string            `json:"candidate_id"`
	VacancyID   string            `json:"vacancy_id"`
	CoverLetter string            `json:"cover_letter"`
	Status      ApplicationStatus `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func (a Application) CanTransitionTo(next ApplicationStatus) bool {
	switch a.Status {
	case ApplicationDraft:
		return next == ApplicationGenerated
	case ApplicationGenerated:
		return next == ApplicationQueued || next == ApplicationSent || next == ApplicationRejected
	case ApplicationQueued:
		return next == ApplicationSent || next == ApplicationFailed || next == ApplicationManual
	case ApplicationSent:
		return next == ApplicationRejected
	case ApplicationFailed, ApplicationManual:
		return next == ApplicationQueued || next == ApplicationRejected
	default:
		return false
	}
}
