package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/oskarpedosk/baltijas-kauss/internal/helpers"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"
)

func (m *Repository) History(w http.ResponseWriter, r *http.Request) {
	draftID := 0
	drafts, err := m.DB.GetDrafts()
	if err != nil {
		helpers.ServerError(w, err)
	}

	if r.URL.Query().Has("draft") {
		draftID, err = strconv.Atoi(r.URL.Query().Get("draft"))
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
	} else {
		if len(drafts) > 0 {
			draftID = drafts[0].DraftID
		}
	}

	if draftID == 0 {
		render.Template(w, r, "history.page.tmpl", &models.TemplateData{})
		return
	}

	var draft []models.DraftPick
	draft, err = m.DB.GetDraft(draftID)
	if err != nil {
		helpers.ServerError(w, err)
	}

	teams, err := m.DB.GetTeams()
	if err != nil {
		helpers.ServerError(w, err)
	}

	teamMap := make(map[int]models.Team)
	for _, team := range teams[1:] {
		team.Name = strings.ToUpper(team.Name)
		teamMap[team.TeamID] = team
	}

	round := 1
	prevTeam := 0
	var row = []models.DraftPick{}
	var draftOrder = []models.Team{}
	var ordererDraft = [][]models.DraftPick{}

	for _, pick := range draft {
		if prevTeam == pick.TeamID {
			if round%2 == 0 {
				reverse(row)
			}
			ordererDraft = append(ordererDraft, row)
			row = []models.DraftPick{}
			round++
		}
		row = append(row, pick)
		if round == 1 {
			draftOrder = append(draftOrder, teamMap[pick.TeamID])
		}
		prevTeam = pick.TeamID
	}
	if round%2 == 0 {
		reverse(row)
	}
	ordererDraft = append(ordererDraft, row)

	data := make(map[string]interface{})
	data["drafts"] = drafts
	data["teams"] = draftOrder
	data["draft"] = ordererDraft
	data["activeDraft"] = draftID

	render.Template(w, r, "history.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func reverse(s []models.DraftPick) {
	for i := 0; i < len(s)/2; i++ {
		j := len(s) - i - 1
		s[i], s[j] = s[j], s[i]
	}
}
