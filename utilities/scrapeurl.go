package utilities

import (
	"strings"

	"github.com/gocolly/colly"
)

var sliceOfPlayers []PlayerData

// NBA player
type PlayerData struct {
	Rank      			int      `json:"rank"`
	Name      			string   `json:"name"`
	Positions 			[]string `json:"positions"`
	Team      			string   `json:"team"`
	Height    			[]int    `json:"feet_inch"`
	OverallRating      	int      `json:"overall_rating"`
	ThreePointRating	int      `json:"3pt_rating"`
	DunkRating      	int      `json:"dunk_rating"`
}

func ScrapeDataFromURL(scrapeUrl string) []PlayerData {
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
					playerPositions[i] = TrimSpace(playerPositions[i])
				}

				// Get player height from player info
				playerHeightData := strings.Split(playerInfo[1], "'")
				for i := range playerHeightData {
					playerHeightData[i] = TrimSpace(playerHeightData[i])
				}
				playerHeightFeet := playerHeightData[0]
				playerHeightInches := Trim(playerHeightData[1], "\"")
				playerHeight := []int{}
				playerHeight = append(playerHeight, Atoi(playerHeightFeet), Atoi(playerHeightInches))

				// Get player team from player info
				playerTeam := TrimSpace(playerInfo[2])

				// Get player ratings
				playerOverallRating := Atoi(el.ChildText("td:nth-child(3)"))
				player3ptRating := Atoi(el.ChildText("td:nth-child(4)"))
				playerDunkRating := Atoi(el.ChildText("td:nth-child(5)"))

				// Add data to struct
				player := PlayerData{
					Rank:     			i,
					Name:      			playerName,
					Positions: 			playerPositions,
					Team:      			playerTeam,
					Height:   			playerHeight,
					OverallRating:   	playerOverallRating,
					ThreePointRating:   player3ptRating,
					DunkRating:   		playerDunkRating,
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
