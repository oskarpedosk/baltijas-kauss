package main

import (
	"2K22/utilities"
	"fmt"
	"net/http"
)

func main() {
	scrapedData := utilities.ScrapeDataFromURL("https://www.2kratings.com/lists/top-100-highest-nba-2k-ratings")
	utilities.WriteToJson("player_data.json", scrapedData)
	players := utilities.ReadJson("player_data.json")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		for i := 0; i < len(players); i++ {
			// Print player info
			fmt.Fprint(w, "#", players[i].Rank)
			fmt.Fprint(w, " ", players[i].Name)
			fmt.Fprint(w, " ", players[i].Team)
			fmt.Fprintln(w, "")

			// Print player positions
			for j := 0; j < len(players[i].Positions); j++ {
				fmt.Fprint(w, players[i].Positions[j])
				if j != len(players[i].Positions)-1 {
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
			if i != len(players)-1 {
				fmt.Fprintln(w, "")
			}
		}
	})
	_ = http.ListenAndServe(":8080", nil)
}
