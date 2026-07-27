package relics

import (
	"context"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	GetRelics(ctx context.Context) ([]repo.RelicsView, error)
	GetRelic(ctx context.Context, id int32) (repo.RelicsView, error)
	GetRelicsByName(ctx context.Context, name string) ([]repo.RelicsView, error)
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) *service {
	return &service{repo: repo}
}

func (s *service) GetRelics(ctx context.Context) ([]repo.RelicsView, error) {
	return s.repo.GetRelics(ctx)
}

func (s *service) GetRelic(ctx context.Context, id int32) (repo.RelicsView, error) {
	return s.repo.GetRelic(ctx, id)
}

func (s *service) GetRelicsByName(ctx context.Context, name string) ([]repo.RelicsView, error) {
	return s.repo.GetRelicsByName(ctx, pgtype.Text{String: name, Valid: true})
}
