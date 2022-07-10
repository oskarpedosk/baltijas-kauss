package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

// NBA player
type PlayerInfo struct {
	PlayerID 			string 
	PlayerName 			string
	PlayerPositions 	[]string
	PlayerTeam 			string
}
type Height struct {
	Feet 	int
	Inches 	int
}
type PlayerStats struct {
	PlayerHeight 		Height
	PlayerOverallRating int
	Player3ptRating 	int
	PlayerDunkRating 	int
}
//type Player struct {
//	info PlayerInfo
//	stats PlayerStats
//}

func trimSpace(str string) string {
	str = strings.TrimSpace(str)
	return str
}

func trim(str string, condition string) string {
	str = strings.Trim(str, condition)
	return str
}

func main() {
	// creating a file
	/*
    fName := "data.csv"
    file, err := os.Create(fName)
    if err != nil {
        log.Fatalf("Could not create file, err: %q", err)
        return
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()
	*/
	scrapeUrl := "https://www.2kratings.com/lists/top-100-highest-nba-2k-ratings"

	c := colly.NewCollector()
    c.OnHTML("div.table-responsive tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			if el.Text != "" {
				playerID := trim(el.ChildText("td.counter"), ".")
				playerName := el.ChildText("td:nth-child(2) div.entries span.entry-font")
				playerInfo := strings.Split(el.ChildText("td:nth-child(2) div.entries span.entry-subtext-font"), "|")
				playerHeightData := strings.Split(playerInfo[1], "'")
				for i := range playerHeightData {
					playerHeightData[i] = trimSpace(playerHeightData[i])
				}
				playerHeightFeet := playerHeightData[0]
				playerHeightInches := trim(playerHeightData[1], "\"")
				playerHeight := make([]string, 2)
				playerHeight = append(playerHeight, playerHeightFeet, playerHeightInches)
				playerTeam := trimSpace(playerInfo[2])
				playerPosition := strings.Split(playerInfo[0], "/")
				playerOverallRating := el.ChildText("td:nth-child(3)")
				player3ptRating := el.ChildText("td:nth-child(4)")
				playerDunkRating := el.ChildText("td:nth-child(5)")
				for i := range playerPosition {
					playerPosition[i] = trimSpace(playerPosition[i])
				}
				
				fmt.Println(playerID, playerName, playerHeight, playerPosition, playerOverallRating, player3ptRating, playerDunkRating, playerTeam)
			}
			
			
			/*
			if playerID != "" {
				writer.Write([]string {
					playerID,
					playerName,
					playerOverallRating,
					
				})
			}
			*/
		})
        
        fmt.Println("Scrapping Complete")	
    })
    c.Visit(scrapeUrl)
}