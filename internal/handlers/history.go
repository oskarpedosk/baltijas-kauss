package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/oskarpedosk/baltijas-kauss/internal/forms"
	"github.com/oskarpedosk/baltijas-kauss/internal/helpers"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"
)

func (m *Repository) History(w http.ResponseWriter, r *http.Request) {
	draftID := 0
	if r.URL.Query().Has("draft") {
		var err error
		draftID, err = strconv.Atoi(r.URL.Query().Get("draft"))
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
	}

	var draft []models.DraftPick
	if draftID != 0 {
		var err error
		draft, err = m.DB.GetDraft(draftID)
		if err != nil {
			helpers.ServerError(w, err)
		}
	}

	drafts, err := m.DB.GetDrafts()
	if err != nil {
		helpers.ServerError(w, err)
	}

	teams, err := m.DB.GetTeams()
	if err != nil {
		helpers.ServerError(w, err)
	}

	var draftOrder = []models.Team{}
	counter := 0
	for _, pick := range draft {
		for _, team := range teams {
			if pick.TeamID == team.TeamID {
				draftOrder = append(draftOrder, team)
				counter++
				break
			}
		}
		if counter == len(teams) - 1 {
			break
		}
	}

	fmt.Println(draftOrder)

	data := make(map[string]interface{})
	data["drafts"] = drafts
	data["teams"] = draftOrder
	data["draft"] = draftPicks
	data["activeDraft"] = draftID

	render.Template(w, r, "history.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}
