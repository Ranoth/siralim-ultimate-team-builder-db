package creatures

import (
	"context"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	GetCreatures(ctx context.Context) ([]repo.CreaturesView, error)
	GetCreature(ctx context.Context, id int32) (repo.CreaturesView, error)
	GetCreaturesByName(ctx context.Context, name string) ([]repo.CreaturesView, error)
	GetCreaturesByTraitName(ctx context.Context, name string) ([]repo.CreaturesView, error)
	GetCreaturesByClassName(ctx context.Context, name string) ([]repo.CreaturesView, error)
	GetCreaturesByRaceName(ctx context.Context, name string) ([]repo.CreaturesView, error)
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) *service {
	return &service{repo: repo}
}

func (s *service) GetCreatures(ctx context.Context) ([]repo.CreaturesView, error) {
	return s.repo.GetCreatures(ctx)
}

func (s *service) GetCreature(ctx context.Context, id int32) (repo.CreaturesView, error) {
	return s.repo.GetCreature(ctx, id)
}

func (s *service) GetCreaturesByName(ctx context.Context, name string) ([]repo.CreaturesView, error) {
	return s.repo.GetCreaturesByName(ctx, pgtype.Text{String: name, Valid: true})
}

func (s *service) GetCreaturesByTraitName(ctx context.Context, name string) ([]repo.CreaturesView, error) {
	return s.repo.GetCreaturesByTraitName(ctx, pgtype.Text{String: name, Valid: true})
}

func (s *service) GetCreaturesByClassName(ctx context.Context, name string) ([]repo.CreaturesView, error) {
	return s.repo.GetCreaturesByClassName(ctx, pgtype.Text{String: name, Valid: true})
}

func (s *service) GetCreaturesByRaceName(ctx context.Context, name string) ([]repo.CreaturesView, error) {
	return s.repo.GetCreaturesByRaceName(ctx, pgtype.Text{String: name, Valid: true})
}
