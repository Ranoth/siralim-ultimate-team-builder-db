package relics

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

func (h *handler) GetRelics(w http.ResponseWriter, r *http.Request) {
	httpx.HandleList[repo.GetRelicsRow](h.service.GetRelics)(w, r)
}

func (h *handler) GetRelic(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByID[repo.GetRelicRow](h.service.GetRelic, "relic")(w, r)
}

func (h *handler) GetRelicsByName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.GetRelicsByNameRow](h.service.GetRelicsByName, "relics")(w, r)
}
