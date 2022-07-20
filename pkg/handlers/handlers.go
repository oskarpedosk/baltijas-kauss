package handlers

import (
	"2K22/pkg/render"
	"2K22/utilities"
	"fmt"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home.page.tmpl")
}

func Players(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "players.page.tmpl")
	players := utilities.ReadJson("../../player_data.json")

	for i := 0; i < len(players); i++ {
		// Print player info
		fmt.Fprint(w, "#", players[i].Rank)
		fmt.Fprint(w, " ", players[i].Name)
		fmt.Fprint(w, " ", players[i].Team)
		fmt.Fprintln(w, "")

		// Print player positions
		for j := 0; j < len(players[i].Positions); j++ {
			fmt.Fprint(w, players[i].Positions[j])
			if j < len(players[i].Positions)-1 {
				fmt.Fprint(w, "/")
			}
		}
		// Print player height
		fmt.Fprintln(w, "")
		fmt.Fprint(w, players[i].Height[0], "'", players[i].Height[1], "\"")

		// Print player ratings
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, players[i].OverallRating, "OVR")
		fmt.Fprintln(w, players[i].ThreePointRating, "3PT")
		fmt.Fprintln(w, players[i].DunkRating, "DUNK")
		if i < len(players)-1 {
			fmt.Fprintln(w, "")
		}
	}
}
