package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/PalDeepak810/evenly/backend/internal/config"
	"github.com/PalDeepak810/evenly/backend/internal/middleware"
	"github.com/PalDeepak810/evenly/backend/internal/platform/database"
	"github.com/PalDeepak810/evenly/backend/internal/platform/logger"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// 1. Load application configuration
	cfg, err := config.Load()

	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// 2. Initialize structured logger.
	appLogger := logger.New(cfg.AppEnv)

	appLogger.Info("Starting Evenly backend")

	// 3. Initialize PostgresSQL connection pool.
	dbCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	pool, err := database.NewPostgresPool(
		dbCtx,
		cfg.DatabaseURL,
	)

	if err != nil {
		return fmt.Errorf("Initialize database: %w", err)
	}
	defer pool.Close()

	appLogger.Info("PostgreSQL connected sucessfully")

	// 4. Initialize HTTP router.
	mux := http.NewServeMux()

	// 5. Register middleware.
	var handler http.Handler = mux

	handler = middleware.RequestLogger(appLogger)(handler)
	handler = middleware.RequestID(handler)

	// 6. Initialize HTTP Server.

	server := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	appLogger.Info(
		"Starting server",
		slog.String("address", cfg.HTTPAddr),
	)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("Http Server failed: %w", err)
	}

	return nil
}

var _ = os.Stdout
