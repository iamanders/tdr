package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
)

// Setup application routes
func (app *application) SetupRoutes() {
	// Middlewares
	// app.Router.Use(middleware.RequestID)
	app.Router.Use(middleware.RealIP)
	app.Router.Use(middleware.Logger)
	app.Router.Use(middleware.Recoverer)
	app.Router.Use(middleware.Timeout(60 * time.Second))

	// CSRF middleware
	csrfProtection := csrf.Protect(
		[]byte(app.Config.CsrfKey),
		csrf.Secure(false),   // Needed?
		csrf.HttpOnly(false), // Needed?
		csrf.Path("/"),       // Needed?
	)
	app.Router.Use(csrfProtection)

	// Static route
	if app.Config.Env == "production" {
		// Production - Embed assets
		var staticFS = http.FS(staticFiles)
		fileServer := http.FileServer(staticFS)
		app.Router.Handle("/static/*", fileServer)
	} else {
		// Dev - Read files from disk
		fs := http.FileServer(http.Dir("static/"))
		app.Router.Handle("/static/*", http.StripPrefix("/static/", fs))
	}

	// Routes
	app.Router.Get("/", app.GetHome)

	// Day
	app.Router.Get("/day", app.GetDay)
	app.Router.Get("/day/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}", app.GetDay)

	// Week
	app.Router.Get("/week", app.GetWeek)
	app.Router.Get("/week/{year:[0-9]{4}}-{week:[0-9]{1,2}}", app.GetWeek)

	// Time
	app.Router.Get("/time/{timeId}/edit", app.GetTimeEdit)
	app.Router.Post("/time/{timeId}/update", app.PostTimeUpdate)
	app.Router.Post("/time/{timeId}/delete", app.PostTimeDelete)
	app.Router.Post("/time/store", app.PostTimeStore)
	app.Router.Post("/time/{timeId}/stop", app.PostTimeStop)
}
