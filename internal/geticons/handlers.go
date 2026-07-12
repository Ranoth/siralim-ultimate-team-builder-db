package geticons

import (
	"net/http"

	"github.com/Ranoth/siralim-ultimate-team-builder-db/internal/httpx"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) GetCreatureIconById(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByIDInPath[[]byte](h.service.GetCreatureIconById, "creature")(w, r)
}

func (h *handler) GetMaterialIconById(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByIDInPath[[]byte](h.service.GetMaterialIconById, "material")(w, r)
}

func (h *handler) GetRelicIconById(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByIDInPath[[]byte](h.service.GetRelicIconById, "relic")(w, r)
}
