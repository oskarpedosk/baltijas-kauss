package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/oskarpedosk/baltijas-kauss/internal/forms"
	"github.com/oskarpedosk/baltijas-kauss/internal/helpers"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"
)

// Player is the single player handler
func (m *Repository) Player(w http.ResponseWriter, r *http.Request) {
	playerID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	player, err := m.DB.GetPlayer(playerID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	teams, err := m.DB.GetTeams()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var playersTeam models.Team
	for _, team := range teams {
		if team.TeamID == player.TeamID {
			playersTeam = team
			break
		}
	}

	data := make(map[string]interface{})
	data["player"] = player
	data["team"] = playersTeam
	data["teams"] = teams[1:]
	data["FA"] = teams[0]

	render.Template(w, r, "player.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) PostPlayer(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	playerID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	teamID, err := strconv.Atoi(r.FormValue("team_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	err = m.DB.AddPlayer(playerID, teamID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

func (m *Repository) Players(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	fmt.Println(r.URL.String())
	fmt.Println(r.RequestURI)
	fmt.Println(r.URL.RawQuery)

	page := 1
	perPage := 20

	baseURL := r.URL.String() + "&"

	pagination, err := m.DB.GetPaginationData(page, perPage, "players", baseURL)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	teams, err := m.DB.GetTeams()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	players, err := m.DB.FilterPlayers(perPage, pagination.Offset, r.URL.Query())
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	fmt.Println(page)
	// players, err := m.DB.GetPlayers(perPage, pagination.Offset)
	// if err != nil {
	// 	helpers.ServerError(w, err)
	// 	return
	// }

	var playersWithTeamInfo []models.PlayerWithTeamInfo

	for _, player := range players {
		for _, team := range teams {
			if player.TeamID == team.TeamID {
				playerWithTeamInfo := models.PlayerWithTeamInfo{
					Player: player,
					Team:   team,
				}
				playersWithTeamInfo = append(playersWithTeamInfo, playerWithTeamInfo)
				break
			}
		}
	}

	ranking := []int{}
	for i := 1; i <= len(players); i++ {
		ranking = append(ranking, i+pagination.Offset)
	}

	data := make(map[string]interface{})
	data["players"] = playersWithTeamInfo
	data["teams"] = teams[1:]
	data["FA"] = teams[0]
	data["ranking"] = ranking
	data["pagination"] = pagination

	render.Template(w, r, "players.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostPlayers(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	switch r.FormValue("action") {
	case "change_team":
		playerID, err := strconv.Atoi(r.FormValue("player_id"))
		if err != nil {
			helpers.ServerError(w, err)
		}

		teamID, err := strconv.Atoi(r.FormValue("team_id"))
		if err != nil {
			helpers.ServerError(w, err)
		}

		player := models.Player{
			PlayerID:         playerID,
			TeamID:           teamID,
			AssignedPosition: 0,
		}

		err = m.DB.SwitchTeam(player)
		if err != nil {
			helpers.ServerError(w, err)
		}
	case "filter":
		v := url.Values{}
		v.Add("ovrh", r.FormValue("ovrh"))
		v.Add("ovrl", r.FormValue("ovrl"))
		v.Add("hh", r.FormValue("hh"))
		v.Add("hl", r.FormValue("hl"))

		redirectURL := r.URL.Path + "?" + v.Encode()

		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

func (m *Repository) PostUpdatePlayer(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	playerID := r.FormValue("player_id")
	ratingsURL := r.FormValue("ratings_url")

	go func(playerID, ratingsURL string) {
		filePath := "./static/js/script/updateplayer.js"
		cmd := exec.Command("node", filePath, playerID, ratingsURL)
		output, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		}

		var player models.Player
		json.Unmarshal(output, &player)

		err = m.DB.UpdatePlayer(player)
		if err != nil {
			helpers.ServerError(w, err)
		}
	}(playerID, ratingsURL)

	http.Redirect(w, r, "/players", http.StatusSeeOther)
}
