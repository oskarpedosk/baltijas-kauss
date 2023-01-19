package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
	"github.com/oskarpedosk/baltijas-kauss/internal/config"
	"github.com/oskarpedosk/baltijas-kauss/internal/driver"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {
	gob.Register(models.NBATeamInfo{})

	// Change this to true when in production
	app.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	templateCache, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = templateCache
	app.UseCache = true

	repo := NewRepo(&app, &driver.DB{})
	NewHandlers(repo)
	render.NewRenderer(&app)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.SignIn)

	mux.Get("/nba", Repo.NBAHome)

	mux.Get("/nba/players", Repo.NBAPlayers)

	mux.Get("/nba/teams", Repo.NBATeams)
	mux.Post("/nba/teams", Repo.PostNBATeams)
	mux.Get("/nba/teams-json", Repo.NBATeamsAvailabilityJSON)
	mux.Get("/nba/team-info-summary", Repo.NBATeamInfoSummary)

	mux.Get("/nba/results", Repo.NBAResults)
	mux.Post("/nba/results", Repo.PostNBAResults)

	fileServer := http.FileServer(http.Dir("../../static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// NoSurf adds CSFR protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csfrHandler := nosurf.New(next)

	csfrHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csfrHandler
}

// SessionLoad loads and saves session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// Get all the files named *.page.tmpl from ../../templates/
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// Range through all the files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = templateSet
	}

	return myCache, nil
}
