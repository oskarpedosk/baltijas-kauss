package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/oskarpedosk/baltijas-kauss/pkg/config"
	"github.com/oskarpedosk/baltijas-kauss/pkg/handlers"
	"github.com/oskarpedosk/baltijas-kauss/pkg/render"
	"github.com/oskarpedosk/baltijas-kauss/utilities"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"
const nba2KDataFileName = "nba2k_player_data"

var app config.AppConfig
var session *scs.SessionManager

func main() {

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
	}

	app.TemplateCache = templateCache
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplate(&app)

	needsScraping := false
	if needsScraping {
		scrapedData := utilities.ScrapeNBA2KData()
		utilities.WriteToJson(nba2KDataFileName, scrapedData)
	}

	fmt.Printf("Starting application on port%s\n", portNumber)

	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = serve.ListenAndServe()
	log.Fatal(err)
}
