package artifacts

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Get("/", h.GetArtifacts)
	r.Get("/id", h.GetArtifact)
	r.Get("/name", h.GetArtifactsByName)
	r.Post("/create", h.CreateArtifact)
	r.Delete("/delete", h.DeleteArtifact)
}
