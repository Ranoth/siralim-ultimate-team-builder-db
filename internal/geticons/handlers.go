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
	httpx.HandleGetIconRawBytesByIdInPath(h.service.GetCreatureIconById, "creature")(w, r)
}

func (h *handler) GetMaterialIconById(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetIconRawBytesByIdInPath(h.service.GetMaterialIconById, "material")(w, r)
}

func (h *handler) GetRelicIconById(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetIconRawBytesByIdInPath(h.service.GetRelicIconById, "relic")(w, r)
}

func (h *handler) GetRaceIconById(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetIconRawBytesByIdInPath(h.service.GetRaceIconById, "race")(w, r)
}

func (h *handler) GetClassIconById(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetIconRawBytesByIdInPath(h.service.GetClassIconById, "class")(w, r)
}
