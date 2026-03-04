package dbseeder

import (
	"context"
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

func (s *Seeder) Seed() {
	config := newSeederConfig()
	jsonParser := newJSONParser(s.logger, config)
	normalizer := newNormalizer(jsonParser, s.logger, config)
	inserter := newInserter(s.logger, config, s.queries)

	jsonParser.parseAndStore()
	normalizer.normalize()
	inserter.insert()
}
