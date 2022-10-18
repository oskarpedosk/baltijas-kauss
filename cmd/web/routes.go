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
	mux.Post("/", handlers.Repo.PostSignIn)
	mux.Get("/logout", handlers.Repo.Logout)

	mux.Get("/nba", handlers.Repo.NBAHome)

	mux.Get("/nba/players", handlers.Repo.NBAPlayers)
	mux.Post("/nba/players", handlers.Repo.PostNBAPlayers)

	mux.Get("/nba/teams", handlers.Repo.NBATeams)
	mux.Post("/nba/teams", handlers.Repo.PostNBATeams)
	// mux.Post("/nba/teams-json", handlers.Repo.NBATeamsAvailabilityJSON)
	mux.Get("/nba/team-info-summary", handlers.Repo.NBATeamInfoSummary)

	mux.Get("/nba/results", handlers.Repo.NBAResults)
	mux.Post("/nba/results", handlers.Repo.PostNBAResults)

	mux.Get("/nba/draft", handlers.Repo.NBADraft)

	fileServer := http.FileServer(http.Dir("../../static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
		// add routes here for logged in user
		
		mux.Get("/nba_teams", handlers.Repo.AdminNBATeams)
		mux.Get("/nba_players", handlers.Repo.AdminNBAPlayers)
		mux.Get("/nba_results", handlers.Repo.AdminNBAResults)
		
	})

	return mux
}
