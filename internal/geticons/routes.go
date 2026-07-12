package geticons

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, service Service) error {
	h := NewHandler(service)

	r.Get("/creatures/{id}", h.GetCreatureIconById)
	r.Get("/materials/{id}", h.GetMaterialIconById)
	r.Get("/relics/{id}", h.GetRelicIconById)
	r.Get("/races/{id}", h.GetRaceIconById)
	r.Get("/classes/{id}", h.GetClassIconById)

	return nil
}
