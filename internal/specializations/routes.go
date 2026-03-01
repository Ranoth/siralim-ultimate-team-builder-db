package specializations

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Get("/", h.GetSpecializations)
	r.Get("/id", h.GetSpecialization)
	r.Get("/name", h.GetSpecializationsByName)
	r.Post("/create", h.CreateSpecialization)
	r.Delete("/delete", h.DeleteSpecialization)
}
