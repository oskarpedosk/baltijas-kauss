package handlers

import (
	"net/http"

	"github.com/oskarpedosk/baltijas-kauss/internal/helpers"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"
)

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	seasons, err := m.DB.GetSeasons()
	if err != nil {
		helpers.ServerError(w, err)
	}

	results, err := m.DB.GetSeasonResults(0)
	if err != nil {
		helpers.ServerError(w, err)
	}

	teams, err := m.DB.GetTeams()
	if err != nil {
		helpers.ServerError(w, err)
	}

	var teamsAndPlayers = []models.TeamAndPlayers{}
	for _, team := range teams {
		if team.TeamID != 1 {
			players, err := m.DB.GetTeamPlayers(team.TeamID)
			if err != nil {
				helpers.ServerError(w, err)
			}
			var teamAndPlayers = models.TeamAndPlayers{
				Team:    team,
				Players: players,
			}
			teamsAndPlayers = append(teamsAndPlayers, teamAndPlayers)
		}
	}
	standings := CalculateStandings(teams[1:], results)

	data := make(map[string]interface{})
	data["teams"] = teamsAndPlayers
	data["standings"] = standings
	data["activeSeason"] = seasons[0].SeasonID

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
