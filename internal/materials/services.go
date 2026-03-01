package materials

import (
	"context"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	GetMaterials(ctx context.Context) ([]repo.GetMaterialsRow, error)
	GetMaterial(ctx context.Context, id int32) (repo.GetMaterialRow, error)
	GetMaterialsByName(ctx context.Context, name string) ([]repo.GetMaterialsRow, error)
	CreateMaterial(ctx context.Context, params repo.CreateMaterialParams) (repo.Material, error)
	DeleteMaterial(ctx context.Context, id int32) error
	CreateMaterialStat(ctx context.Context, params repo.CreateMaterialStatParams) (int32, error)
	DeleteMaterialStat(ctx context.Context, id int32) error
	CreateStat(ctx context.Context, statType repo.StatType) (int32, error)
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

func (s *service) GetMaterialsByName(ctx context.Context, name string) ([]repo.GetMaterialsRow, error) {
	materials, err := s.repo.GetMaterialsByName(ctx, pgtype.Text{String: name, Valid: true})
	if err != nil {
		return nil, err
	}
	results := make([]repo.GetMaterialsRow, len(materials))
	for i, _ := range materials {
		results[i] = repo.GetMaterialsRow{}
	}
	return results, nil
}

func (s *service) CreateMaterial(ctx context.Context, params repo.CreateMaterialParams) (repo.Material, error) {
	createRows, err := s.repo.CreateMaterial(ctx, params)
	if err != nil {
		return repo.Material{}, err
	}
	materialRows, err := s.repo.GetMaterial(ctx, createRows.ID)
	if err != nil {
		return repo.Material{}, err
	}
	return repo.Material{
		ID:          materialRows.ID,
		Name:        materialRows.Name,
		Description: materialRows.Description,
		Icon:        materialRows.Icon,
		Type:        materialRows.Type,
	}, nil
}

func (s *service) DeleteMaterial(ctx context.Context, id int32) error {
	return s.repo.DeleteMaterial(ctx, id)
}

// func UpdateMaterialStats(ctx context.Context, repo repo.Querier, params repo.UpdateMaterialStatParams) (repo.UpdateMaterialStatParams, error) {
// 	err := repo.UpdateMaterialStat(ctx, params)
// 	if err != nil {
// 		return params, err
// 	}
// 	return params, nil
// }

func (s *service) CreateMaterialStat(ctx context.Context, params repo.CreateMaterialStatParams) (int32, error) {
	stat, err := s.repo.CreateMaterialStat(ctx, params)
	if err != nil {
		return 0, err
	}
	return stat, nil
}

func (s *service) DeleteMaterialStat(ctx context.Context, id int32) error {
	return s.repo.DeleteMaterialStat(ctx, id)
}

func (s *service) CreateStat(ctx context.Context, statType repo.StatType) (int32, error) {
	id, err := s.repo.CreateStat(ctx, statType)
	if err != nil {
		return 0, err
	}
	return id, nil
}
