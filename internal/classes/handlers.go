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
	httpx.HandleList[repo.ClassesView](h.service.GetClasses)(w, r)
}

func (h *handler) GetClass(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByID[repo.ClassesView](h.service.GetClass, "class")(w, r)
}

func (h *handler) GetClassesByName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.ClassesView](h.service.GetClassesByName, "classes")(w, r)
}
