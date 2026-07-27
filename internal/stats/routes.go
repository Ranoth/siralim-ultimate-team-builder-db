package stats

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Get("/", h.GetStats)
	r.Get("/id", h.GetStat)
	r.Get("/type", h.GetStatsByType)
}
