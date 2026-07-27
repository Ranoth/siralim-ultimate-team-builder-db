package spellproperties

import (
	"context"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	GetSpellProperties(ctx context.Context) ([]repo.SpellProperty, error)
	GetSpellProperty(ctx context.Context, id int32) (repo.SpellProperty, error)
	GetSpellPropertiesByName(ctx context.Context, name string) ([]repo.SpellProperty, error)
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) *service {
	return &service{repo: repo}
}

func (s *service) GetSpellProperties(ctx context.Context) ([]repo.SpellProperty, error) {
	return s.repo.GetSpellProperties(ctx)
}

func (s *service) GetSpellProperty(ctx context.Context, id int32) (repo.SpellProperty, error) {
	return s.repo.GetSpellProperty(ctx, id)
}

func (s *service) GetSpellPropertiesByName(ctx context.Context, name string) ([]repo.SpellProperty, error) {
	return s.repo.GetSpellPropertiesByName(ctx, pgtype.Text{String: name, Valid: true})
}
