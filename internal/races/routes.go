package races

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Get("/", h.GetRace)
	r.Get("/name", h.GetRacesByName)
	r.Get("/trait", h.GetRacesByTraitName)
	r.Get("/class", h.GetRacesByClassName)
	r.Get("/creature", h.GetRacesByCreatureName)
	r.Post("/create", h.CreateRace)
	r.Delete("/delete", h.DeleteRace)
}
