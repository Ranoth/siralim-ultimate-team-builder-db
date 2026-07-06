package materials

import (
	"context"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	GetMaterials(ctx context.Context) ([]repo.GetMaterialsRow, error)
	GetMaterial(ctx context.Context, id int32) (repo.GetMaterialRow, error)
	GetMaterialsByName(ctx context.Context, name string) ([]repo.GetMaterialsByNameRow, error)
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) *service {
	return &service{repo: repo}
}

func (s *service) GetMaterials(ctx context.Context) ([]repo.GetMaterialsRow, error) {
	return s.repo.GetMaterials(ctx)
}

func (s *service) GetMaterial(ctx context.Context, id int32) (repo.GetMaterialRow, error) {
	return s.repo.GetMaterial(ctx, id)
}

func (s *service) GetMaterialsByName(ctx context.Context, name string) ([]repo.GetMaterialsByNameRow, error) {
	return s.repo.GetMaterialsByName(ctx, pgtype.Text{String: name, Valid: true})
}
