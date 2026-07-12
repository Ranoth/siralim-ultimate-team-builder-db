package materials

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

func (h *handler) GetMaterials(w http.ResponseWriter, r *http.Request) {
	httpx.HandleList[repo.MaterialsView](h.service.GetMaterials)(w, r)
}

func (h *handler) GetMaterial(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByID[repo.MaterialsView](h.service.GetMaterial, "material")(w, r)
}

func (h *handler) GetMaterialsByName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.MaterialsView](h.service.GetMaterialsByName, "materials")(w, r)
}
