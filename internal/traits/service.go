package traits

import (
	"context"

	repo "github.com/Ranoth/SUTBDB/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	GetTraits(ctx context.Context) ([]repo.Trait, error)
	GetTrait(ctx context.Context, id int32) (repo.Trait, error)
	GetTraitsByName(ctx context.Context, name string) ([]repo.Trait, error)
	CreateTrait(ctx context.Context, trait repo.CreateTraitParams) (repo.Trait, error)
	DeleteTrait(ctx context.Context, id int32) error
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) *service {
	return &service{repo: repo}
}

func (s *service) GetTraits(ctx context.Context) ([]repo.Trait, error) {
	return s.repo.GetTraits(ctx)
}

func (s *service) GetTrait(ctx context.Context, id int32) (repo.Trait, error) {
	return s.repo.GetTrait(ctx, id)
}

func (s *service) GetTraitsByName(ctx context.Context, name string) ([]repo.Trait, error) {
	return s.repo.GetTraitsByName(ctx, pgtype.Text{String: name, Valid: true})
}

func (s *service) CreateTrait(ctx context.Context, params repo.CreateTraitParams) (repo.Trait, error) {
	id, err := s.repo.CreateTrait(ctx, params)
	if err != nil {
		return repo.Trait{}, err
	}
	return s.repo.GetTrait(ctx, id)
}

func (s *service) DeleteTrait(ctx context.Context, id int32) error {
	return s.repo.DeleteTrait(ctx, id)
}
