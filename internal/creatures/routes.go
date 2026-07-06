package creatures

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, service Service) error {
	h := NewHandler(service)

	r.Get("/", h.GetCreatures)
	r.Get("/id", h.GetCreature)
	r.Get("/name", h.GetCreaturesByName)
	r.Get("/trait", h.GetCreaturesByTraitName)
	r.Get("/class", h.GetCreaturesByClassName)
	r.Get("/race", h.GetCreaturesByRaceName)

	return nil
}
