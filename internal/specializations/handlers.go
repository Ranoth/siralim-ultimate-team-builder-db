package specializations

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

func (h *handler) GetSpecializations(w http.ResponseWriter, r *http.Request) {
	httpx.HandleList[repo.Specialization](h.service.GetSpecializations)(w, r)
}

func (h *handler) GetSpecialization(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByID[repo.Specialization](h.service.GetSpecialization, "specialization")(w, r)
}

func (h *handler) GetSpecializationsByName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.Specialization](h.service.GetSpecializationsByName, "specializations")(w, r)
}

func (h *handler) CreateSpecialization(w http.ResponseWriter, r *http.Request) {
	httpx.HandleCreate[repo.CreateSpecializationParams, repo.Specialization](h.service.CreateSpecialization, "specialization")(w, r)
}

func (h *handler) DeleteSpecialization(w http.ResponseWriter, r *http.Request) {
	httpx.HandleDelete[repo.Specialization](h.service.DeleteSpecialization, "specialization")(w, r)
}
