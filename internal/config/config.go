package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Env           string        `env:"ENV" envDefault:"dev"`
	ServerPort    int           `env:"SERVER_PORT"   envDefault:"8080"`
	DatabaseURL   string        `env:"DATABASE_URL"  envDefault:"postgres://user:pass@localhost:5432/uploaddb?sslmode=disable"`
	LogLevel      string        `env:"LOG_LEVEL"     envDefault:"info"`
	RateLimit     float64       `env:"RATE_LIMIT"    envDefault:"10"`
	RateBurst     int           `env:"RATE_BURST"    envDefault:"20"`
	Auth0Domain   string        `env:"AUTH0_DOMAIN"`
	Auth0Audience string        `env:"AUTH0_AUDIENCE"`
	AWSRegion     string        `env:"AWS_REGION"`
	AWSBucket     string        `env:"AWS_S3_BUCKET"`
	ReadTimeout   time.Duration `env:"READ_TIMEOUT"  envDefault:"5s"`
	WriteTimeout  time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout   time.Duration `env:"IDLE_TIMEOUT"  envDefault:"60s"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	return &cfg, nil
}
