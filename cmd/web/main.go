package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/oskarpedosk/baltijas-kauss/internal/config"
	"github.com/oskarpedosk/baltijas-kauss/internal/handlers"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"
const nba2KDataFileName = "nba2k_player_data"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting application on port%s\n", portNumber)

	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = serve.ListenAndServe()
	log.Fatal(err)
}

func run() error {
	// What am I going to put in the session
	gob.Register(models.TeamInfo{})

	// Change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return err
	}

	app.TemplateCache = templateCache
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	return nil
}
