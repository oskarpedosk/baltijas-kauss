package main

import (
	"fmt"
	"net/http"


)

var portNumber = ":8080"

func main() {
	// scrapedData := utilities.ScrapeDataFromURL("https://www.2kratings.com/lists/top-100-highest-nba-2k-ratings")
	// utilities.WriteToJson("player_data.json", scrapedData)
	http.HandleFunc("/", Home)
	http.HandleFunc("/player_ratings", Players) 

	fmt.Printf("Starting application on port%s\n", portNumber)
	_ = http.ListenAndServe(portNumber, nil)
}

