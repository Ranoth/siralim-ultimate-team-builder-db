package traits

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

func (h *handler) GetTraits(w http.ResponseWriter, r *http.Request) {
	httpx.HandleList[repo.Trait](h.service.GetTraits)(w, r)
}

func (h *handler) GetTrait(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByID[repo.Trait](h.service.GetTrait, "trait")(w, r)
}

func (h *handler) GetTraitsByName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.Trait](h.service.GetTraitsByName, "traits")(w, r)
}
