package classes

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

func (h *handler) GetClasses(w http.ResponseWriter, r *http.Request) {
	httpx.HandleList[repo.Class](h.service.GetClasses)(w, r)
}

func (h *handler) GetClass(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByID[repo.Class](h.service.GetClass, "class")(w, r)
}

func (h *handler) GetClassesByName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.Class](h.service.GetClassesByName, "classes")(w, r)
}

func (h *handler) CreateClass(w http.ResponseWriter, r *http.Request) {
	httpx.HandleCreate[repo.CreateClassParams, repo.Class](h.service.CreateClass, "class")(w, r)
}

func (h *handler) DeleteClass(w http.ResponseWriter, r *http.Request) {
	httpx.HandleDelete[repo.Class](h.service.DeleteClass, "class")(w, r)
}
