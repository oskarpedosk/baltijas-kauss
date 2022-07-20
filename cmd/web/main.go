package main

import (
	"2K22/pkg/config"
	"2K22/pkg/handlers"
	"2K22/pkg/render"
	"fmt"
	"log"
	"net/http"
)

var portNumber = ":8080"

func main() {
	var app config.AppConfig

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = templateCache

	render.NewTemplate(&app)

	// scrapedData := utilities.ScrapeDataFromURL("https://www.2kratings.com/lists/top-100-highest-nba-2k-ratings")
	// utilities.WriteToJson("player_data.json", scrapedData)
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/player_ratings", handlers.Players)

	fmt.Printf("Starting application on port%s\n", portNumber)
	_ = http.ListenAndServe(portNumber, nil)
}
