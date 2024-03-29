package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(secureHeaders)
	router.Use(app.sessionManager.LoadAndSave)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("/*", http.StripPrefix("/static/", fileServer))

	router.Get("/", app.home)
	router.Get("/about", app.about)

	// group snippet routes
	router.Route("/snippet", func(r chi.Router) {
		r.Use(app.noSurf)
		r.Use(app.authenticate)

		r.Get("/view/{id}", app.viewSnippet)

		r.Group(func(r chi.Router) {
			r.Use(app.requireAuthentication)
			r.Get("/create", app.createSnippet)
			r.Post("/create", app.createSnippetPost)
		})
	})

	// group user routes
	router.Route("/user", func(r chi.Router) {
		r.Use(app.noSurf)
		r.Use(app.authenticate)

		r.Get("/signup", app.userSignUp)
		r.Post("/signup", app.userSignUpPost)
		r.Get("/login", app.userLogin)
		r.Post("/login", app.userLoginPost)

		r.Group(func(r chi.Router) {
			r.Use(app.requireAuthentication)
			r.Post("/logout", app.userLogout)
		})
	})

	return router
}
