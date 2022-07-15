package main

import (
	"2K22/utilities"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// NBA player
type PlayerData struct {
	Rank      int      `json:"player_rank"`
	Name      string   `json:"player_name"`
	Positions []string `json:"player_positions"`
	Team      string   `json:"player_team"`
	Height    []int    `json:"player_height"`
	Ratings   []int    `json:"player_ratings"`
}

var sliceOfPlayers []PlayerData


func main() {
	scrapedData := scrapeDataFromURL("https://www.2kratings.com/lists/top-100-highest-nba-2k-ratings")
	utilities.WriteToJson("data.json", scrapedData)
}

func scrapeDataFromURL(scrapeUrl string) []PlayerData {
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
				
				// Add struct to slice
				sliceOfPlayers = append(sliceOfPlayers, player)
			}
		})
		// fmt.Println("Scraping Complete")
	})
	c.Visit(scrapeUrl)

	return sliceOfPlayers
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
