package main

import (
	"github.com/oskarpedosk/baltijas-kauss/pkg/config"
	"github.com/oskarpedosk/baltijas-kauss/pkg/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/player_ratings", handlers.Repo.Players)

	fileServer := http.FileServer(http.Dir("../../static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
