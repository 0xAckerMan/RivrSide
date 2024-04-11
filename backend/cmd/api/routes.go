package main

import "github.com/go-chi/chi/v5"

func (app *Application) routes() *chi.Mux {
	r := chi.NewRouter()
	r.NotFound(app.notFoundResponse)
    r.MethodNotAllowed(app.methodNotAllowedResponse)

	routes := chi.NewRouter()
	r.Mount("/api/v1", routes)
	routes.Get("/healthcheck", app.healthcheck)
	routes.Group(func(r chi.Router) {
		r.Get("/packageplans", app.HandleGetAllPackages)
		r.Post("/packageplan", app.HandleCreatePackagePlan)
        r.Get("/packageplan/{id}", app.HandleGetSinglePackage)
        r.Patch("/packageplan/{id}", app.HandleUpdatePackagePlan)
		r.Delete("/packageplan/{id}", app.HandleDeletePackagePlan)
	})

    routes.Group(func(r chi.Router) {
        r.Get("/rooms", app.HandleGetAllRooms)
        r.Get("/room/{id}", app.HandleGetSingleRoom)
        r.Post("/room", app.HandleCreateRoom)
        r.Patch("/room/{id}", app.HandlePartialUpdateRoom)
        r.Delete("/room/{id}", app.HandleDeleteRoom)
    })

	admin := chi.NewRouter()
	routes.Mount("/admin", admin)
	admin.Group(func(r chi.Router) {
		r.Get("/healthcheck", app.healthcheck)
        r.Route("/role", func(r chi.Router) {
            r.Post("/", app.HandleCreateRole)
            r.Get("/", app.HandleGetAllRole)
            r.Patch("/{id}", app.HandleUpdateRole)
            r.Delete("/{id}", app.HandleDeleteRole)
        })
	})

    manager := chi.NewRouter()
    routes.Mount("/manager", manager)
    manager.Group(func(r chi.Router) {
        r.Get("/healthcheck", app.healthcheck)
        r.Route("/tenant", func(r chi.Router) {
            r.Post("/",app.HandleCreateNewTenant)
            r.Get("/", app.HandleGetAllTenants)
        })
    })
	return r
}
