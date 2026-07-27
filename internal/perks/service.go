package perks

import (
	"context"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	GetPerks(ctx context.Context) ([]repo.Perk, error)
	GetPerk(ctx context.Context, id int32) (repo.Perk, error)
	GetPerksByName(ctx context.Context, name string) ([]repo.Perk, error)
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) *service {
	return &service{repo: repo}
}

func (s *service) GetPerks(ctx context.Context) ([]repo.Perk, error) {
	return s.repo.GetPerks(ctx)
}

func (s *service) GetPerk(ctx context.Context, id int32) (repo.Perk, error) {
	return s.repo.GetPerk(ctx, id)
}

func (s *service) GetPerksByName(ctx context.Context, name string) ([]repo.Perk, error) {
	return s.repo.GetPerksByName(ctx, pgtype.Text{String: name, Valid: true})
}
