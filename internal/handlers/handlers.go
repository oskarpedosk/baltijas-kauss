package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/oskarpedosk/baltijas-kauss/internal/config"
	"github.com/oskarpedosk/baltijas-kauss/internal/driver"
	"github.com/oskarpedosk/baltijas-kauss/internal/forms"
	"github.com/oskarpedosk/baltijas-kauss/internal/helpers"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"
	"github.com/oskarpedosk/baltijas-kauss/internal/repository"
	"github.com/oskarpedosk/baltijas-kauss/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository
var positions = []models.NBAPosition{
	{
		Name:   "PG",
		Number: 1,
	},
	{
		Name:   "SG",
		Number: 2,
	},
	{
		Name:   "SF",
		Number: 3,
	},
	{
		Name:   "PF",
		Number: 4,
	},
	{
		Name:   "C",
		Number: 5,
	},
	{
		Name:   "PG",
		Number: 6,
	},
	{
		Name:   "SG",
		Number: 7,
	},
	{
		Name:   "SF",
		Number: 8,
	},
	{
		Name:   "PF",
		Number: 9,
	},
	{
		Name:   "C",
		Number: 10,
	},
	{
		Name:   "Res",
		Number: 11,
	},
	{
		Name:   "Res",
		Number: 12,
	},
}

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) SignIn(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "signin.page.tmpl", &models.TemplateData{})
}

func (m *Repository) NBAHome(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, r, "nba_home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) NBAPlayers(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("action") == "update" {
		playerID, err := strconv.Atoi(r.FormValue("player_id"))
		if err != nil {
			helpers.ServerError(w, err)
		}

		nullInt := true

		teamID, err := strconv.Atoi(r.FormValue("team_id"))
		if err != nil {
			log.Println(err)
			nullInt = false
			// helpers.ServerError(w, err)
		}

		player := models.NBAPlayer{
			PlayerID: playerID,
			TeamID:   sql.NullInt64{int64(teamID), nullInt},
			Assigned: 0,
		}

		err = m.DB.UpdateNBAPlayer(player)
		if err != nil {
			helpers.ServerError(w, err)
		}
	}

	teams, err := m.DB.GetNBATeamInfo()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	players, err := m.DB.GetNBAPlayers()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	badges, err := m.DB.GetNBABadges()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	playersBadges, err := m.DB.GetNBAPlayersBadges()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["nba_players"] = players
	data["nba_teams"] = teams
	data["nba_badges"] = badges
	data["nba_players_badges"] = playersBadges

	render.Template(w, r, "nba_players.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) PostNBAPlayers(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	playerID, err := strconv.Atoi(r.Form.Get("player_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	nullInt := true

	teamID, err := strconv.Atoi(r.Form.Get("team_id"))
	if err != nil {
		log.Println(err)
		nullInt = false
		// helpers.ServerError(w, err)
	}

	player := models.NBAPlayer{
		PlayerID: playerID,
		TeamID:   sql.NullInt64{int64(teamID), nullInt},
		Assigned: 0,
	}

	err = m.DB.UpdateNBAPlayer(player)
	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "nba_players", player)

	http.Redirect(w, r, "/nba/players", http.StatusSeeOther)
}

