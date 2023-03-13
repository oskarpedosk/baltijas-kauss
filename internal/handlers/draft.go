package handlers

import (
	"net/http"

	"github.com/oskarpedosk/baltijas-kauss/internal/helpers"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"
)

func (m *Repository) NBADraft(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	players, err := m.DB.GetPlayers(200, 0)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	teams, err := m.DB.GetTeams()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["nba_players"] = players
	data["nba_teams"] = teams
	data["positions"] = positions

	render.Template(w, r, "draft.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
