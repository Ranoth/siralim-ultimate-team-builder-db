package spells

import (
	"context"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	GetSpells(ctx context.Context) ([]repo.Spell, error)
	GetSpell(ctx context.Context, id int32) (repo.Spell, error)
	GetSpellsByName(ctx context.Context, name string) ([]repo.Spell, error)
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) *service {
	return &service{repo: repo}
}

func (s *service) GetSpells(ctx context.Context) ([]repo.Spell, error) {
	return s.repo.GetSpells(ctx)
}

func (s *service) GetSpell(ctx context.Context, id int32) (repo.Spell, error) {
	return s.repo.GetSpell(ctx, id)
}

func (s *service) GetSpellsByName(ctx context.Context, name string) ([]repo.Spell, error) {
	return s.repo.GetSpellsByName(ctx, pgtype.Text{String: name, Valid: true})
}
