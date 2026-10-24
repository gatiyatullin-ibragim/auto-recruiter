package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXCandidateRepository struct{ pool *pgxpool.Pool }

func NewPGXCandidateRepository(pool *pgxpool.Pool) *PGXCandidateRepository {
	return &PGXCandidateRepository{pool: pool}
}

func (r *PGXCandidateRepository) GetByID(ctx context.Context, id string) (domain.Candidate, error) {
	var candidate domain.Candidate
	var skills, criteria []byte
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, email, resume_text, skills, criteria, created_at, updated_at
		FROM candidates WHERE id = $1`, id).Scan(
		&candidate.ID, &candidate.Name, &candidate.Email, &candidate.ResumeText,
		&skills, &criteria, &candidate.CreatedAt, &candidate.UpdatedAt)
	if err != nil {
		return domain.Candidate{}, fmt.Errorf("get candidate: %w", err)
	}
	if err := json.Unmarshal(skills, &candidate.Skills); err != nil {
		return domain.Candidate{}, fmt.Errorf("decode candidate skills: %w", err)
	}
	if err := json.Unmarshal(criteria, &candidate.Criteria); err != nil {
		return domain.Candidate{}, fmt.Errorf("decode candidate criteria: %w", err)
	}
	return candidate, nil
}

func (r *PGXCandidateRepository) Save(ctx context.Context, candidate domain.Candidate) error {
	skills, err := json.Marshal(candidate.Skills)
	if err != nil {
		return fmt.Errorf("encode candidate skills: %w", err)
	}
	criteria, err := json.Marshal(candidate.Criteria)
	if err != nil {
		return fmt.Errorf("encode candidate criteria: %w", err)
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO candidates (id, name, email, resume_text, skills, criteria)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name, email = EXCLUDED.email,
			resume_text = EXCLUDED.resume_text, skills = EXCLUDED.skills,
			criteria = EXCLUDED.criteria, updated_at = now()`,
		candidate.ID, candidate.Name, candidate.Email, candidate.ResumeText, skills, criteria)
	if err != nil {
		return fmt.Errorf("save candidate: %w", err)
	}
	return nil
}
