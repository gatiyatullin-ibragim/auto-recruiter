package domain

import "errors"

var (
	ErrNotFound            = errors.New("resource not found")
	ErrInvalidInput        = errors.New("invalid input")
	ErrExternalUnavailable = errors.New("external source unavailable")
	ErrConflict            = errors.New("resource conflict")
	ErrCacheMiss           = errors.New("cache miss")
)
