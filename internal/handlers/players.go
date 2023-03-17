package handlers

import (
	"encoding/json"
	"fmt"
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

	data := make(map[string]interface{})
	data["player"] = player
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
			filePath := "./static/js/script/updateplayer.js"
			cmd := exec.Command("node", filePath, playerID, ratingsURL)
			output, err := cmd.Output()
			if err != nil {
				fmt.Println(err)
				m.App.Session.Put(r.Context(), "warning", err)
				http.Redirect(w, r, "/players", http.StatusSeeOther)
				return
			}


			// Parse the output as an array of two objects
			var data []json.RawMessage
			err = json.Unmarshal(output, &data)
			if err != nil {
				fmt.Println(err)
				m.App.Session.Put(r.Context(), "warning", err)
				http.Redirect(w, r, "/players", http.StatusSeeOther)
				return
			}

			// Unmarshal the first object as a Player
			var player models.Player
			err = json.Unmarshal(data[0], &player)
			if err != nil {
				fmt.Println(err)
				m.App.Session.Put(r.Context(), "warning", err)
				http.Redirect(w, r, "/players", http.StatusSeeOther)
				return
			}

			// Unmarshal the second object as a slice of Badges
			var badges []models.Badge
			err = json.Unmarshal(data[1], &badges)
			if err != nil {
				fmt.Println(err)
			}

			err = m.DB.UpdatePlayer(player)
			if err != nil {
				helpers.ServerError(w, err)
			}

			err = m.DB.UpdatePlayerBadges(player, badges)
			if err != nil {
				helpers.ServerError(w, err)
			}
	}

	m.App.Session.Put(r.Context(), "warning", "Updating player ID: " + playerID)
	http.Redirect(w, r, "/players", http.StatusSeeOther)
}
