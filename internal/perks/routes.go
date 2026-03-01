package perks

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Get("/", h.GetPerks)
	r.Get("/id", h.GetPerk)
	r.Get("/name", h.GetPerksByName)
	r.Post("/create", h.CreatePerk)
	r.Delete("/delete", h.DeletePerk)
}
