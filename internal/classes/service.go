package classes

import (
	"context"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	GetClasses(ctx context.Context) ([]repo.ClassesView, error)
	GetClass(ctx context.Context, id int32) (repo.ClassesView, error)
	GetClassesByName(ctx context.Context, name string) ([]repo.ClassesView, error)
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) *service {
	return &service{repo: repo}
}

func (s *service) GetClasses(ctx context.Context) ([]repo.ClassesView, error) {
	return s.repo.GetClasses(ctx)
}

func (s *service) GetClass(ctx context.Context, id int32) (repo.ClassesView, error) {
	return s.repo.GetClass(ctx, id)
}

func (s *service) GetClassesByName(ctx context.Context, name string) ([]repo.ClassesView, error) {
	return s.repo.GetClassesByName(ctx, pgtype.Text{String: name, Valid: true})
}
