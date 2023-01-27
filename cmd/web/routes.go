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

	mux.Get("/", handlers.Repo.Login)
	mux.Post("/", handlers.Repo.PostLogin)
	mux.Get("/logout", handlers.Repo.Logout)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/", func(mux chi.Router) {
		mux.Use(Auth)
		
		mux.Get("/home", handlers.Repo.NBAHome)

		mux.Get("/{src}/{id}", handlers.Repo.Player)
		mux.Post("/{src}/{id}", handlers.Repo.PostPlayer)
		
		mux.Get("/players", handlers.Repo.Players)
		mux.Post("/players", handlers.Repo.PostPlayers)

		mux.Get("/players/page={page}", handlers.Repo.Players)
		mux.Post("/players/page={page}", handlers.Repo.PostPlayers)

		mux.Get("/teams", handlers.Repo.NBATeams)
		mux.Post("/teams", handlers.Repo.PostNBATeams)
		// mux.Post("/nba/teams-json", handlers.Repo.NBATeamsAvailabilityJSON)
		mux.Get("/team-info-summary", handlers.Repo.NBATeamInfoSummary)

		mux.Get("/standings", handlers.Repo.NBAResults)
		mux.Post("/standings", handlers.Repo.PostNBAResults)

		mux.Get("/draft", handlers.Repo.NBADraft)
		mux.Get("/ws", handlers.Repo.WsEndPoint)
	})

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(AuthAdmin)
		// add routes here for admin
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
		
		mux.Get("/teams", handlers.Repo.AdminNBATeams)
		mux.Get("/players", handlers.Repo.AdminNBAPlayers)
		mux.Get("/{src}/{id}", handlers.Repo.AdminShowNBAPlayer)
		mux.Get("/standings", handlers.Repo.AdminNBAResults)
		
	})

	return mux
}