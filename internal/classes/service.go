package classes

import (
	"context"

	repo "github.com/Ranoth/SUTBDB/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	GetClasses(ctx context.Context) ([]repo.Class, error)
	GetClass(ctx context.Context, id int32) (repo.Class, error)
	GetClassesByName(ctx context.Context, name string) ([]repo.Class, error)
	CreateClass(ctx context.Context, params repo.CreateClassParams) (repo.Class, error)
	DeleteClass(ctx context.Context, id int32) error
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) *service {
	return &service{repo: repo}
}

func (s *service) GetClasses(ctx context.Context) ([]repo.Class, error) {
	return s.repo.GetClasses(ctx)
}

func (s *service) GetClass(ctx context.Context, id int32) (repo.Class, error) {
	return s.repo.GetClass(ctx, id)
}

func (s *service) GetClassesByName(ctx context.Context, name string) ([]repo.Class, error) {
	return s.repo.GetClassesByName(ctx, pgtype.Text{String: name, Valid: true})
}

func (s *service) CreateClass(ctx context.Context, params repo.CreateClassParams) (repo.Class, error) {
	id, err := s.repo.CreateClass(ctx, params)
	if err != nil {
		return repo.Class{}, err
	}
	return s.repo.GetClass(ctx, id)
}

func (s *service) DeleteClass(ctx context.Context, id int32) error {
	return s.repo.DeleteClass(ctx, id)
}