func (m *Repository) NBATeams(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("action") == "drop" {
		playerID, err := strconv.Atoi(r.FormValue("playerID"))
		if err != nil {
			helpers.ServerError(w, err)
		}

		err = m.DB.DropNBAPlayer(playerID)
		if err != nil {
			helpers.ServerError(w, err)
		}
	} else if r.FormValue("action") == "add" {
		playerID, err := strconv.Atoi(r.FormValue("playerID"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		teamID, err := strconv.Atoi(r.FormValue("teamID"))
		if err != nil {
			helpers.ServerError(w, err)
		}

		err = m.DB.AddNBAPlayer(playerID, teamID)
		if err != nil {
			helpers.ServerError(w, err)
		}
	}

	if r.FormValue("player_id") != "" {
		playerID, err := strconv.Atoi(r.FormValue("player_id"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		teamID, err := strconv.Atoi(r.FormValue("team_id"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		assigned, err := strconv.Atoi(r.FormValue("assigned"))
		if err != nil {
			helpers.ServerError(w, err)
		}
		player := models.NBAPlayer{
			PlayerID: playerID,
			TeamID:   sql.NullInt64{int64(teamID), true},
			Assigned: assigned,
		}
		err = m.DB.AssignNBAPlayer(player)
		if err != nil {
			helpers.ServerError(w, err)
		}
	}

	var emptyTeamInfo models.NBATeamInfo
	data := make(map[string]interface{})

	teams, err := m.DB.GetNBATeamInfo()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	players, err := m.DB.GetNBAPlayers()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["teamInfo"] = emptyTeamInfo
	data["nba_players"] = players
	data["nba_teams"] = teams
	data["positions"] = positions

	render.Template(w, r, "nba_teams.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostNBATeams(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	text := r.Form.Get("dark_text")
	darkText := "false"
	if text == "true" {
		darkText = "true"
	}

	teamID, err := strconv.Atoi(r.Form.Get("team_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	teamInfo := models.NBATeamInfo{
		ID:           teamID,
		Name:         r.Form.Get("team_name"),
		Abbreviation: r.Form.Get("abbreviation"),
		Color1:       r.Form.Get("team_color1"),
		Color2:       r.Form.Get("team_color2"),
		DarkText:     darkText,
	}

	form := forms.New(r.PostForm)

	form.Required("team_name", "abbreviation")
	form.MaxLength("abbreviation", 4)

	if !form.Valid() {
		data := make(map[string]interface{})

		teams, err := m.DB.GetNBATeamInfo()
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		players, err := m.DB.GetNBAPlayers()
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		data["nba_players"] = players
		data["nba_teams"] = teams
		data["teamInfo"] = teamInfo
		data["positions"] = positions

		render.Template(w, r, "nba_teams.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.UpdateNBATeamInfo(teamInfo)
	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "team_info", teamInfo)

	http.Redirect(w, r, "/nba/teams", http.StatusSeeOther)
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// Handles request for availability and sends JSON response
func (m *Repository) NBATeamsAvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(out)
}

func (m *Repository) NBATeamInfoSummary(w http.ResponseWriter, r *http.Request) {
	teamInfo, ok := m.App.Session.Get(r.Context(), "team_info").(models.NBATeamInfo)
	if !ok {
		m.App.ErrorLog.Println("Can't get error from session")
		m.App.Session.Put(r.Context(), "error", "Cant get team info from session")
		http.Redirect(w, r, "/nba/teams", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "team_info")
	data := make(map[string]interface{})
	data["team_info"] = teamInfo

	render.Template(w, r, "nba_team_info.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) NBAResults(w http.ResponseWriter, r *http.Request) {
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
			HomeTeam:  homeTeam,
			HomeScore: homeScore,
			AwayScore: awayScore,
			AwayTeam:  awayTeam,
		}
		err = m.DB.AddNBAResult(result)
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
		timestampString := r.FormValue("timestamp")
		layout := "2006-01-02 15:04:05 -0700 MST"
		timestamp, err := time.Parse(layout, timestampString)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		result := models.Result{
			HomeTeam:  homeTeam,
			HomeScore: homeScore,
			AwayScore: awayScore,
			AwayTeam:  awayTeam,
			Time:      timestamp,
		}
		err = m.DB.UpdateNBAResult(result)
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
			Time: timestamp,
		}
		err = m.DB.DeleteNBAResult(result)
		if err != nil {
			helpers.ServerError(w, err)
		}
	}

	var emptyStandings models.Result
	data := make(map[string]interface{})

	teams, err := m.DB.GetNBATeamInfo()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	standings, err := m.DB.GetNBAStandings()
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

	render.Template(w, r, "nba_results.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostNBAResults(w http.ResponseWriter, r *http.Request) {
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
		HomeTeam:  homeTeam,
		HomeScore: homeScore,
		AwayScore: awayScore,
		AwayTeam:  awayTeam,
	}

	form := forms.New(r.PostForm)

	form.Required("home_team", "home_score", "away_score", "away_team")
	form.IsDuplicate("home_team", "away_team", "Home and away have to be different")
	form.IsDuplicate("home_score", "away_score", "Score can't be a draw")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["NBAresult"] = result

		render.Template(w, r, "nba_results.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.AddNBAResult(result)
	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "nba_result", result)

	http.Redirect(w, r, "/nba/results", http.StatusSeeOther)
}

func (m *Repository) NBADraft(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "nba_draft.page.tmpl", &models.TemplateData{})
}
