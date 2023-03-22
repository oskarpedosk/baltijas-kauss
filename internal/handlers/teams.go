package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/oskarpedosk/baltijas-kauss/internal/forms"
	"github.com/oskarpedosk/baltijas-kauss/internal/helpers"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"
)

func (m *Repository) Team(w http.ResponseWriter, r *http.Request) {
	teamID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	teams, err := m.DB.GetTeams()
	if err != nil {
		helpers.ServerError(w, err)
	}

	players, err := m.DB.GetTeamPlayers(teamID)
	if err != nil {
		helpers.ServerError(w, err)
	}

	data := make(map[string]interface{})
	data["players"] = players
	data["teams"] = teams[1:]
	data["starters"] = []string{"1", "2", "3", "4", "5"}

	for _, team := range teams {
		if team.TeamID == teamID {
			data["team"] = team
			break
		}
	}

	render.Template(w, r, "teams.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostTeam(w http.ResponseWriter, r *http.Request) {
	teamID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	err = r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	switch r.FormValue("action") {
	case "updateTeam":
		teamInfo := models.Team{
			TeamID:       teamID,
			Name:         r.FormValue("team_name"),
			Abbreviation: r.FormValue("abbreviation"),
			Color1:       r.FormValue("team_color1"),
			Color2:       r.FormValue("team_color2"),
			TextColor:    r.FormValue("text_color"),
		}

		form := forms.New(r.PostForm)

		form.Required("team_name", "abbreviation")
		form.AlphaNumeric("team_name", "abbreviation")
		form.MaxLength("team_name", 20)
		form.IsUpper("abbreviation")
		form.MaxLength("abbreviation", 4)

		if !form.Valid() {
			data := make(map[string]interface{})

			teams, err := m.DB.GetTeams()
			if err != nil {
				helpers.ServerError(w, err)
			}

			data["teams"] = teams[1:]

			for _, team := range teams {
				if team.TeamID == teamID {
					data["team"] = team
					break
				}
			}

			errMsg := form.Errors.Get("team_name")
			if errMsg == "" {
				errMsg = form.Errors.Get("abbreviation")
			}
			m.App.Session.Put(r.Context(), "error", errMsg)
			render.Template(w, r, "teams.page.tmpl", &models.TemplateData{
				Form: form,
				Data: data,
			})
			return
		}

		err = m.DB.UpdateTeam(teamInfo)
		if err != nil {
			helpers.ServerError(w, err)
		}

		m.App.Session.Put(r.Context(), "flash", "Team updated successfully!")
		http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
	case "updatePosition":
		playerID, err := strconv.Atoi(r.FormValue("player_id"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		position, err := strconv.Atoi(r.FormValue("position"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		removePlayerStr := r.FormValue("remove_player_id")

		if removePlayerStr != "" {
			removePlayer, err := strconv.Atoi(removePlayerStr)
			if err != nil {
				helpers.ServerError(w, err)
			}
			if removePlayer != playerID {
				m.DB.AssignPosition(removePlayer, 0)
			}
		}
		m.DB.AssignPosition(playerID, position)
	default:
		m.App.Session.Put(r.Context(), "warning", "Error, wrong post")
		http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
	}
}
