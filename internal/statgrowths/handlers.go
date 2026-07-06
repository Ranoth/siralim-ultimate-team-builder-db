package statgrowths

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

func (h *handler) GetStatGrowthsByCreatureId(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByID[[]repo.GetStatGrowthsByCreatureIdRow](h.service.GetStatGrowthsByCreatureId, "stat growths")(w, r)
}
