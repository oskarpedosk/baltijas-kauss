package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"reflect"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/oskarpedosk/baltijas-kauss/internal/forms"
	"github.com/oskarpedosk/baltijas-kauss/internal/helpers"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"
)

// mac, windows or ubuntu
const systemOS = "ubuntu"

var queryFilters = map[string]string{
	"team":   "TeamID",
	"ovrl":   "OverallMin",
	"ovrh":   "OverallMax",
	"hl":     "HeightMin",
	"hh":     "HeightMax",
	"wl":     "WeightMin",
	"wh":     "WeightMax",
	"3ptl":   "ThreePointShotMin",
	"3pth":   "ThreePointShotMax",
	"ddunkl": "DrivingDunkMin",
	"ddunkh": "DrivingDunkMax",
	"athl":   "AthleticismMin",
	"athh":   "AthleticismMax",
	"perdl":  "PerimeterDefenseMin",
	"perdh":  "PerimeterDefenseMax",
	"intdl":  "InteriorDefenseMin",
	"intdh":  "InteriorDefenseMax",
	"rebl":   "ReboundingMin",
	"rebh":   "ReboundingMax",
	"p1":     "Position1",
	"p2":     "Position2",
	"p3":     "Position3",
	"p4":     "Position4",
	"p5":     "Position5",
	"limit":  "Limit",
	"offset": "Offset",
}

