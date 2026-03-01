package races

import (
	"context"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	GetRaces(ctx context.Context) ([]repo.Race, error)
	GetRace(ctx context.Context, id int32) (repo.Race, error)
	GetRacesByName(ctx context.Context, name string) ([]repo.Race, error)
	GetRacesByTraitName(ctx context.Context, traitName string) ([]repo.Race, error)
	GetRacesByClassName(ctx context.Context, className string) ([]repo.Race, error)
	GetRacesByCreatureName(ctx context.Context, creatureName string) ([]repo.Race, error)
	CreateRace(ctx context.Context, params repo.CreateRaceParams) (repo.Race, error)
	DeleteRace(ctx context.Context, id int32) error
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) *service {
	return &service{repo: repo}
}

func (s *service) GetRaces(ctx context.Context) ([]repo.Race, error) {
	return s.repo.GetRaces(ctx)
}

func (s *service) GetRace(ctx context.Context, id int32) (repo.Race, error) {
	return s.repo.GetRace(ctx, id)
}

func (s *service) GetRacesByName(ctx context.Context, name string) ([]repo.Race, error) {
	return s.repo.GetRacesByName(ctx, pgtype.Text{String: name, Valid: true})
}

func (s *service) GetRacesByTraitName(ctx context.Context, traitName string) ([]repo.Race, error) {
	return s.repo.GetRacesByTraitName(ctx, pgtype.Text{String: traitName, Valid: true})
}

func (s *service) GetRacesByClassName(ctx context.Context, className string) ([]repo.Race, error) {
	return s.repo.GetRacesByClassName(ctx, pgtype.Text{String: className, Valid: true})
}

func (s *service) GetRacesByCreatureName(ctx context.Context, creatureName string) ([]repo.Race, error) {
	return s.repo.GetRacesByCreatureName(ctx, pgtype.Text{String: creatureName, Valid: true})
}

func (s *service) CreateRace(ctx context.Context, params repo.CreateRaceParams) (repo.Race, error) {
	id, err := s.repo.CreateRace(ctx, params)
	if err != nil {
		return repo.Race{}, err
	}
	return s.repo.GetRace(ctx, id)
}

func (s *service) DeleteRace(ctx context.Context, id int32) error {
	return s.repo.DeleteRace(ctx, id)
}
