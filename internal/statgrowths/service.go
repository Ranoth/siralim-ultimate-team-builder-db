package statgrowths

import (
	"context"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
)

type Service interface {
	GetStatGrowthsByCreatureId(ctx context.Context, id int32) ([]repo.GetStatGrowthsByCreatureIdRow, error)
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) *service {
	return &service{repo: repo}
}

func (s *service) GetStatGrowthsByCreatureId(ctx context.Context, creatureId int32) ([]repo.GetStatGrowthsByCreatureIdRow, error) {
	return s.repo.GetStatGrowthsByCreatureId(ctx, creatureId)
}
