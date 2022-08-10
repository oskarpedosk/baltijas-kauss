package handlers

import (
	"github.com/oskarpedosk/baltijas-kauss/pkg/config"
	"github.com/oskarpedosk/baltijas-kauss/pkg/models"
	"github.com/oskarpedosk/baltijas-kauss/pkg/render"
	"github.com/oskarpedosk/baltijas-kauss/utilities"
	"net/http"
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

func (m *Repository) NBAHome(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "nbahome.page.tmpl", &models.TemplateData{})
}

func (m *Repository) NBAPlayers(w http.ResponseWriter, r *http.Request) {

	playerData := make(map[string]interface{})
	players := utilities.ReadNBAPlayerData("../../static/jsondata/nba2k_player_data.json")

	for i := 0; i < len(players); i++ {
		playerData[players[i].FirstName + players[i].LastName] = players[i]
	}

	stringMap := make(map[string]string)
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "players.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		NBAPlayerData: playerData,
	})
}

func (m *Repository) SignIn(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "signin.page.tmpl", &models.TemplateData{})
}