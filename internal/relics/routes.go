package relics

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Get("/", h.GetRelics)
	r.Get("/id", h.GetRelic)
	r.Get("/name", h.GetRelicsByName)
	r.Post("/create", h.CreateRelic)
	r.Delete("/delete", h.DeleteRelic)
}
