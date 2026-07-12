package geticons

import (
	"context"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
)

type Service interface {
	GetCreatureIconById(ctx context.Context, id int32) ([]byte, error)
	GetMaterialIconById(ctx context.Context, id int32) ([]byte, error)
	GetRelicIconById(ctx context.Context, id int32) ([]byte, error)
	GetRaceIconById(ctx context.Context, id int32) ([]byte, error)
	GetClassIconById(ctx context.Context, id int32) ([]byte, error)
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) *service {
	return &service{repo: repo}
}

func (s *service) GetCreatureIconById(ctx context.Context, id int32) ([]byte, error) {
	return s.repo.GetCreatureIconById(ctx, id)
}

func (s *service) GetMaterialIconById(ctx context.Context, id int32) ([]byte, error) {
	return s.repo.GetMaterialIconById(ctx, id)
}

func (s *service) GetRelicIconById(ctx context.Context, id int32) ([]byte, error) {
	return s.repo.GetRelicIconById(ctx, id)
}

func (s *service) GetRaceIconById(ctx context.Context, id int32) ([]byte, error) {
	return s.repo.GetRaceIconById(ctx, id)
}

func (s *service) GetClassIconById(ctx context.Context, id int32) ([]byte, error) {
	return s.repo.GetClassIconById(ctx, id)
}
