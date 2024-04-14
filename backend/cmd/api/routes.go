package main

import (
	"expvar"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Application) routes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.NotFound(app.notFoundResponse)
	r.MethodNotAllowed(app.methodNotAllowedResponse)

	r.Get("/", app.ApiHome)
	r.Get("/debug/vars", func(w http.ResponseWriter, r *http.Request) {
		expvar.Handler().ServeHTTP(w, r)
	})
	routes := chi.NewRouter()
	r.Mount("/api/v1", routes)
	routes.Get("/healthcheck", app.healthcheck)
	routes.Post("/login", app.UserLogin)
	routes.Post("/signout", app.SignOut)

	routes.Patch("/user/{id}", app.HandleUpdateUser)

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
		// manager endpoints
		r.Get("/rooms/vaccant", app.HandleGetAllVacantRooms)
		r.Get("/rooms/male_vaccant", app.HandleGetMaleVacantRooms)
		r.Get("/rooms/female_vaccant", app.HandleGetFemaleVacantRooms)
		r.Get("/room/{id}/tenants", app.HandleGetRoomTenants)
	})

	admin := chi.NewRouter()
	routes.Mount("/admin", admin)
	admin.Group(func(r chi.Router) {
		r.Get("/healthcheck", app.healthcheck)
		r.Delete("/user/{id}", app.HandleDeleteUser)
		r.Route("/role", func(r chi.Router) {
			r.Post("/", app.HandleCreateRole)
			r.Get("/", app.HandleGetAllRole)
			r.Patch("/{id}", app.HandleUpdateRole)
			r.Delete("/{id}", app.HandleDeleteRole)
		})

		r.Route("/manager", func(r chi.Router) {
			r.Post("/", app.HandleCreateManager)
			r.Get("/", app.HandleGetAllManagers)
		})
	})

	manager := chi.NewRouter()
	routes.Mount("/manager", manager)
	manager.Group(func(r chi.Router) {
		r.Get("/healthcheck", app.healthcheck)
		r.Route("/tenant", func(r chi.Router) {
			r.Post("/", app.HandleCreateNewTenant)
			r.Get("/", app.HandleGetAllTenants)
			r.Get("/active", app.HandleGetAllActiveTenants)
			r.Get("/{id}", app.HandleGetTenantInfo)
		})
	})

	tenant := chi.NewRouter()
	routes.Mount("/tenant", tenant)
	tenant.Group(func(r chi.Router) {
		r.Use(app.UserMiddleware)
		r.Get("/me", app.HandleGetProfileInfo)
		r.Patch("/me/update", app.HandleUpdateTenantInfo)
	})
	return r
}
