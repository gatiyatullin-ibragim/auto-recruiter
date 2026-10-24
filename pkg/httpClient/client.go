package httpClient

import (
	"net/http"
	"time"
)

func New(timeout time.Duration) *http.Client {
	// New — создаёт HTTP-клиент с общим deadline на запрос.
	return &http.Client{Timeout: timeout}
}
