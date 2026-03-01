package classes

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, service Service) error {
	handler := NewHandler(service)

	r.Get("/", handler.GetClasses)
	r.Get("/id", handler.GetClass)
	r.Get("/name", handler.GetClassesByName)
	r.Post("/create", handler.CreateClass)
	r.Delete("/delete", handler.DeleteClass)

	return nil
}
