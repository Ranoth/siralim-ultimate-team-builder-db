package artifacts

import (
	"context"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	GetArtifacts(ctx context.Context) ([]repo.Artifact, error)
	GetArtifact(ctx context.Context, id int32) (repo.Artifact, error)
	GetArtifactsByName(ctx context.Context, name string) ([]repo.Artifact, error)
	CreateArtifact(ctx context.Context, params repo.CreateArtifactParams) (repo.Artifact, error)
	DeleteArtifact(ctx context.Context, id int32) error
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) *service {
	return &service{repo: repo}
}

func (s *service) GetArtifacts(ctx context.Context) ([]repo.Artifact, error) {
	return s.repo.GetArtifacts(ctx)
}

func (s *service) GetArtifact(ctx context.Context, id int32) (repo.Artifact, error) {
	return s.repo.GetArtifact(ctx, id)
}

func (s *service) GetArtifactsByName(ctx context.Context, name string) ([]repo.Artifact, error) {
	return s.repo.GetArtifactsByName(ctx, pgtype.Text{String: name, Valid: true})
}

func (s *service) CreateArtifact(ctx context.Context, params repo.CreateArtifactParams) (repo.Artifact, error) {
	id, err := s.repo.CreateArtifact(ctx, params)
	if err != nil {
		return repo.Artifact{}, err
	}
	return s.repo.GetArtifact(ctx, id)
}

func (s *service) DeleteArtifact(ctx context.Context, id int32) error {
	return s.repo.DeleteArtifact(ctx, id)
}
