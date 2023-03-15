package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/oskarpedosk/baltijas-kauss/internal/helpers"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"
)

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

func (m *Repository) AdminNBATeams(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-nba-teams.page.tmpl", &models.TemplateData{})
}

func (m *Repository) PostAdminPlayers(w http.ResponseWriter, r *http.Request) {
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

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	playerCount, err := m.DB.CountPlayers()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	limitStr := r.FormValue("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if limit == 0 {
		filter.Limit = playerCount
	} else if 0 < limit && limit <= playerCount {
		filter.Limit = limit
	}

	players, err := m.DB.GetPlayers(filter)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	filePath := "./static/js/script/updateplayer.js"
	success := 0
	for _, player := range players {
		cmd := exec.Command("node", filePath, strconv.Itoa(player.PlayerID), player.RatingsURL)
		output, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Parse the output as an array of two objects
		var data []json.RawMessage
		err = json.Unmarshal(output, &data)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Unmarshal the first object as a Player
		var player models.Player
		err = json.Unmarshal(data[0], &player)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Unmarshal the second object as a slice of Badges
		var badges []models.Badge
		err = json.Unmarshal(data[1], &badges)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = m.DB.UpdatePlayer(player)
		if err != nil {
			continue
			helpers.ServerError(w, err)
		}

		err = m.DB.UpdatePlayerBadges(player, badges)
		if err != nil {
			continue
			helpers.ServerError(w, err)
		}
		success++
	}

	msg := fmt.Sprintf("%d of %d updated successfully", success, filter.Limit)
	m.App.Session.Put(r.Context(), "warning", msg)
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func (m *Repository) AdminNBAPlayers(w http.ResponseWriter, r *http.Request) {
	// players, err := m.DB.GetPlayers(120, 0)
	// if err != nil {
	// 	helpers.ServerError(w, err)
	// 	return
	// }

	data := make(map[string]interface{})
	// data["nba_players"] = players

	render.Template(w, r, "admin-nba-players.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminNBAResults(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-nba-results.page.tmpl", &models.TemplateData{})
}

// Shows a single players stats
func (m *Repository) AdminShowNBAPlayer(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[3])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	player, err := m.DB.GetPlayer(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	teams, err := m.DB.GetTeams()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["nba_player"] = player
	data["nba_teams"] = teams

	render.Template(w, r, "admin-nba-player.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
