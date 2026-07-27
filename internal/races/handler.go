package races

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

func (h *handler) GetRaces(w http.ResponseWriter, r *http.Request) {
	httpx.HandleList[repo.RacesView](h.service.GetRaces)(w, r)
}

func (h *handler) GetRace(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByID[repo.RacesView](h.service.GetRace, "race")(w, r)
}

func (h *handler) GetRacesByName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.RacesView](h.service.GetRacesByName, "races")(w, r)
}

func (h *handler) GetRacesByTraitName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.RacesView](h.service.GetRacesByTraitName, "races")(w, r)
}

func (h *handler) GetRacesByClassName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.RacesView](h.service.GetRacesByClassName, "races")(w, r)
}

func (h *handler) GetRacesByCreatureName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.RacesView](h.service.GetRacesByCreatureName, "races")(w, r)
}
