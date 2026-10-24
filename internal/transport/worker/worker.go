package worker

import "context"

// Job — единица фоновой работы, запускаемая runner-ом.
type Job interface{ Run(context.Context) error }

// Runner — последовательный исполнитель фоновых job-ов.
type Runner struct{ jobs []Job }

// NewRunner — принимает job-ы через variadic dependency injection.
func NewRunner(jobs ...Job) *Runner { return &Runner{jobs: jobs} }

// Run — выполняет job-ы и останавливается на первой ошибке.
func (r *Runner) Run(ctx context.Context) error {
	for _, job := range r.jobs {
		if err := job.Run(ctx); err != nil {
			return err
		}
	}
	return nil
}
