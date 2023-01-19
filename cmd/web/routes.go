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

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/", func(mux chi.Router) {
		mux.Use(Auth)
		
		mux.Get("/nba", handlers.Repo.NBAHome)

		mux.Get("/nba/{src}/{id}", handlers.Repo.Player)
		mux.Get("/nba/players", handlers.Repo.NBAPlayers)
		mux.Get("/nba/players/page={page}", handlers.Repo.NBAPlayers)
		mux.Post("/nba/players", handlers.Repo.PostNBAPlayers)

		mux.Get("/nba/teams", handlers.Repo.NBATeams)
		mux.Post("/nba/teams", handlers.Repo.PostNBATeams)
		// mux.Post("/nba/teams-json", handlers.Repo.NBATeamsAvailabilityJSON)
		mux.Get("/nba/team-info-summary", handlers.Repo.NBATeamInfoSummary)

		mux.Get("/nba/standings", handlers.Repo.NBAResults)
		mux.Post("/nba/standings", handlers.Repo.PostNBAResults)

		mux.Get("/nba/draft", handlers.Repo.NBADraft)
		mux.Get("/ws", handlers.Repo.WsEndPoint)
	})

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(AuthAdmin)
		// add routes here for admin
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
		
		mux.Get("/nba_teams", handlers.Repo.AdminNBATeams)
		mux.Get("/nba_players", handlers.Repo.AdminNBAPlayers)
		mux.Get("/{src}/{id}", handlers.Repo.AdminShowNBAPlayer)
		mux.Get("/nba_standings", handlers.Repo.AdminNBAResults)
		
	})

	return mux
}