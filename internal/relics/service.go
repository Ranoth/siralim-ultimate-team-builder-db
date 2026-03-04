package relics

import (
	"context"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	GetRelics(ctx context.Context) ([]repo.GetRelicsRow, error)
	GetRelic(ctx context.Context, id int32) (repo.GetRelicRow, error)
	GetRelicsByName(ctx context.Context, name string) ([]repo.GetRelicsByNameRow, error)
	CreateRelic(ctx context.Context, params repo.CreateRelicParams) (repo.GetRelicRow, error)
	DeleteRelic(ctx context.Context, id int32) error
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) *service {
	return &service{repo: repo}
}

func (s *service) GetRelics(ctx context.Context) ([]repo.GetRelicsRow, error) {
	return s.repo.GetRelics(ctx)
}

func (s *service) GetRelic(ctx context.Context, id int32) (repo.GetRelicRow, error) {
	return s.repo.GetRelic(ctx, id)
}

func (s *service) GetRelicsByName(ctx context.Context, name string) ([]repo.GetRelicsByNameRow, error) {
	return s.repo.GetRelicsByName(ctx, pgtype.Text{String: name, Valid: true})
}

func (s *service) CreateRelic(ctx context.Context, params repo.CreateRelicParams) (repo.GetRelicRow, error) {
	// Convert repo.GetRelicRow to repo.Relic if needed
	// or adjust your query to return the correct type
	_, err := s.repo.CreateRelic(ctx, params)
	if err != nil {
		return repo.GetRelicRow{}, err
	}

	// Map GetRelicRow to Relic if necessary
	return repo.GetRelicRow{
		// Map fields from relic (GetRelicRow) to repo.GetRelicRow
	}, nil
}

func (s *service) DeleteRelic(ctx context.Context, id int32) error {
	return s.repo.DeleteRelic(ctx, id)
}
