package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chiMid "github.com/go-chi/chi/v5/middleware"
	"github.com/ifaisalabid1/file-upload-service/internal/config"
	"github.com/ifaisalabid1/file-upload-service/internal/handler"
	"github.com/ifaisalabid1/file-upload-service/internal/logger"
	"github.com/ifaisalabid1/file-upload-service/internal/middleware"
	"github.com/ifaisalabid1/file-upload-service/internal/repository"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	logger := logger.New(cfg.Env, cfg.LogLevel)
	slog.SetDefault(logger)

	logger.Info("starting service...", slog.Int("port", cfg.ServerPort), slog.String("log_level", cfg.LogLevel))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := repository.NewPool(ctx, cfg.DatabaseURL, logger)
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	logger.Info("connected to postgres")

	healthHandler := &handler.HealthHandler{
		DB: pool,
	}

	r := chi.NewRouter()
	r.Use(chiMid.RequestID)
	r.Use(chiMid.RealIP)
	r.Use(middleware.Logger(logger))
	r.Use(chiMid.Recoverer)

	r.Get("/health", healthHandler.HealthCheckHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Info("shutting down", slog.String("signal", sig.String()))

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("forced shutdown", slog.String("error", err.Error()))
	}

	logger.Info("server stopped gracefully")
}
