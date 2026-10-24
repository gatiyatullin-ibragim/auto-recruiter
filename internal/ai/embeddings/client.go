package embeddings

import "context"

type Client interface {
	// Client — абстракция получения числового вектора для текста.
	Embed(context.Context, string) ([]float32, error)
}
