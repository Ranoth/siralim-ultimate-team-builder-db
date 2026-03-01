package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	repo "github.com/Ranoth/siralim-ultimate-team-builder-db/internal/adapters/postgresql/sqlc"
	"github.com/Ranoth/siralim-ultimate-team-builder-db/internal/artifacts"
	"github.com/Ranoth/siralim-ultimate-team-builder-db/internal/classes"
	"github.com/Ranoth/siralim-ultimate-team-builder-db/internal/creatures"
	"github.com/Ranoth/siralim-ultimate-team-builder-db/internal/dbseeder"
	"github.com/Ranoth/siralim-ultimate-team-builder-db/internal/materials"
	"github.com/Ranoth/siralim-ultimate-team-builder-db/internal/perks"
	"github.com/Ranoth/siralim-ultimate-team-builder-db/internal/races"
	"github.com/Ranoth/siralim-ultimate-team-builder-db/internal/specializations"
	"github.com/Ranoth/siralim-ultimate-team-builder-db/internal/spellproperties"
	"github.com/Ranoth/siralim-ultimate-team-builder-db/internal/spells"
	"github.com/Ranoth/siralim-ultimate-team-builder-db/internal/stats"
	"github.com/Ranoth/siralim-ultimate-team-builder-db/internal/traits"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root."))
	})

	querier := repo.New(app.db)

	routes := []route{
		{
			path: "/traits",
			register: func(r chi.Router, q repo.Querier) {
				traits.RegisterRoutes(r, traits.NewService(q))
			},
		},
		{
			path: "/creatures",
			register: func(r chi.Router, q repo.Querier) {
				creatures.RegisterRoutes(r, creatures.NewService(q))
			},
		},
		{
			path: "/classes",
			register: func(r chi.Router, q repo.Querier) {
				classes.RegisterRoutes(r, classes.NewService(q))
			},
		},
		{
			path: "/races",
			register: func(r chi.Router, q repo.Querier) {
				races.RegisterRoutes(r, races.NewService(q))
			},
		},
		{
			path: "/spells",
			register: func(r chi.Router, q repo.Querier) {
				spells.RegisterRoutes(r, spells.NewService(q))
			},
		},
		{
			path: "/spell-properties",
			register: func(r chi.Router, q repo.Querier) {
				spellproperties.RegisterRoutes(r, spellproperties.NewService(q))
			},
		},
		{
			path: "/perks",
			register: func(r chi.Router, q repo.Querier) {
				perks.RegisterRoutes(r, perks.NewService(q))
			},
		},
		{
			path: "/specializations",
			register: func(r chi.Router, q repo.Querier) {
				specializations.RegisterRoutes(r, specializations.NewService(q))
			},
		},
		{
			path: "/artifacts",
			register: func(r chi.Router, q repo.Querier) {
				artifacts.RegisterRoutes(r, artifacts.NewService(q))
			},
		},
		{
			path: "/stats",
			register: func(r chi.Router, q repo.Querier) {
				stats.RegisterRoutes(r, stats.NewService(q))
			},
		},
		{
			path: "/materials",
			register: func(r chi.Router, q repo.Querier) {
				materials.RegisterRoutes(r, materials.NewService(q))
			},
		},
	}

	for _, route := range routes {
		rt := route
		r.Route(rt.path, func(r chi.Router) {
			rt.register(r, querier)
		})
	}

	return r
}

func (app *application) run(h http.Handler) error {
	app.srv = &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Starting server on %s", app.config.addr)

	dbseeder.Run()

	return app.srv.ListenAndServe()
}

func (app *application) shutdown(ctx context.Context) error {
	slog.Info("Shutting down HTTP server...")
	if err := app.srv.Shutdown(ctx); err != nil {
		return err
	}

	slog.Info("Closing database connection...")
	if err := app.db.Close(ctx); err != nil {
		return err
	}

	return nil
}

type application struct {
	config config
	db     *pgx.Conn
	srv    *http.Server
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

type route struct {
	path     string
	register func(r chi.Router, q repo.Querier)
}
