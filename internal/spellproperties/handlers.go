package spellproperties

import (
	"net/http"

	repo "github.com/Ranoth/SUTBDB/internal/adapters/postgresql/sqlc"
	"github.com/Ranoth/SUTBDB/internal/httpx"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) GetSpellProperties(w http.ResponseWriter, r *http.Request) {
	httpx.HandleList[repo.SpellProperty](h.service.GetSpellProperties)(w, r)
}

func (h *handler) GetSpellProperty(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByID[repo.SpellProperty](h.service.GetSpellProperty, "spell property")(w, r)
}

func (h *handler) GetSpellPropertiesByName(w http.ResponseWriter, r *http.Request) {
	httpx.HandleGetByName[repo.SpellProperty](h.service.GetSpellPropertiesByName, "spell properties")(w, r)
}

func (h *handler) CreateSpellProperty(w http.ResponseWriter, r *http.Request) {
	httpx.HandleCreate[repo.CreateSpellPropertyParams, repo.SpellProperty](h.service.CreateSpellProperty, "spell property")(w, r)
}

func (h *handler) DeleteSpellProperty(w http.ResponseWriter, r *http.Request) {
	httpx.HandleDelete[repo.SpellProperty](h.service.DeleteSpellProperty, "spell property")(w, r)
}
