package stats

import (
	"context"
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

func (h *handler) GetStats(w http.ResponseWriter, r *http.Request) {
	httpx.HandleList[repo.Stat](h.service.GetStats)(w, r)
}

func (h *handler) GetStat(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByID[repo.Stat](h.service.GetStat, "stat")(w, r)
}

func (h *handler) GetStatsByType(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.Stat](h.service.GetStatsByType, "stats")(w, r)
}

func (h *handler) CreateStat(w http.ResponseWriter, r *http.Request) {
	type createStatRequest struct {
		Type repo.StatType `json:"type"`
	}

	httpx.HandleCreate[createStatRequest, repo.Stat](func(ctx context.Context, request createStatRequest) (repo.Stat, error) {
		return h.service.CreateStat(ctx, request.Type)
	}, "stat")(w, r)
}

func (h *handler) DeleteStat(w http.ResponseWriter, r *http.Request) {
	httpx.HandleDelete[repo.Stat](h.service.DeleteStat, "stat")(w, r)
}
