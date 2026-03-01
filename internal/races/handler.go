package races

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

func (h *handler) GetRaces(w http.ResponseWriter, r *http.Request) {
	httpx.HandleList[repo.Race](h.service.GetRaces)(w, r)
}

func (h *handler) GetRace(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByID[repo.Race](h.service.GetRace, "race")(w, r)
}

func (h *handler) GetRacesByName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.Race](h.service.GetRacesByName, "races")(w, r)
}

func (h *handler) GetRacesByTraitName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.Race](h.service.GetRacesByTraitName, "races")(w, r)
}

func (h *handler) GetRacesByClassName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.Race](h.service.GetRacesByClassName, "races")(w, r)
}

func (h *handler) GetRacesByCreatureName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.Race](h.service.GetRacesByCreatureName, "races")(w, r)
}

func (h *handler) CreateRace(w http.ResponseWriter, r *http.Request) {
	httpx.HandleCreate[repo.CreateRaceParams, repo.Race](h.service.CreateRace, "race")(w, r)
}

func (h *handler) DeleteRace(w http.ResponseWriter, r *http.Request) {
	httpx.HandleDelete[repo.Race](h.service.DeleteRace, "race")(w, r)
}
