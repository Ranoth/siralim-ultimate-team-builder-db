package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
	"github.com/Ranoth/siralim-ultimate-team-builder-db/internal/dbseeder"
	"github.com/Ranoth/siralim-ultimate-team-builder-db/internal/env"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	cfg := config{
		addr: env.GetString("ADDRESS", ":8080"),
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=postgres user=postgres password=postgres dbname=sutbdb sslmode=disable"),
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	logger.Info("Connected to database", "dsn", cfg.db.dsn)

	// Check if database is already seeded before loading JSON
	queries := repo.New(conn)
	seeder := dbseeder.NewSeeder(queries)

	isSeeded, err := seeder.IsAlreadySeeded(ctx)
	if err != nil {
		logger.Debug("Error checking if database is seeded", "error", err)
	}

	if !isSeeded {
		// Load and process JSON data only if database is not seeded
		dbseeder.Run()
		transformedData := dbseeder.GetCorrelatedTables()

		if transformedData != nil {
			if err := seeder.SeedDatabase(ctx, transformedData); err != nil {
				logger.Error("Failed to seed database", "error", err)
				return
			}
		}
	} else {
		logger.Info("Database already seeded, skipping data loading and seeding")
	}

	api := application{
		config: cfg,
		db:     conn,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := api.run(api.mount()); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error", "error", err)
		}
	}()

	<-quit

	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := api.shutdown(shutdownCtx); err != nil {
		logger.Error("Error during shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Server shut down gracefully")
}
