package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ifaisalabid1/file-upload-service/internal/config"
	"github.com/ifaisalabid1/file-upload-service/internal/logger"
	"github.com/ifaisalabid1/file-upload-service/internal/repository"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", slog.String("error", err.Error()))
	}

	logger := logger.New(cfg.Env, cfg.LogLevel)
	slog.SetDefault(logger)

	pool, err := repository.NewPool(ctx, cfg.DatabaseURL, logger)
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	logger.Info("connected to postgres")
}
