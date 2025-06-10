package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/mayloo89/bamos/internal/config"
	"github.com/mayloo89/bamos/internal/handler"
)

func routes(app *config.AppConfig, repo *handler.Repository) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", repo.Home)
	mux.Get("/colectivos/vehiclePositionsSimple", repo.VehiclePositionsSimple)
	// mux.Get("/colectivos/feed-gtfs-frequency", repo.FeedGtfsFrequency)

	mux.Get("/colectivos/search", repo.SearchLine)
	mux.Post("/colectivos/search", repo.PostSearchLine)

	// Allowed Parking
	mux.Get("/transit/allowed-parking", repo.AllowedParking)
	mux.Post("/transit/allowed-parking", repo.PostAllowedParking)

	return mux
}
