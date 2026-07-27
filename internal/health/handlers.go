package health

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

func (h *handler) GetHealth(w http.ResponseWriter, r *http.Request) {
	status, healthy := h.service.Check(r.Context())

	httpCode := http.StatusOK
	if !healthy {
		httpCode = http.StatusServiceUnavailable
	}

	httpx.WriteJSON(w, httpCode, status)
}
