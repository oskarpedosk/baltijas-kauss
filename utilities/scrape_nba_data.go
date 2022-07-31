package utilities

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

var sliceOfNBAPlayers []NBAPlayerData

// NBA player
type NBAPlayerData struct {
	FirstName      string         `json:"first_name"`
	LastName       string         `json:"last_name"`
	Team           string         `json:"team"`
	Archetype      string         `json:"archetype"`
	Positions      []string       `json:"positions"`
	Height         int            `json:"height"`
	Weight         int            `json:"weight"`
	Stats          map[string]int `json:"stats"`
	PlayerURL      string         `json:"player_url"`
	PlayerImageURL string         `json:"player_image_url"`
	BadgeCount     []int          `json:"badge_count"`
	Badges         []SingleBadge  `json:"badges"`
}

// NBA player badge
type SingleBadge struct {
	BadgeImageURL string
	Name          string
	Type          string
	Info          string
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
			fmt.Println("---------------------------")
			fmt.Println(teamURL + " scraping complete.")
			fmt.Println("")
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
		firstName := ""
		lastName := ""
		team := ""
		archetype := ""
		positionsSlice := []string{}
		heightInt := 0
		weightInt := 0
		// Get player image URL
		e.ForEach("div.profile-photo", func(_ int, el *colly.HTMLElement) {
			imageURL = el.ChildAttr("img.header-image", "src")
		})

		e.ForEach("div.player-info", func(_ int, el *colly.HTMLElement) {
			if el.Text != "" {
				// Get player name
				fullName := strings.ReplaceAll(el.ChildText("h1:nth-child(1)"), "’", "'")
				nameString := strings.SplitN(fullName, " ", 2)
				firstName = nameString[0]
				lastName = nameString[1]

				// Get player team, remove "Team: " from the beginning of string
				team = strings.Split(el.ChildText("p:nth-child(3)"), ": ")[1]

				// Get player archetype, remove "Archetype: " from the beginning of string
				archetype = strings.Split(el.ChildText("p:nth-child(4)"), ": ")[1]
				// Get player positions in a slice

				playerPositionsString := strings.Split(el.ChildText("p:nth-child(5)"), ": ")[1]
				positionsSlice = strings.Split(playerPositionsString, "/")
				for i := range positionsSlice {
					positionsSlice[i] = TrimSpace(positionsSlice[i])
				}

				heightAndWeightExists := true
				if len(strings.Split(el.ChildText("p:nth-child(6)"), "|")) != 2 {
					heightAndWeightExists = false
				}

				// Get player height
				heightString := strings.Split(el.ChildText("p:nth-child(6)"), "|")[0]
				height := ""
				for _, char := range heightString {
					if char == '(' {
						height = ""
					} else if char == 'c' {
						break
					} else {
						height += string(char)
					}
				}
				// Convert to int
				heightInt = Atoi(height)

				// Get player weight
				if heightAndWeightExists {
					weightString := strings.Split(el.ChildText("p:nth-child(6)"), "|")[1]
					weight := ""
					for _, char := range weightString {
						if char == '(' {
							weight = ""
						} else if char == 'k' {
							break
						} else {
							weight += string(char)
						}
					}
					// Convert to int
					weightInt = Atoi(weight)
				}
			}
		})

		// Get player overall rating
		playerStats := make(map[string]int)
		e.ForEach("span.attribute-box-player", func(_ int, el *colly.HTMLElement) {
			if el.Text != "" {
				overallRating := Atoi(el.Text)
				playerStats["Overall_Rating"] = overallRating
			}
		})

		// Add player badge counts to slice
		badgeCount := []int{}
		e.ForEach("div.badges-container span.badge-count", func(_ int, el *colly.HTMLElement) {
			if el.Text != "" {
				badgeCount = append(badgeCount, Atoi(el.Text))
			}
		})

		// Get table header stats, intangibles, potential and total attributes
		e.ForEach("div[id=nav-attributes] div.card-header", func(_ int, el *colly.HTMLElement) {
			if el.Text != "" {
				statsRating := Atoi(strings.Replace((strings.SplitN(el.Text, " ", 2)[0]), ",", "", -1))
				statsName := strings.ReplaceAll(strings.SplitN(el.Text, " ", 2)[1], " ", "_")
				playerStats[statsName] = statsRating
			}
		})
		// Get table body stats
		e.ForEach("div[id=nav-attributes] div.card-body li.mb-1", func(_ int, el *colly.HTMLElement) {
			if el.Text != "" {
				statsRating := Atoi(strings.SplitN(el.Text, " ", 2)[0])
				statsName := strings.ReplaceAll(strings.SplitN(el.Text, " ", 2)[1], " ", "_")
				statsName = strings.ReplaceAll(statsName, "-", "_")
				playerStats[statsName] = statsRating
			}
		})

		// Get player badges
		allBadges := []SingleBadge{}
		e.ForEach("div[id=pills-all] div.badge-card", func(_ int, el *colly.HTMLElement) {
			badgeImageURL := ""
			badgeName := ""
			badgeType := ""
			badgeInfo := ""
			if el.Text != "" {
				badgeImageURL = "https://www.2kratings.com" + el.ChildAttr("img", "data-src")
				badgeName = el.ChildText("h4")
				badgeType = el.ChildText("span.badge")
				badgeInfo = strings.ReplaceAll(el.ChildText("p.badge-description"), "’", "'")
				
				singleBadge := SingleBadge{
					BadgeImageURL: badgeImageURL,
					Name: badgeName,
					Type: badgeType,
					Info: badgeInfo,
				}
				allBadges = append(allBadges, singleBadge)
			}
		})


		fmt.Println("Player: " + firstName + " " + lastName + " scraping Complete")

		// Add data to struct
		nbaPlayer := NBAPlayerData{
			FirstName:      firstName,
			LastName:       lastName,
			Team:           team,
			Archetype:      archetype,
			Positions:      positionsSlice,
			Height:         heightInt,
			Weight:         weightInt,
			Stats:          playerStats,
			PlayerURL:      playerURL,
			PlayerImageURL: imageURL,
			BadgeCount:     badgeCount,
			Badges:         allBadges,
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
