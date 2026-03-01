package materials

import (
	"context"

	repo "github.com/Ranoth/SUTBDB/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	GetMaterials(ctx context.Context) ([]repo.Material, error)
	GetMaterial(ctx context.Context, id int32) (repo.Material, error)
	GetMaterialsByName(ctx context.Context, name string) ([]repo.Material, error)
	CreateMaterial(ctx context.Context, material repo.CreateMaterialParams) (repo.Material, error)
	DeleteMaterial(ctx context.Context, id int32) error
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) *service {
	return &service{repo: repo}
}

func (s *service) GetMaterials(ctx context.Context) ([]repo.Material, error) {
	return s.repo.GetMaterials(ctx)
}

func (s *service) GetMaterial(ctx context.Context, id int32) (repo.Material, error) {
	return s.repo.GetMaterial(ctx, id)
}

func (s *service) GetMaterialsByName(ctx context.Context, name string) ([]repo.Material, error) {
	return s.repo.GetMaterialsByName(ctx, pgtype.Text{String: name, Valid: true})
}

func (s *service) CreateMaterial(ctx context.Context, params repo.CreateMaterialParams) (repo.Material, error) {
	id, err := s.repo.CreateMaterial(ctx, params)
	if err != nil {
		return repo.Material{}, err
	}
	return s.repo.GetMaterial(ctx, id)
}

func (s *service) DeleteMaterial(ctx context.Context, id int32) error {
	return s.repo.DeleteMaterial(ctx, id)
}
