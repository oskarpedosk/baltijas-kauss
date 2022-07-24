package handlers

import (
	"2K22/pkg/config"
	"2K22/pkg/models"
	"2K22/pkg/render"
	"2K22/utilities"
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

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Players(w http.ResponseWriter, r *http.Request) {

	players := utilities.ReadJson("../../player_data.json")
	playerData := make(map[string]interface{})
	playerData["players"] = players

	stringMap := make(map[string]string)
	stringMap["test"] = "hello again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "players.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		PlayerData: playerData,
	})
}
