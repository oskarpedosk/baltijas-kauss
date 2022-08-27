package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/oskarpedosk/baltijas-kauss/internal/config"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"
	"github.com/oskarpedosk/baltijas-kauss/utilities"
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

	playerData := make(map[string]interface{})
	players := utilities.ReadNBAPlayerData("../../static/jsondata/nba2k_player_data.json")

	for i := 0; i < len(players); i++ {
		playerData[players[i].FirstName+players[i].LastName] = players[i]
	}

	stringMap := make(map[string]string)
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, r, "nba_players.page.tmpl", &models.TemplateData{
		StringMap:     stringMap,
		NBAPlayerData: playerData,
	})
}

func (m *Repository) NBATeams(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "nba_teams.page.tmpl", &models.TemplateData{})
}

func (m *Repository) PostNBATeams(w http.ResponseWriter, r *http.Request) {
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
		OK: true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(out)
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
