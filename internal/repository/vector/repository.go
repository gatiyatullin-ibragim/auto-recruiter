package vector

import "context"

type Repository interface {
	// Repository — порт векторного хранилища для semantic search.
	Upsert(context.Context, string, []float32, map[string]string) error
	Search(context.Context, []float32, int) ([]string, error)
}
