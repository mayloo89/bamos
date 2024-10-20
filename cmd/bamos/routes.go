package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/mayloo89/bamos/internal/config"
	"github.com/mayloo89/bamos/internal/handler"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", handler.Repo.Home)
	mux.Get("/colectivos/vehiclePositionsSimple", handler.Repo.VehiclePositionsSimple)
	mux.Get("/colectivos/feed-gtfs-frequency", handler.Repo.FeedGtfsFrequency)

	mux.Get("/colectivos/search", handler.Repo.SearchLine)
	mux.Post("/colectivos/search", handler.Repo.PostSearchLine)

	return mux
}
