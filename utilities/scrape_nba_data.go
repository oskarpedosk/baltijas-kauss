package utilities

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

var NBAPlayersSlice []NBAPlayerData
var playersBadgesSlice []PlayersBadges
var badgesSlice []Badges
var playerIndex int = 1
var badgeIndex int = 1

const nba2KDataURL = "https://www.2kratings.com/current-teams"

// NBA player
type NBAPlayerData struct {
	PlayerID       int             `json:"player_id"`
	FirstName      string          `json:"first_name"`
	LastName       string          `json:"last_name"`
	Team           string          `json:"team"`
	Archetype      string          `json:"archetype"`
	Positions      []string        `json:"positions"`
	Height         int             `json:"height"`
	Weight         int             `json:"weight"`
	Stats          map[string]int  `json:"stats"`
	PlayerURL      string          `json:"player_url"`
	PlayerImageURL string          `json:"player_image_url"`
	BadgeCount     []int           `json:"badge_count"`
	Badges         []PlayersBadges `json:"badges"`
}

// NBA player badge
type PlayersBadges struct {
	PlayerID int    `json:"player_id"`
	BadgeID  int    `json:"badge_id"`
	Name     string `json:"name"`
	Level    string `json:"level"`
}

// NBA badges
type Badges struct {
	BadgeID        int    `json:"badge_id"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	Info           string `json:"info"`
	BronzeImageURL string `json:"bronze_img_url"`
	SilverImageURL string `json:"silver_img_url"`
	GoldImageURL   string `json:"gold_img_url"`
	HOFImageURL    string `json:"hof_img_url"`
}

func ScrapeNBA2KData() []NBAPlayerData {
	c := colly.NewCollector(
		colly.AllowedDomains("www.2kratings.com"),
	)

	// Get all badges
	badgesSlice = scrapeBadges()

	// Scrape all teams for player urls
	c.OnHTML("div.table-responsive tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			if el.Text != "" {
				// Get team URL for scraping
				teamURL := el.ChildAttr("a", "href")
				// Scrape players from a team
				NBAPlayersSlice = scrapePlayerURLFromTeam(teamURL)
			}
		})
		fmt.Println("Scraping Complete")
	})
	c.Visit(nba2KDataURL)

	return NBAPlayersSlice
}

func scrapePlayerURLFromTeam(teamURL string) []NBAPlayerData {
	c := colly.NewCollector(
		colly.AllowedDomains("www.2kratings.com"),
	)

	// Scrape only first table
	firstTable := true
	c.OnHTML("table.table tbody", func(e *colly.HTMLElement) {
		if firstTable {
			e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
				if el.Text != "" {
					// Get player URL for scraping
					playerURL := el.ChildAttr("a", "href")
					// Scrape player info and stats as a slice
					NBAPlayersSlice = scrapePlayerStats(playerURL)
					playerIndex++
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

	return NBAPlayersSlice
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

				// Check if player weight exists
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

		// Add player badge counts to badgeCount slice
		badgeCount := []int{}
		e.ForEach("div.badges-container span.badge-count", func(_ int, el *colly.HTMLElement) {
			if el.Text != "" {
				badgeCount = append(badgeCount, Atoi(el.Text))
			}
		})

		// Get table header stats, intangibles, potential and total attributes and add them to a playerStats map
		e.ForEach("div#nav-attributes div.card-header", func(_ int, el *colly.HTMLElement) {
			if el.Text != "" {
				statsRating := Atoi(strings.Replace((strings.SplitN(el.Text, " ", 2)[0]), ",", "", -1))
				statsName := strings.ReplaceAll(strings.SplitN(el.Text, " ", 2)[1], " ", "_")
				playerStats[statsName] = statsRating
			}
		})
		// Get table body stats and add them to a playerStats map
		e.ForEach("div#nav-attributes div.card-body li.mb-1", func(_ int, el *colly.HTMLElement) {
			if el.Text != "" {
				statsRating := Atoi(strings.SplitN(el.Text, " ", 2)[0])
				statsName := strings.ReplaceAll(strings.SplitN(el.Text, " ", 2)[1], " ", "_")
				statsName = strings.ReplaceAll(statsName, "-", "_")
				playerStats[statsName] = statsRating
			}
		})

		// Get single player badges
		playerBadges := []PlayersBadges{}
		e.ForEach("div[id=pills-all] div.badge-card", func(_ int, el *colly.HTMLElement) {
			badgeImageURL := ""
			badgeName := ""
			badgeType := ""
			badgeInfo := ""
			level := ""
			// Get player badges
			if el.Text != "" {
				badgeImageURL = "https://www.2kratings.com" + el.ChildAttr("img", "data-src")
				badgeName = el.ChildText("h4")
				badgeType = el.ChildText("span.badge")
				badgeInfo = strings.ReplaceAll(el.ChildText("p.badge-description"), "’", "'")
				tempBadgeIndex := 0
				// Add badge info to badges struct
				if strings.Contains(badgeImageURL, "_bronze") {
					level = "Bronze"
					for i, value := range badgesSlice {
						if value.Name == badgeName {
							value.Type = badgeType
							value.Info = badgeInfo
							value.BronzeImageURL = badgeImageURL
							tempBadgeIndex = i + 1
						}
					}
					
				} else if strings.Contains(badgeImageURL, "_silver") {
					level = "Silver"
					for i, value := range badgesSlice {
						if value.Name == badgeName {
							value.Type = badgeType
							value.Info = badgeInfo
							value.SilverImageURL = badgeImageURL
							tempBadgeIndex = i + 1
						}
					}
				} else if strings.Contains(badgeImageURL, "_gold") {
					level = "Gold"
					for i, value := range badgesSlice {
						if value.Name == badgeName {
							value.Type = badgeType
							value.Info = badgeInfo
							value.GoldImageURL = badgeImageURL
							tempBadgeIndex = i + 1
						}
					}
				} else if strings.Contains(badgeImageURL, "_hof") {
					level = "HOF"
					for i, value := range badgesSlice {
						if value.Name == badgeName {
							value.Type = badgeType
							value.Info = badgeInfo
							value.HOFImageURL = badgeImageURL
							tempBadgeIndex = i + 1
						}
					}
				}

				// Add data to PlayersBadges struct
				singleBadge := PlayersBadges{
					PlayerID: playerIndex,
					BadgeID:  tempBadgeIndex,
					Name:     badgeName,
					Level:    level,
				}
				// Append data to slices
				playerBadges = append(playerBadges, singleBadge)
				playersBadgesSlice = append(playersBadgesSlice, singleBadge)
			}
		})

		fmt.Println("Player: " + firstName + " " + lastName + " scraping Complete")

		// Add data to NBAPlayerData struct
		nbaPlayer := NBAPlayerData{
			PlayerID:       playerIndex,
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
			Badges:         playerBadges,
		}

		// Append players data to all players slice
		NBAPlayersSlice = append(NBAPlayersSlice, nbaPlayer)

	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())

	})
	c.Visit(playerURL)

	return NBAPlayersSlice
}

func scrapeBadges() []Badges {
	c := colly.NewCollector(
		colly.AllowedDomains("www.2kratings.com"),
	)

	// Scrape all badges names
	c.OnHTML("nav#sidebar]", func(e *colly.HTMLElement) {
		e.ForEach("li.sidebar-item", func(_ int, el *colly.HTMLElement) {
			if el.Text != "" {
				// Get badge name
				badgeName := el.Text
				// Add badge name and BadgeID to badges struct
				singleBadge := Badges{
					BadgeID: badgeIndex,
					Name:    badgeName,
				}
				badgeIndex++
				// Append badge to all badges slice
				badgesSlice = append(badgesSlice, singleBadge)
			}
		})
		fmt.Println("Badges scraping Complete")
	})

	c.Visit(nba2KDataURL)

	return badgesSlice
}
