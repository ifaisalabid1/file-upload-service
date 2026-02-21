package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	AWS      AWSConfig
	App      AppConfig
}

type ServerConfig struct {
	Port         string
	Environment  string
	LogLevel     slog.Level
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host              string
	Port              int
	User              string
	Password          string
	DBName            string
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
}

type AWSConfig struct {
	Region    string
	AccessKey string
	SecretKey string
	S3Bucket  string
	S3BaseURL string
}

type AppConfig struct {
	MaxFileSize      int64
	AllowedFileTypes []string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Environment:  getEnv("ENVIRONMENT", "development"),
			LogLevel:     parseLogLevel(getEnv("LOG_LEVEL", "info")),
			ReadTimeout:  parseDuration(getEnv("READ_TIMEOUT", "15s"), 15*time.Second),
			WriteTimeout: parseDuration(getEnv("WRITE_TIMEOUT", "15s"), 15*time.Second),
			IdleTimeout:  parseDuration(getEnv("IDLE_TIMEOUT", "60s"), 60*time.Second),
		},

		Database: DatabaseConfig{
			Host:              getEnv("DB_HOST", "localhost"),
			Port:              parseInt(getEnv("DB_PORT", "5432")),
			User:              getEnv("DB_USER", "postgres"),
			Password:          getEnv("DB_PASSWORD", "postgres"),
			DBName:            getEnv("DB_NAME", "fileupload"),
			MaxConns:          int32(parseInt(getEnv("DB_MAX_CONNS", "10"))),
			MinConns:          int32(parseInt(getEnv("DB_MIN_CONNS", "2"))),
			MaxConnLifetime:   parseDuration(getEnv("DB_MAX_CONN_LIFETIME", "1h"), time.Hour),
			MaxConnIdleTime:   parseDuration(getEnv("DB_MAX_CONN_IDLE_TIME", "30m"), 30*time.Minute),
			HealthCheckPeriod: parseDuration(getEnv("DB_HEALTH_CHECK_PERIOD", "1m"), time.Minute),
		},

		AWS: AWSConfig{
			Region:    getEnv("AWS_REGION", "us-east-1"),
			AccessKey: getEnv("AWS_ACCESS_KEY", ""),
			SecretKey: getEnv("AWS_SECRET_KEY", ""),
			S3Bucket:  getEnv("S3_BUCKET", ""),
			S3BaseURL: getEnv("S3_BASE_URL", ""),
		},

		App: AppConfig{
			MaxFileSize: parseInt64(getEnv("MAX_FILE_SIZE", "10_485_760")),
			AllowedFileTypes: []string{
				"image/jpeg",
				"image/jpg",
				"image/png",
				"image/webp",
				"image/gif",
				"application/pdf",
				"text/plain",
			},
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	var missing []string

	if c.Database.Host == "" {
		missing = append(missing, "DB_HOST is required")
	}
	if c.Database.User == "" {
		missing = append(missing, "DB_USER is required")
	}
	if c.Database.Password == "" {
		missing = append(missing, "DB_Password is required")
	}
	if c.Database.DBName == "" {
		missing = append(missing, "DB_NAME is required")
	}
	if c.AWS.Region == "" {
		missing = append(missing, "AWS_REGION is required")
	}
	if c.AWS.S3Bucket == "" {
		missing = append(missing, "S3_BUCKET is required")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables %s", strings.Join(missing, ", "))
	}

	return nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.DBName)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func parseInt(value string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return i
}

func parseInt64(value string) int64 {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}

	return i
}

func parseDuration(value string, defaultValue time.Duration) time.Duration {
	d, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return d
}

func parseLogLevel(level string) slog.Level {

	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
