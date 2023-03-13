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

	data := make(map[string]interface{})
	data["teams"] = teams[1:]

	for _, team := range teams {
		if team.TeamID == teamID {
			data["team"] = team
			break
		}
	}

	render.Template(w, r, "team.page.tmpl", &models.TemplateData{
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

	text := r.FormValue("dark_text")
	if text == "" {
		text = "false"
	}

	teamInfo := models.Team{
		TeamID:       teamID,
		Name:         r.FormValue("team_name"),
		Abbreviation: r.FormValue("abbreviation"),
		Color1:       r.FormValue("team_color1"),
		Color2:       r.FormValue("team_color2"),
		DarkText:     text,
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
		render.Template(w, r, "team.page.tmpl", &models.TemplateData{
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
}
