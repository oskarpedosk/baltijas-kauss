package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// NBA player
type PlayerData struct {
	Rank      int      `json:"rank"`
	Name      string   `json:"name"`
	Positions []string `json:"positions"`
	Team      string   `json:"team"`
	Height    []int    `json:"height"`
	Ratings   []int    `json:"ratings"`
}

func main() {
	scrapeDataFromURL("https://www.2kratings.com/lists/top-100-highest-nba-2k-ratings")

}

func scrapeDataFromURL(scrapeUrl string) {
	// Create a file
	fName := "data.json"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Could not create file, err: %q", err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	c := colly.NewCollector()
	c.OnHTML("div.table-responsive tbody", func(e *colly.HTMLElement) {
		i := 1
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			if el.Text != "" {
				// Get player name
				playerName := el.ChildText("td:nth-child(2) div.entries span.entry-font")

				// Get player info (positions, height, team)
				playerInfo := strings.Split(el.ChildText("td:nth-child(2) div.entries span.entry-subtext-font"), "|")

				// Get player positions from player info
				playerPositions := strings.Split(playerInfo[0], "/")
				for i := range playerPositions {
					playerPositions[i] = trimSpace(playerPositions[i])
				}

				// Get player height from player info
				playerHeightData := strings.Split(playerInfo[1], "'")
				for i := range playerHeightData {
					playerHeightData[i] = trimSpace(playerHeightData[i])
				}
				playerHeightFeet := playerHeightData[0]
				playerHeightInches := trim(playerHeightData[1], "\"")
				playerHeight := []int{}
				playerHeight = append(playerHeight, atoi(playerHeightFeet), atoi(playerHeightInches))

				// Get player team from player info
				playerTeam := trimSpace(playerInfo[2])

				// Get player ratings
				playerRatings := []int{}
				playerOverallRating := atoi(el.ChildText("td:nth-child(3)"))
				player3ptRating := atoi(el.ChildText("td:nth-child(4)"))
				playerDunkRating := atoi(el.ChildText("td:nth-child(5)"))
				playerRatings = append(playerRatings, playerOverallRating, player3ptRating, playerDunkRating)

				// Add data to struct
				player := PlayerData{
					Rank: i,
					Name: playerName,
					Positions: playerPositions,
					Team: playerTeam,
					Height: playerHeight,
					Ratings: playerRatings,
				}
				i++
				fmt.Println(player)

			}
		})
		// fmt.Println("Scraping Complete")
	})
	c.Visit(scrapeUrl)
}

func trimSpace(str string) string {
	str = strings.TrimSpace(str)
	return str
}

func trim(str string, condition string) string {
	str = strings.Trim(str, condition)
	return str
}

func atoi(str string) int {
	integer, _ := strconv.Atoi(str)
	return integer
}
