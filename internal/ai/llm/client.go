package llm

import "context"

type Client interface {
	// Client — абстракция текстовой генерации; сервисы не знают провайдера.
	Complete(context.Context, string) (string, error)
}
