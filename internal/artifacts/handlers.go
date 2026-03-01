package artifacts

import (
	"net/http"

	repo "github.com/Ranoth/SUTBDB/internal/adapters/postgresql/sqlc"
	"github.com/Ranoth/SUTBDB/internal/httpx"
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

func (h *handler) CreateArtifact(w http.ResponseWriter, r *http.Request) {
	httpx.HandleCreate[repo.CreateArtifactParams, repo.Artifact](h.service.CreateArtifact, "artifact")(w, r)
}

func (h *handler) DeleteArtifact(w http.ResponseWriter, r *http.Request) {
	httpx.HandleDelete[repo.Artifact](h.service.DeleteArtifact, "artifact")(w, r)
}
