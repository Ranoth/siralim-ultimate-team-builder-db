package perks

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

func (h *handler) GetPerks(w http.ResponseWriter, r *http.Request) {
	httpx.HandleList[repo.Perk](h.service.GetPerks)(w, r)
}

func (h *handler) GetPerk(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByID[repo.Perk](h.service.GetPerk, "perk")(w, r)
}

func (h *handler) GetPerksByName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.Perk](h.service.GetPerksByName, "perks")(w, r)
}

func (h *handler) CreatePerk(w http.ResponseWriter, r *http.Request) {
	httpx.HandleCreate[repo.CreatePerkParams, repo.Perk](h.service.CreatePerk, "perk")(w, r)
}

func (h *handler) DeletePerk(w http.ResponseWriter, r *http.Request) {
	httpx.HandleDelete[repo.Perk](h.service.DeletePerk, "perk")(w, r)
}
