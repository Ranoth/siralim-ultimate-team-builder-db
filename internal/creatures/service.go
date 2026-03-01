package creatures

import (
	"context"

	repo "github.com/Ranoth/SUTBDB/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	GetCreatures(ctx context.Context) ([]repo.Creature, error)
	GetCreature(ctx context.Context, id int32) (repo.Creature, error)
	GetCreaturesByName(ctx context.Context, name string) ([]repo.Creature, error)
	GetCreaturesByTraitName(ctx context.Context, name string) ([]repo.Creature, error)
	GetCreaturesByClassName(ctx context.Context, name string) ([]repo.Creature, error)
	GetCreaturesByRaceName(ctx context.Context, name string) ([]repo.Creature, error)
	CreateCreature(ctx context.Context, params repo.CreateCreatureParams) (repo.Creature, error)
	DeleteCreature(ctx context.Context, id int32) error
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) *service {
	return &service{repo: repo}
}

func (s *service) GetCreatures(ctx context.Context) ([]repo.Creature, error) {
	return s.repo.GetCreatures(ctx)
}

func (s *service) GetCreature(ctx context.Context, id int32) (repo.Creature, error) {
	return s.repo.GetCreature(ctx, id)
}

func (s *service) GetCreaturesByName(ctx context.Context, name string) ([]repo.Creature, error) {
	return s.repo.GetCreaturesByName(ctx, pgtype.Text{String: name, Valid: true})
}

func (s *service) GetCreaturesByTraitName(ctx context.Context, name string) ([]repo.Creature, error) {
	return s.repo.GetCreaturesByTraitName(ctx, pgtype.Text{String: name, Valid: true})
}

func (s *service) GetCreaturesByClassName(ctx context.Context, name string) ([]repo.Creature, error) {
	return s.repo.GetCreaturesByClassName(ctx, pgtype.Text{String: name, Valid: true})
}

func (s *service) GetCreaturesByRaceName(ctx context.Context, name string) ([]repo.Creature, error) {
	return s.repo.GetCreaturesByRaceName(ctx, pgtype.Text{String: name, Valid: true})
}

func (s *service) CreateCreature(ctx context.Context, params repo.CreateCreatureParams) (repo.Creature, error) {
	id, err := s.repo.CreateCreature(ctx, params)
	if err != nil {
		return repo.Creature{}, err
	}
	return s.repo.GetCreature(ctx, id)
}

func (s *service) DeleteCreature(ctx context.Context, id int32) error {
	return s.repo.DeleteCreature(ctx, id)
}
