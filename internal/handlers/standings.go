package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/oskarpedosk/baltijas-kauss/internal/forms"
	"github.com/oskarpedosk/baltijas-kauss/internal/helpers"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"
)

func (m *Repository) Standings(w http.ResponseWriter, r *http.Request) {
	m.DB.NewSeason()
	if r.FormValue("action") == "add" {
		homeTeam, err := strconv.Atoi(r.FormValue("home_team_id"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		homeScore, err := strconv.Atoi(r.FormValue("home_score"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		awayScore, err := strconv.Atoi(r.FormValue("away_score"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		awayTeam, err := strconv.Atoi(r.FormValue("away_team_id"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		result := models.Result{
			HomeTeamID: homeTeam,
			HomeScore:  homeScore,
			AwayScore:  awayScore,
			AwayTeamID: awayTeam,
		}
		err = m.DB.AddResult(result)
		if err != nil {
			helpers.ServerError(w, err)
		}
	} else if r.FormValue("action") == "update" {
		homeTeam, err := strconv.Atoi(r.FormValue("home_team_id"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		homeScore, err := strconv.Atoi(r.FormValue("home_score"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		awayScore, err := strconv.Atoi(r.FormValue("away_score"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		awayTeam, err := strconv.Atoi(r.FormValue("away_team_id"))
		if err != nil {
			helpers.ServerError(w, err)
		}

		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		result := models.Result{
			HomeTeamID: homeTeam,
			HomeScore:  homeScore,
			AwayScore:  awayScore,
			AwayTeamID: awayTeam,
		}
		err = m.DB.UpdateResult(result)
		if err != nil {
			helpers.ServerError(w, err)
		}

	} else if r.FormValue("action") == "delete" {
		timestampString := r.FormValue("timestamp")
		layout := "2006-01-02 15:04:05 -0700 MST"
		timestamp, err := time.Parse(layout, timestampString)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		result := models.Result{
			CreatedAt: timestamp,
		}
		err = m.DB.DeleteResult(result)
		if err != nil {
			helpers.ServerError(w, err)
		}
	}

	var emptyStandings models.Result
	data := make(map[string]interface{})

	teams, err := m.DB.GetTeams()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	standings, err := m.DB.GetStandings()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	lastResults, err := m.DB.GetLastResults(10)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["result"] = emptyStandings
	data["teams"] = teams
	data["standings"] = standings
	data["last_results"] = lastResults

	render.Template(w, r, "standings.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostStandings(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	homeTeam, err := strconv.Atoi(r.Form.Get("home_team"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	homeScore, err := strconv.Atoi(r.Form.Get("home_score"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	awayScore, err := strconv.Atoi(r.Form.Get("away_score"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	awayTeam, err := strconv.Atoi(r.Form.Get("away_team"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	result := models.Result{
		HomeTeamID: homeTeam,
		HomeScore:  homeScore,
		AwayScore:  awayScore,
		AwayTeamID: awayTeam,
	}

	form := forms.New(r.PostForm)

	form.Required("home_team", "home_score", "away_score", "away_team")
	form.IsDuplicate("home_team", "away_team", "Home and away have to be different")
	form.IsDuplicate("home_score", "away_score", "Score can't be a draw")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["NBAresult"] = result

		render.Template(w, r, "standings.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.AddResult(result)
	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "result", result)

	http.Redirect(w, r, "/results", http.StatusSeeOther)
}
