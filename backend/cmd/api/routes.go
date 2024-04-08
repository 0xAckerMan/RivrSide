package main

import "github.com/go-chi/chi/v5"

func (app *Application) routes() *chi.Mux{
    r := chi.NewRouter()

    routes := chi.NewRouter()
    r.Mount("/api/v1", routes)
    routes.Get("/healthcheck", app.healthcheck)
    routes.Get("/packages", app.handle_getAllPackages)
    return r
}
