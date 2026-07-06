package creatures

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

func (h *handler) GetCreatures(w http.ResponseWriter, r *http.Request) {
	httpx.HandleList[repo.GetCreaturesRow](h.service.GetCreatures)(w, r)
}

func (h *handler) GetCreature(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByID[repo.GetCreatureRow](h.service.GetCreature, "creature")(w, r)
}

func (h *handler) GetCreaturesByName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.Creature](h.service.GetCreaturesByName, "creatures")(w, r)
}

func (h *handler) GetCreaturesByTraitName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.Creature](h.service.GetCreaturesByTraitName, "creatures")(w, r)
}

func (h *handler) GetCreaturesByClassName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.Creature](h.service.GetCreaturesByClassName, "creatures")(w, r)
}

func (h *handler) GetCreaturesByRaceName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.Creature](h.service.GetCreaturesByRaceName, "creatures")(w, r)
}
