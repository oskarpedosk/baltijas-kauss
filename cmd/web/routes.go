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

		mux.Get("/players/{id}", handlers.Repo.Player)
		mux.Post("/players/{id}", handlers.Repo.PostPlayer)

		mux.Post("/update", handlers.Repo.PostUpdatePlayer)
		
		mux.Get("/players", handlers.Repo.Players)
		mux.Post("/players", handlers.Repo.PostPlayers)

		mux.Get("/team/{id}", handlers.Repo.Team)
		mux.Post("/team/{id}", handlers.Repo.PostTeam)

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