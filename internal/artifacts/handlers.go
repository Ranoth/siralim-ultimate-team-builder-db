package artifacts

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

func (h *handler) GetArtifacts(w http.ResponseWriter, r *http.Request) {
	httpx.HandleList[repo.Artifact](h.service.GetArtifacts)(w, r)
}

func (h *handler) GetArtifact(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByID[repo.Artifact](h.service.GetArtifact, "artifact")(w, r)
}

func (h *handler) GetArtifactsByName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.Artifact](h.service.GetArtifactsByName, "artifacts")(w, r)
}
