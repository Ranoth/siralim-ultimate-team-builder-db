package specializations

import (
	"context"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	GetSpecializations(ctx context.Context) ([]repo.Specialization, error)
	GetSpecialization(ctx context.Context, id int32) (repo.Specialization, error)
	GetSpecializationsByName(ctx context.Context, name string) ([]repo.Specialization, error)
	CreateSpecialization(ctx context.Context, params repo.CreateSpecializationParams) (repo.Specialization, error)
	DeleteSpecialization(ctx context.Context, id int32) error
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) *service {
	return &service{repo: repo}
}

func (s *service) GetSpecializations(ctx context.Context) ([]repo.Specialization, error) {
	return s.repo.GetSpecializations(ctx)
}

func (s *service) GetSpecialization(ctx context.Context, id int32) (repo.Specialization, error) {
	return s.repo.GetSpecialization(ctx, id)
}

func (s *service) GetSpecializationsByName(ctx context.Context, name string) ([]repo.Specialization, error) {
	return s.repo.GetSpecializationsByName(ctx, pgtype.Text{String: name, Valid: true})
}

func (s *service) CreateSpecialization(ctx context.Context, params repo.CreateSpecializationParams) (repo.Specialization, error) {
	id, err := s.repo.CreateSpecialization(ctx, params)
	if err != nil {
		return repo.Specialization{}, err
	}
	return s.repo.GetSpecialization(ctx, id)
}

func (s *service) DeleteSpecialization(ctx context.Context, id int32) error {
	return s.repo.DeleteSpecialization(ctx, id)
}
