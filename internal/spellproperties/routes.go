package spellproperties

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Get("/", h.GetSpellProperties)
	r.Get("/id", h.GetSpellProperty)
	r.Get("/name", h.GetSpellPropertiesByName)
	r.Post("/create", h.CreateSpellProperty)
	r.Delete("/delete", h.DeleteSpellProperty)
}
