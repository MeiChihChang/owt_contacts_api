package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *application) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Post("/authenticate", app.Authenticate)

	mux.Route("/contact", func(mux chi.Router) {
		mux.Use(app.authRequired)

		mux.Get("/all", app.AllContacts)
		mux.Get("/get/{id}", app.GetContact)
		mux.Put("/new", app.InsertContact)
		mux.Patch("/update/{id}", app.UpdateContact)
		mux.Delete("/delete/{id}", app.DeleteContact)
	})

	mux.Route("/skill", func(mux chi.Router) {
		mux.Use(app.authRequired)

		mux.Get("/all", app.AllSkills)
		mux.Get("/get/{id}", app.GetSkill)
		mux.Put("/new", app.InsertSkill)
		mux.Patch("/update/{id}", app.UpdateSkill)
		mux.Delete("/delete/{id}", app.DeleteSkill)
	})

	return mux
}
