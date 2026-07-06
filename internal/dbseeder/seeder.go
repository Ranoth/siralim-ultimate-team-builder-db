package dbseeder

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
)

type Seeder struct {
	conn    *pgxpool.Pool
	queries *repo.Queries
	logger  *slog.Logger
}

func NewSeeder(conn *pgxpool.Pool, queries *repo.Queries, logger *slog.Logger) *Seeder {
	return &Seeder{
		conn:    conn,
		queries: queries,
		logger:  logger,
	}
}

func (s *Seeder) CheckIfSeeded() (bool, error) {
	count, err := s.queries.GetStatsCount(context.Background())
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Seeder) Seed() error {
	config := newSeederConfig()
	jsonParser := newJSONParser(s.logger, config)
	normalizer := newNormalizer(jsonParser, s.logger, config)

	if err := jsonParser.parseAndStore(); err != nil {
		return fmt.Errorf("error parsing JSON files: %w", err)
	}
	if err := normalizer.normalize(); err != nil {
		return fmt.Errorf("error normalizing data: %w", err)
	}

	ctx := context.Background()
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	inserter := newInserter(s.logger, config, s.queries.WithTx(tx))
	if err := inserter.insert(); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return fmt.Errorf("insert failed: %w; rollback failed: %v", err, rollbackErr)
		}
		return fmt.Errorf("error inserting data into database: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
