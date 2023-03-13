package handlers

import (
	"net/http"
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

func (m *Repository) AdminNBAPlayers(w http.ResponseWriter, r *http.Request) {
	players, err := m.DB.GetPlayers(120, 0)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["nba_players"] = players

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
