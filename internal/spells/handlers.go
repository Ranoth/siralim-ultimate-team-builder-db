package spells

import (
	"net/http"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
	"github.com/Ranoth/siralim-ultimate-team-builder-db/internal/httpx"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) GetSpells(w http.ResponseWriter, r *http.Request) {
	httpx.HandleList[repo.Spell](h.service.GetSpells)(w, r)
}

func (h *handler) GetSpell(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByID[repo.Spell](h.service.GetSpell, "spell")(w, r)
}

func (h *handler) GetSpellsByName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.Spell](h.service.GetSpellsByName, "spells")(w, r)
}

func (h *handler) CreateSpell(w http.ResponseWriter, r *http.Request) {
	httpx.HandleCreate[repo.CreateSpellParams, repo.Spell](h.service.CreateSpell, "spell")(w, r)
}

func (h *handler) DeleteSpell(w http.ResponseWriter, r *http.Request) {
	httpx.HandleDelete[repo.Spell](h.service.DeleteSpell, "spell")(w, r)
}
