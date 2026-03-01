package creatures

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, service Service) error {
	handler := NewHandler(service)

	r.Get("/", handler.GetCreatures)
	r.Get("/id", handler.GetCreature)
	r.Get("/name", handler.GetCreaturesByName)
	r.Get("/trait", handler.GetCreaturesByTraitName)
	r.Get("/class", handler.GetCreaturesByClassName)
	r.Get("/race", handler.GetCreaturesByRaceName)
	r.Post("/create", handler.CreateCreature)
	r.Delete("/delete", handler.DeleteCreature)

	return nil
}
