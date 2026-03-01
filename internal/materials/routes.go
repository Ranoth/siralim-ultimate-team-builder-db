package materials

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Get("/", h.GetMaterials)
	r.Get("/id", h.GetMaterial)
	r.Get("/name", h.GetMaterialsByName)
	r.Post("/create", h.CreateMaterial)
	r.Delete("/delete", h.DeleteMaterial)
}