// Player is the single player handler
func (m *Repository) Player(w http.ResponseWriter, r *http.Request) {
	playerID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	player, err := m.DB.GetPlayer(playerID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	badges, err := m.DB.GetPlayerBadges(playerID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	teams, err := m.DB.GetTeams()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var playersTeam models.Team
	for _, team := range teams {
		if team.TeamID == player.TeamID {
			playersTeam = team
			break
		}
	}

	ADP, err := m.DB.GetADP(playerID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["player"] = player
	data["ADP"] = ADP
	data["badges"] = badges
	data["team"] = playersTeam
	data["teams"] = teams[1:]
	data["FA"] = teams[0]

	render.Template(w, r, "player.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) PostPlayer(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	playerID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	teamID, err := strconv.Atoi(r.FormValue("team_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	err = m.DB.AddPlayer(playerID, teamID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

func (m *Repository) Players(w http.ResponseWriter, r *http.Request) {
	filter := models.Filter{
		TeamID:              0,
		HeightMin:           150,
		HeightMax:           250,
		WeightMin:           50,
		WeightMax:           150,
		OverallMin:          1,
		OverallMax:          99,
		ThreePointShotMin:   1,
		ThreePointShotMax:   99,
		DrivingDunkMin:      1,
		DrivingDunkMax:      99,
		AthleticismMin:      1,
		AthleticismMax:      99,
		PerimeterDefenseMin: 1,
		PerimeterDefenseMax: 99,
		InteriorDefenseMin:  1,
		InteriorDefenseMax:  99,
		ReboundingMin:       1,
		ReboundingMax:       99,
		Position1:           1,
		Position2:           1,
		Position3:           1,
		Position4:           1,
		Position5:           1,
		Limit:               20,
		Offset:              0,
		Era:                 2,
		Col1:                "overall",
		Col2:                "\"attributes/TotalAttributes\"",
		Order:               "desc",
	}

	for key, value := range r.URL.Query() {
		if key == "col" {
			if value[0] == "lname" {
				filter.Col1 = "last_name"
			} else if value[0] == "ovr" {
				filter.Col1 = "overall"
				filter.Col2 = "\"attributes/TotalAttributes\""
			} else if value[0] == "3pt" {
				filter.Col1 = "\"attributes/ThreePointShot\""
			} else if value[0] == "ddunk" {
				filter.Col1 = "\"attributes/DrivingDunk\""
			} else if value[0] == "ath" {
				filter.Col1 = "\"attributes/Athleticism\""
			} else if value[0] == "perd" {
				filter.Col1 = "\"attributes/PerimeterDefense\""
			} else if value[0] == "intd" {
				filter.Col1 = "\"attributes/InteriorDefense\""
			} else if value[0] == "reb" {
				filter.Col1 = "\"attributes/Rebounding\""
			} else if value[0] == "bdg" {
				filter.Col1 = "total_badges"
			} else if value[0] == "total" {
				filter.Col1 = "\"attributes/TotalAttributes\""
			}
			if filter.Col1 != "overall" {
				filter.Col2 = "overall"
			}
		} else if key == "sort" {
			filter.Order = value[0]
		} else if key == "search" {
			filter.Search = value[0]
		} else if key == "era" {
			if value[0] == "current" {
				filter.Era = 1
			} else if value[0] == "legends" {
				filter.Era = 0
			}
		} else {
			queryInt, err := strconv.Atoi(value[0])
			if err != nil {
				continue
			}
			if fieldName, ok := queryFilters[key]; ok {
				field := reflect.ValueOf(&filter).Elem().FieldByName(fieldName)
				field.SetInt(int64(queryInt))
			}
		}
	}

	teams, err := m.DB.GetTeams()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	players, err := m.DB.GetPlayers(filter)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var playersWithTeamInfo []models.PlayerWithTeamInfo

	for _, player := range players {
		for _, team := range teams {
			if player.TeamID == team.TeamID {
				playerWithTeamInfo := models.PlayerWithTeamInfo{
					Player: player,
					Team:   team,
				}
				playersWithTeamInfo = append(playersWithTeamInfo, playerWithTeamInfo)
				break
			}
		}
	}

	ranking := []int{}
	for i := 1; i <= len(players); i++ {
		ranking = append(ranking, i+filter.Offset)
	}

	data := make(map[string]interface{})
	data["players"] = playersWithTeamInfo
	data["teams"] = teams[1:]
	data["FA"] = teams[0]
	data["ranking"] = ranking

	render.Template(w, r, "players.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostPlayers(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	switch r.FormValue("action") {
	case "change_team":
		playerID, err := strconv.Atoi(r.FormValue("player_id"))
		if err != nil {
			helpers.ServerError(w, err)
		}

		teamID, err := strconv.Atoi(r.FormValue("team_id"))
		if err != nil {
			helpers.ServerError(w, err)
		}

		player := models.Player{
			PlayerID:         playerID,
			TeamID:           teamID,
			AssignedPosition: 0,
		}

		err = m.DB.SwitchTeam(player)
		if err != nil {
			helpers.ServerError(w, err)
		}
	}
}

func (m *Repository) PostUpdatePlayer(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	playerID := r.FormValue("player_id")
	ratingsURL := r.FormValue("ratings_url")

	if len(playerID) > 0 && len(ratingsURL) > 0 {
		go func(playerID, ratingsURL string) {
			filePath := ""
			switch systemOS {
			case "mac":
				filePath = "./static/js/script/updateplayer.js"
			case "windows":
				filePath = ".\\static\\js\\script\\updateplayer.js"
			case "ubuntu":
				filePath = "/var/www/bkauss/static/js/script/updateplayer.js"
			}

			cmd := exec.Command("node", filePath, systemOS, playerID, ratingsURL)
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("Command failed with error: %v\n", err)
			}

			// Parse the output as an array of two objects
			var data []json.RawMessage
			err = json.Unmarshal(output, &data)
			if err != nil {
				log.Println(data)
				log.Printf("Parsing scraper output err: %v\n", err)
				return
			}

			// Unmarshal the first object as a Player
			var player models.Player
			err = json.Unmarshal(data[0], &player)
			if err != nil {
				log.Println(data[0])
				log.Printf("Unmarshaling first object to Player err: %v\n", err)
				return
			}

			// Unmarshal the second object as a slice of Badges
			var badges []models.Badge
			err = json.Unmarshal(data[1], &badges)
			if err != nil {
				log.Println(data[1])
				log.Printf("Unmarshaling second object to slice of Badges err: %v\n", err)
				return
			}

			err = m.DB.UpdatePlayer(player)
			if err != nil {
				log.Printf("m.DB.UpdatePlayer err: %v\n", err)
				return
			}

			err = m.DB.UpdatePlayerBadges(player, badges)
			if err != nil {
				log.Printf("m.DB.UpdatePlayerBadges err: %v\n", err)
				return
			}
			log.Println(playerID, ratingsURL, "updated")
		}(playerID, ratingsURL)
	}

	m.App.Session.Put(r.Context(), "info", "Updating player ID: "+playerID)
	http.Redirect(w, r, "/players/"+playerID, http.StatusSeeOther)
}

func (m *Repository) SearchPlayers(w http.ResponseWriter, r *http.Request) {
	filter := models.Filter{
		TeamID:              0,
		HeightMin:           150,
		HeightMax:           250,
		WeightMin:           50,
		WeightMax:           150,
		OverallMin:          1,
		OverallMax:          99,
		ThreePointShotMin:   1,
		ThreePointShotMax:   99,
		DrivingDunkMin:      1,
		DrivingDunkMax:      99,
		AthleticismMin:      1,
		AthleticismMax:      99,
		PerimeterDefenseMin: 1,
		PerimeterDefenseMax: 99,
		InteriorDefenseMin:  1,
		InteriorDefenseMax:  99,
		ReboundingMin:       1,
		ReboundingMax:       99,
		Position1:           1,
		Position2:           1,
		Position3:           1,
		Position4:           1,
		Position5:           1,
		Limit:               10,
		Offset:              0,
		Era:                 2,
		Search:              r.URL.Query().Get("query"),
		Col1:                "overall",
		Col2:                "\"attributes/TotalAttributes\"",
		Order:               "desc",
	}

	players, err := m.DB.GetPlayers(filter)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	_ = helpers.WriteJson(w, http.StatusOK, players, nil)
}
