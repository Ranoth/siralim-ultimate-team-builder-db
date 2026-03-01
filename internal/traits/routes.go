package traits

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, service Service) error {
	handler := NewHandler(service)

	r.Get("/", handler.GetTraits)
	r.Get("/id", handler.GetTrait)
	r.Get("/name", handler.GetTraitsByName)
	r.Post("/create", handler.CreateTrait)
	r.Delete("/delete", handler.DeleteTrait)

	return nil
}
