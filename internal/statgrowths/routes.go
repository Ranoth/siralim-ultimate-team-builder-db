package statgrowths

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, service Service) error {
	h := NewHandler(service)

	r.Get("/creature", h.GetStatGrowthsByCreatureId)

	return nil
}
