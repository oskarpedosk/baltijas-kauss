package utilities

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

var sliceOfNBAPlayers []NBAPlayerData

// NBA player
type NBAPlayerData struct {
	Name           string         `json:"name"`
	Team           string         `json:"team"`
	Archetype      string         `json:"archetype"`
	Positions      []string       `json:"positions"`
	Height         int            `json:"height"`
	Weight         int            `json:"weight"`
	BadgeCount     []int          `json:"badge_count"`
	Stats          map[string]int `json:"stats"`
	PlayerURL      string         `json:"player_url"`
	PlayerImageURL string         `json:"player_image_url"`
}

func ScrapeNBA2KData(scrapeUrl string) []NBAPlayerData {
	c := colly.NewCollector(
		colly.AllowedDomains("www.2kratings.com"),
	)
	c.OnHTML("div.table-responsive tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			if el.Text != "" {
				teamURL := el.ChildAttr("a", "href")
				sliceOfNBAPlayers = scrapePlayerURLFromTeam(teamURL)
			}
		})
		fmt.Println("Scraping Complete")
	})
	c.Visit(scrapeUrl)

	return sliceOfNBAPlayers
}

func scrapePlayerURLFromTeam(teamURL string) []NBAPlayerData {
	c := colly.NewCollector(
		colly.AllowedDomains("www.2kratings.com"),
	)
	firstTable := true
	c.OnHTML("table.table tbody", func(e *colly.HTMLElement) {
		if firstTable {
			e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
				if el.Text != "" {
					playerURL := el.ChildAttr("a", "href")
					sliceOfNBAPlayers = scrapePlayerStats(playerURL)
				}
			})
			firstTable = false
			fmt.Println("Players scraping Complete")
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())

	})
	c.Visit(teamURL)

	return sliceOfNBAPlayers
}

func scrapePlayerStats(playerURL string) []NBAPlayerData {
	c := colly.NewCollector(
		colly.AllowedDomains("www.2kratings.com"),
	)
	c.OnHTML("div.main", func(e *colly.HTMLElement) {
		// Get player info
		imageURL := ""
		playerName := ""
		playerTeam := ""
		playerArchetype := ""
		playerPositionsSlice := []string{}
		playerHeightInt := 0
		playerWeightInt := 0
		// Get player image URL
		e.ForEach("div.profile-photo", func(_ int, el *colly.HTMLElement) {
			imageURL = el.ChildAttr("img.header-image", "src")
		})
		
		e.ForEach("div.player-info", func(_ int, el *colly.HTMLElement) {
			if el.Text != "" {
				// Get player name
				playerName = el.ChildText("h1:nth-child(1)")

				// Get player team, remove "Team: " from the beginning of string
				playerTeam = strings.Split(el.ChildText("p:nth-child(3)"), ": ")[1]

				// Get player archetype, remove "Archetype: " from the beginning of string
				playerArchetype = strings.Split(el.ChildText("p:nth-child(4)"), ": ")[1]
				// Get player positions in a slice

				playerPositionsString := strings.Split(el.ChildText("p:nth-child(5)"), ": ")[1]
				playerPositionsSlice = strings.Split(playerPositionsString, "/")
				for i := range playerPositionsSlice {
					playerPositionsSlice[i] = TrimSpace(playerPositionsSlice[i])
				}

				heightAndWeightExists := true
				if len(strings.Split(el.ChildText("p:nth-child(6)"), "|")) != 2 {
					heightAndWeightExists = false
				}

				// Get player height
				playerHeightString := strings.Split(el.ChildText("p:nth-child(6)"), "|")[0]
				playerHeight := ""
				for _, char := range playerHeightString {
					if char == '(' {
						playerHeight = ""
					} else if char == 'c' {
						break
					} else {
						playerHeight += string(char)
					}
				}
				// Convert to int
				playerHeightInt = Atoi(playerHeight)

				// Get player weight
				if heightAndWeightExists {
					playerWeightString := strings.Split(el.ChildText("p:nth-child(6)"), "|")[1]
					playerWeight := ""
					for _, char := range playerWeightString {
						if char == '(' {
							playerWeight = ""
						} else if char == 'k' {
							break
						} else {
							playerWeight += string(char)
						}
					}
					// Convert to int
					playerWeightInt = Atoi(playerWeight)
				}
			}
		})

		// Get player overall rating
		playerStats := make(map[string]int)
		e.ForEach("span.attribute-box-player", func(_ int, el *colly.HTMLElement) {
			if el.Text != "" {
				overallRating := Atoi(el.Text)
				playerStats["Overall"] = overallRating
			}
		})

		// Add player badge counts to slice
		badgeCount := []int{}
		e.ForEach("div.badges-container span.badge-count", func(_ int, el *colly.HTMLElement) {
			if el.Text != "" {
				badgeCount = append(badgeCount, Atoi(el.Text))
			}
		})

		// Get player stats
		e.ForEach("div[id=nav-attributes] div.card-body li.mb-1", func(_ int, el *colly.HTMLElement) {
			if el.Text != "" {
				statsRating := Atoi(strings.SplitN(el.Text, " ", 2)[0])
				statsName := strings.ReplaceAll(strings.SplitN(el.Text, " ", 2)[1], " ", "_")
				statsName = strings.ReplaceAll(statsName, "-", "_")
				playerStats[statsName] = statsRating
			}
		})
		// Get player intangibles, potential and total attributes
		e.ForEach("div[id=nav-attributes] div.card-horizontal h5.card-title", func(_ int, el *colly.HTMLElement) {
			if el.Text != "" {
				// Get total attributes without a comma
				statsRating := Atoi(strings.Replace((strings.SplitN(el.Text, " ", 2)[0]), ",", "", -1))
				statsName := strings.ReplaceAll(strings.SplitN(el.Text, " ", 2)[1], " ", "_")
				playerStats[statsName] = statsRating
			}
		})
		fmt.Println("Player: " + playerName + " scraping Complete")

		// Add data to struct
		nbaPlayer := NBAPlayerData{
			Name:           playerName,
			Team:           playerTeam,
			Archetype:      playerArchetype,
			Positions:      playerPositionsSlice,
			Height:         playerHeightInt,
			Weight:         playerWeightInt,
			BadgeCount:     badgeCount,
			Stats:          playerStats,
			PlayerURL:      playerURL,
			PlayerImageURL: imageURL,
		}

		sliceOfNBAPlayers = append(sliceOfNBAPlayers, nbaPlayer)

	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())

	})
	c.Visit(playerURL)

	return sliceOfNBAPlayers
}
