package materials

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

func (h *handler) GetMaterials(w http.ResponseWriter, r *http.Request) {
	httpx.HandleList[repo.Material](h.service.GetMaterials)(w, r)
}

func (h *handler) GetMaterial(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByID[repo.Material](h.service.GetMaterial, "material")(w, r)
}

func (h *handler) GetMaterialsByName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.Material](h.service.GetMaterialsByName, "materials")(w, r)
}

func (h *handler) CreateMaterial(w http.ResponseWriter, r *http.Request) {
	httpx.HandleCreate[repo.CreateMaterialParams, repo.Material](h.service.CreateMaterial, "material")(w, r)
}

func (h *handler) DeleteMaterial(w http.ResponseWriter, r *http.Request) {
	httpx.HandleDelete[repo.Material](h.service.DeleteMaterial, "material")(w, r)
}
