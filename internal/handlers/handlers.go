package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/oskarpedosk/baltijas-kauss/internal/config"
	"github.com/oskarpedosk/baltijas-kauss/internal/forms"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) SignIn(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "signin.page.tmpl", &models.TemplateData{})
}

func (m *Repository) NBAHome(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, r, "nba_home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) NBAPlayers(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, r, "nba_players.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *Repository) NBATeams(w http.ResponseWriter, r *http.Request) {
	var emptyTeamInfo models.TeamInfo
	data := make(map[string]interface{})
	data["teamInfo"] = emptyTeamInfo

	render.RenderTemplate(w, r, "nba_teams.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostNBATeams(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	
	teamInfo := models.TeamInfo{
		TeamName:     r.Form.Get("team_name"),
		Abbreviation: r.Form.Get("abbreviation"),
		TeamColor:    r.Form.Get("team_color"),
		DarkText:     r.Form.Get("text_color"),
	}

	form := forms.New(r.PostForm)

	form.Required("team_name", "abbreviation")
	form.MaxLength("abbreviation", 4, r)


	if !form.Valid() {
		data := make(map[string]interface{})
		data["teamInfo"] = teamInfo

		render.RenderTemplate(w, r, "nba_teams.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "team_info", teamInfo)

	http.Redirect(w, r, "/nba/team-info-summary", http.StatusSeeOther)

	team_name := r.Form.Get("team_name")
	abbreviation := r.Form.Get("abbreviation")
	team_color := r.Form.Get("team_color")
	text_color := r.Form.Get("text_color")
	w.Write([]byte(fmt.Sprintf("team name is: %s and abbreviation is: %s and team color is: %s and text color is %s", team_name, abbreviation, team_color, text_color)))
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
		log.Println(err)
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(out)
}

func (m *Repository) NBATeamInfoSummary(w http.ResponseWriter, r *http.Request) {
	teamInfo, ok := m.App.Session.Get(r.Context(), "team_info").(models.TeamInfo)
	if !ok {
		log.Println("Cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Cant get team info from session")
		http.Redirect(w, r, "/nba/teams", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "team_info")
	data := make(map[string]interface{})
	data["team_info"] = teamInfo

	render.RenderTemplate(w, r, "nba_team_info.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) NBAResults(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "nba_results.page.tmpl", &models.TemplateData{})
}

func (m *Repository) PostNBAResults(w http.ResponseWriter, r *http.Request) {
	home_team := r.Form.Get("home_team")
	home_score := r.Form.Get("home_score")
	away_team := r.Form.Get("away_team")
	away_score := r.Form.Get("away_score")
	w.Write([]byte(fmt.Sprintf("home is: %s and away is: %s and home score is: %s and away score is %s", home_team, away_team, home_score, away_score)))
}
