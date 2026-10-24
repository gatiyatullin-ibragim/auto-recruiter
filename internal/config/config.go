package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Environment   string `env:"APP_ENV" envDefault:"development"`
	HTTPAddress   string `env:"HTTP_ADDR" envDefault:":8080"`
	DatabaseURL   string `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@localhost:5432/auto_recruiter?sslmode=disable"`
	RedisAddress  string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	LLMProvider   string `env:"LLM_PROVIDER" envDefault:"claude"`
	LLMBaseURL    string `env:"LLM_BASE_URL" envDefault:""`
	LLMModel      string `env:"LLM_MODEL" envDefault:""`
	HTTPTimeout   string `env:"HTTP_TIMEOUT" envDefault:"10s"`
	EncryptionKey string `env:"ENCRYPTION_KEY"`
	SMTPHost      string `env:"SMTP_HOST"`
	SMTPPort      int    `env:"SMTP_PORT" envDefault:"587"`
}

func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}
