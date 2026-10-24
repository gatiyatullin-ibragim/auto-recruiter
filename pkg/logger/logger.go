package logger

import (
	"log/slog"
	"os"
)

func New(environment string) *slog.Logger {
	// New — выбирает text logging для разработки и JSON для production.
	var handler slog.Handler = slog.NewJSONHandler(os.Stdout, nil)
	if environment == "development" {
		handler = slog.NewTextHandler(os.Stdout, nil)
	}
	return slog.New(handler)
}
