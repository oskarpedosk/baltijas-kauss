package main

import (
	"github.com/oskarpedosk/baltijas-kauss/internal/config"
	"github.com/oskarpedosk/baltijas-kauss/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.SignIn)
	mux.Get("/nba", handlers.Repo.NBAHome)
	mux.Get("/nba/players", handlers.Repo.NBAPlayers)

	mux.Get("/nba/teams", handlers.Repo.NBATeams)
	mux.Post("/nba/teams", handlers.Repo.PostNBATeams)
	// mux.Post("/nba/teams-json", handlers.Repo.NBATeamsAvailabilityJSON)
	mux.Get("/nba/team-info-summary", handlers.Repo.NBATeamInfoSummary)

	mux.Get("/nba/results", handlers.Repo.NBAResults)
	mux.Post("/nba/results", handlers.Repo.PostNBAResults)

	fileServer := http.FileServer(http.Dir("../../static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
