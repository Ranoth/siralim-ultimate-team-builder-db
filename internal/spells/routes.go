package spells

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Get("/", h.GetSpells)
	r.Get("/id", h.GetSpell)
	r.Get("/name", h.GetSpellsByName)
	r.Post("/create", h.CreateSpell)
	r.Delete("/delete", h.DeleteSpell)
}
