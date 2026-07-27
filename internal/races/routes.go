package races

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Get("/", h.GetRaces)
	r.Get("/id", h.GetRace)
	r.Get("/name", h.GetRacesByName)
	r.Get("/trait", h.GetRacesByTraitName)
	r.Get("/class", h.GetRacesByClassName)
	r.Get("/creature", h.GetRacesByCreatureName)
}
