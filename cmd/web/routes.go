package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ol-ilyassov/booking-app/pkg/config"
	"github.com/ol-ilyassov/booking-app/pkg/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.AboutUs)

	return mux
}

// As an example: pat mux.
// func routes(app *config.AppConfig) http.Handler {
// 	mux := pat.New()

// 	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
// 	mux.Get("/about", http.HandlerFunc(handlers.Repo.AboutUs))

// 	return mux
// }
