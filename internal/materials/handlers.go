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
	httpx.HandleList[repo.GetMaterialsRow](h.service.GetMaterials)(w, r)
}

func (h *handler) GetMaterial(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByID[repo.GetMaterialRow](h.service.GetMaterial, "material")(w, r)
}

func (h *handler) GetMaterialsByName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.GetMaterialsRow](h.service.GetMaterialsByName, "materials")(w, r)
}

func (h *handler) CreateMaterial(w http.ResponseWriter, r *http.Request) {
	httpx.HandleCreate[repo.CreateMaterialParams, repo.Material](h.service.CreateMaterial, "material")(w, r)
}

func (h *handler) DeleteMaterial(w http.ResponseWriter, r *http.Request) {
	httpx.HandleDelete[repo.GetMaterialRow](h.service.DeleteMaterial, "material")(w, r)
}

func (h *handler) CreateMaterialStat(w http.ResponseWriter, r *http.Request) {
	httpx.HandleCreate[repo.CreateMaterialStatParams, int32](h.service.CreateMaterialStat, "material stat")(w, r)
}

func (h *handler) DeleteMaterialStat(w http.ResponseWriter, r *http.Request) {
	httpx.HandleDelete[int32](h.service.DeleteMaterialStat, "material stat")(w, r)
}
