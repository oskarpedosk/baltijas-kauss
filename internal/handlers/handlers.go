package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/oskarpedosk/baltijas-kauss/internal/config"
	"github.com/oskarpedosk/baltijas-kauss/internal/driver"
	"github.com/oskarpedosk/baltijas-kauss/internal/forms"
	"github.com/oskarpedosk/baltijas-kauss/internal/helpers"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"
	"github.com/oskarpedosk/baltijas-kauss/internal/repository"
	"github.com/oskarpedosk/baltijas-kauss/internal/repository/dbrepo"
	"github.com/oskarpedosk/baltijas-kauss/utilities"
)

// Repo the repository used by the handlers
var Repo *Repository

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
	render.Template(w, r, "nba_players.page.tmpl", &models.TemplateData{})
}

func (m *Repository) NBATeams(w http.ResponseWriter, r *http.Request) {
	var emptyTeamInfo models.NBATeamInfo
	data := make(map[string]interface{})
	data["teamInfo"] = emptyTeamInfo

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

	teamInfo := models.NBATeamInfo{
		TeamName:     r.Form.Get("team_name"),
		Abbreviation: r.Form.Get("abbreviation"),
		TeamColor:    r.Form.Get("team_color"),
		DarkText:     r.Form.Get("text_color"),
	}

	form := forms.New(r.PostForm)

	form.Required("team_name", "abbreviation")
	form.MaxLength("abbreviation", 4)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["teamInfo"] = teamInfo

		render.Template(w, r, "nba_teams.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "team_info", teamInfo)

	http.Redirect(w, r, "/nba/team-info-summary", http.StatusSeeOther)

	//team_name := r.Form.Get("team_name")
	//abbreviation := r.Form.Get("abbreviation")
	//team_color := r.Form.Get("team_color")
	//text_color := r.Form.Get("text_color")
	//w.Write([]byte(fmt.Sprintf("team name is: %s and abbreviation is: %s and team color is: %s and text color is %s", team_name, abbreviation, team_color, text_color)))
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
	var emptyResult models.Result
	data := make(map[string]interface{})
	data["result"] = emptyResult

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

	result := models.Result{
		HomeTeam:  r.Form.Get("home_team"),
		HomeScore: utilities.Atoi(r.Form.Get("home_score")),
		AwayScore: utilities.Atoi(r.Form.Get("away_score")),
		AwayTeam:  r.Form.Get("away_team"),
	}

	form := forms.New(r.PostForm)

	form.Required("home_team", "home_score", "away_score", "away_team")
	form.IsDuplicate("home_team", "away_team", "Home and away have to be different")
	form.IsDuplicate("home_score", "away_score", "Score can't be a draw")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["result"] = result

		render.Template(w, r, "nba_results.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

}
