package main

import (
	"2K22/utilities"
)

func main() {
	scrapedData := utilities.ScrapeDataFromURL("https://www.2kratings.com/lists/top-100-highest-nba-2k-ratings")
	utilities.WriteToJson("player_data.json", scrapedData)
	utilities.ReadJson("player_data.json")
}
