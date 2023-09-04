package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/oskarpedosk/baltijas-kauss/internal/helpers"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"
)

func (m *Repository) AdminHome(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) AdminTeams(w http.ResponseWriter, r *http.Request) {
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

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if limit == 0 {
		filter.Limit = playerCount
	} else if 0 < limit && limit <= playerCount {
		filter.Limit = limit
	}

	offset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	filter.Offset = offset

	players, err := m.DB.GetPlayers(filter)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	filePath := ""
	switch systemOS {
	case "mac":
		filePath = "./static/js/script/updateplayer.js"
	case "windows":
		filePath = ".\\static\\js\\script\\updateplayer.js"
	case "ubuntu":
		filePath = "/var/www/bkauss/static/js/script/updateplayer.js"
	}

	success := 0
	for _, player := range players {
		successMsg := fmt.Sprintf("%s %s update success", player.FirstName, player.LastName)
		failMsg := fmt.Sprintf("%s %s update failed", player.FirstName, player.LastName)
		cmd := exec.Command("node", filePath, strconv.Itoa(player.PlayerID), player.RatingsURL)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(failMsg)
			log.Printf("Command failed with error: %v\n", err)
			continue
		}

		// Parse the output as an array of two objects
		var data []json.RawMessage
		err = json.Unmarshal(output, &data)
		if err != nil {
			log.Println(data)
			log.Printf("Parsing scraper output err: %v\n", err)
			continue
		}

		// Unmarshal the first object as a Player
		var player models.Player
		err = json.Unmarshal(data[0], &player)
		if err != nil {
			log.Println(data[0])
			log.Printf("Unmarshaling first object to Player err: %v\n", err)
			continue
		}

		// Unmarshal the second object as a slice of Badges
		var badges []models.Badge
		err = json.Unmarshal(data[1], &badges)
		if err != nil {
			log.Println(data[1])
			log.Printf("Unmarshaling second object to slice of Badges err: %v\n", err)
			continue
		}

		err = m.DB.UpdatePlayer(player)
		if err != nil {
			log.Println(failMsg)
			log.Printf("m.DB.UpdatePlayer err: %v\n", err)
			continue
		}

		err = m.DB.UpdatePlayerBadges(player, badges)
		if err != nil {
			log.Println(failMsg)
			log.Printf("m.DB.UpdatePlayerBadges err: %v\n", err)
			continue
		}
		success++
		log.Println(successMsg)
	}
	log.Println("------------")
	msg := fmt.Sprintf("%d of %d players updated successfully", success, filter.Limit)
	log.Println(msg)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (m *Repository) AdminPlayers(w http.ResponseWriter, r *http.Request) {
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
		Limit:               1000,
		Offset:              0,
		Era:                 2,
		Col1:                "overall",
		Col2:                "\"attributes/TotalAttributes\"",
		Order:               "desc",
	}

	players, err := m.DB.GetPlayers(filter)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["players"] = players

	render.Template(w, r, "admin-nba-players.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminStandings(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-nba-standings.page.tmpl", &models.TemplateData{})
}

func (m *Repository) PostAdminStandings(w http.ResponseWriter, r *http.Request) {
	m.DB.StartNewSeason()
	m.App.Session.Put(r.Context(), "flash", "New season started")
	http.Redirect(w, r, "/standings", http.StatusSeeOther)
}

func (m *Repository) AdminPlayer(w http.ResponseWriter, r *http.Request) {
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
	data["player"] = player
	data["teams"] = teams

	render.Template(w, r, "admin-nba-player.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) NewPlayer(w http.ResponseWriter, r *http.Request) {
	ratingsURL := r.FormValue("ratings_url")
	err := m.DB.CreateNewPlayer(ratingsURL)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot create player.")
		http.Redirect(w, r, "/admin/players", http.StatusSeeOther)
	}
	m.App.Session.Put(r.Context(), "flash", "Player successfully created!")
	http.Redirect(w, r, "/admin/players", http.StatusSeeOther)
}

func (m *Repository) EditPlayer(w http.ResponseWriter, r *http.Request) {
	playerID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	method := r.Form.Get("method")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	switch method {
	case "edit":
		log.Println(r.Form)
	case "delete":
		err = m.DB.DeletePlayer(playerID)
		if err != nil {
			m.App.Session.Put(r.Context(), "error", fmt.Sprintf("Error deleting playerID: %d", playerID))
			http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
			return
		}
		m.App.Session.Put(r.Context(), "flash", "Player deleted!")
		http.Redirect(w, r, "/admin/players", http.StatusSeeOther)
	}
}
