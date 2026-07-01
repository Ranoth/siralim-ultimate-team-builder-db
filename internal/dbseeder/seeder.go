package dbseeder

import (
	"context"
	"fmt"
	"log/slog"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
)

type Seeder struct {
	queries *repo.Queries
	logger  *slog.Logger
}

func NewSeeder(queries *repo.Queries, logger *slog.Logger) *Seeder {
	return &Seeder{
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
	inserter := newInserter(s.logger, config, s.queries)

	if err := jsonParser.parseAndStore(); err != nil {
		return fmt.Errorf("error parsing JSON files: %w", err)
	}
	if err := normalizer.normalize(); err != nil {
		return fmt.Errorf("error normalizing data: %w", err)
	}
	if err := inserter.insert(); err != nil {
		return fmt.Errorf("error inserting data into database: %w", err)
	}
	return nil
}
