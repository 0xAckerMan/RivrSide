package main

import "github.com/go-chi/chi/v5"

func (app *Application) routes() *chi.Mux{
    r := chi.NewRouter()
    r.Route("/api", func(r chi.Router) {
        r.Get("/healthcheck", app.healthcheck)
    })
    return r
}
