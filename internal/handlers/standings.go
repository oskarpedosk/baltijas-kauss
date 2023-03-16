package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/oskarpedosk/baltijas-kauss/internal/forms"
	"github.com/oskarpedosk/baltijas-kauss/internal/helpers"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"
)

func (m *Repository) Standings(w http.ResponseWriter, r *http.Request) {
	seasonID := 0
	if r.URL.Query().Has("s") {
		seasonID, _ = strconv.Atoi(r.URL.Query().Get("s"))
	}

	seasons, err := m.DB.GetSeasons()
	if err != nil {
		helpers.ServerError(w, err)
	}

	results, err := m.DB.GetSeasonResults(seasonID)
	if err != nil {
		helpers.ServerError(w, err)
	}
	teams, err := m.DB.GetTeams()
	if err != nil {
		helpers.ServerError(w, err)
	}
	var teamsWithoutFA = []models.Team{}
	for _, team := range teams {
		if team.TeamID != 1 {
			teamsWithoutFA = append(teamsWithoutFA, team)
		}
	}
	standings := CalculateStandings(teamsWithoutFA, results)

	activeSeason := seasonID
	if seasonID == 0 {
		activeSeason = seasons[0].SeasonID
	}

	data := make(map[string]interface{})
	data["results"] = results
	data["seasons"] = seasons
	data["teams"] = teamsWithoutFA
	data["standings"] = standings
	data["activeSeason"] = activeSeason

	render.Template(w, r, "standings.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func CalculateStandings(teams []models.Team, results []models.Result) []models.Standings {
	data := make(map[int]models.Standings)
	for _, team := range teams {
		data[team.TeamID] = models.Standings{Team: team}
	}

	// Count wins and losses
	for _, result := range results {
		if result.HomeScore > result.AwayScore {
			standings := data[result.HomeTeam.TeamID]
			standings.HomeWins++
			standings.Streak += "W"
			standings.BasketsFor += result.HomeScore
			standings.BasketsAgainst += result.AwayScore
			data[result.HomeTeam.TeamID] = standings

			standings2 := data[result.AwayTeam.TeamID]
			standings2.AwayLosses++
			standings2.Streak += "L"
			standings2.BasketsFor += result.AwayScore
			standings2.BasketsAgainst += result.HomeScore
			data[result.AwayTeam.TeamID] = standings2
		} else {
			standings := data[result.AwayTeam.TeamID]
			standings.AwayWins++
			standings.Streak += "W"
			standings.BasketsFor += result.AwayScore
			standings.BasketsAgainst += result.HomeScore
			data[result.AwayTeam.TeamID] = standings

			standings2 := data[result.HomeTeam.TeamID]
			standings2.HomeLosses++
			standings2.Streak += "L"
			standings2.BasketsFor += result.HomeScore
			standings2.BasketsAgainst += result.AwayScore
			data[result.HomeTeam.TeamID] = standings2
		}
	}

	var standings = []models.Standings{}

	for _, v := range data {
		teamStats := v
		teamStats.TotalWins = v.HomeWins + v.AwayWins
		teamStats.TotalLosses = v.HomeLosses + v.AwayLosses
		teamStats.TotalLosses = teamStats.HomeLosses + teamStats.AwayLosses
		teamStats.Played = teamStats.TotalWins + teamStats.TotalLosses
		if len(v.Streak) > 0 {
			teamStats.Streak = string(v.Streak[0])
			streakCount := 0
			for _, char := range v.Streak {
				if string(char) == string(v.Streak[0]) {
					streakCount++
				} else {
					break
				}
			}
			teamStats.StreakCount = streakCount
		} else {
			teamStats.Streak = ""
			teamStats.StreakCount = 0
		}
		teamStats.BasketsSum = v.BasketsFor - v.BasketsAgainst

		winPercentage := 0
		forAvg := 0.0
		againstAvg := 0.0
		if teamStats.Played != 0 {
			winPercentage = teamStats.TotalWins * 1000 / teamStats.Played
			forAvg = toFixed(float64(v.BasketsFor)/float64(teamStats.Played), 1)
			againstAvg = toFixed(float64(v.BasketsAgainst)/float64(teamStats.Played), 1)
		}
		teamStats.WinPercentage = winPercentage
		teamStats.ForAvg = forAvg
		teamStats.AgainstAvg = againstAvg

		lastGames := []string{"", "", "", "", ""}
		x := 5
		y := 0

		if len(v.Streak) < x {
			x = len(v.Streak)
		}

		for i := x; i > 0; i-- {
			lastGames[i-1] = string(v.Streak[y])
			y++
		}

		teamStats.LastFive = lastGames
		standings = append(standings, teamStats)
	}

	orderedStandings := order(standings)

	return orderedStandings
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func order(slice []models.Standings) []models.Standings {
	for i := 0; i < len(slice)-1; i++ {
		if slice[i].WinPercentage < slice[i+1].WinPercentage {
			slice[i], slice[i+1] = slice[i+1], slice[i]
			order(slice)
		}
		if slice[i].WinPercentage == slice[i+1].WinPercentage {
			if slice[i].BasketsSum < slice[i+1].BasketsSum {
				slice[i], slice[i+1] = slice[i+1], slice[i]
			}
		}

	}
	return slice
}

// func (m *Repository) Standings(w http.ResponseWriter, r *http.Request) {
// 	m.DB.CreateNewSeason()
// 	if r.FormValue("action") == "add" {
// 		homeTeam, err := strconv.Atoi(r.FormValue("home_team_id"))
// 		if err != nil {
// 			helpers.ServerError(w, err)
// 		}
// 		homeScore, err := strconv.Atoi(r.FormValue("home_score"))
// 		if err != nil {
// 			helpers.ServerError(w, err)
// 		}
// 		awayScore, err := strconv.Atoi(r.FormValue("away_score"))
// 		if err != nil {
// 			helpers.ServerError(w, err)
// 		}
// 		awayTeam, err := strconv.Atoi(r.FormValue("away_team_id"))
// 		if err != nil {
// 			helpers.ServerError(w, err)
// 		}
// 		result := models.Result{
// 			HomeTeamID: homeTeam,
// 			HomeScore:  homeScore,
// 			AwayScore:  awayScore,
// 			AwayTeamID: awayTeam,
// 		}
// 		err = m.DB.AddResult(result)
// 		if err != nil {
// 			helpers.ServerError(w, err)
// 		}
// 	} else if r.FormValue("action") == "update" {
// 		homeTeam, err := strconv.Atoi(r.FormValue("home_team_id"))
// 		if err != nil {
// 			helpers.ServerError(w, err)
// 		}
// 		homeScore, err := strconv.Atoi(r.FormValue("home_score"))
// 		if err != nil {
// 			helpers.ServerError(w, err)
// 		}
// 		awayScore, err := strconv.Atoi(r.FormValue("away_score"))
// 		if err != nil {
// 			helpers.ServerError(w, err)
// 		}
// 		awayTeam, err := strconv.Atoi(r.FormValue("away_team_id"))
// 		if err != nil {
// 			helpers.ServerError(w, err)
// 		}

// 		if err != nil {
// 			helpers.ServerError(w, err)
// 			return
// 		}

// 		result := models.Result{
// 			HomeTeamID: homeTeam,
// 			HomeScore:  homeScore,
// 			AwayScore:  awayScore,
// 			AwayTeamID: awayTeam,
// 		}
// 		err = m.DB.UpdateResult(result)
// 		if err != nil {
// 			helpers.ServerError(w, err)
// 		}

// 	} else if r.FormValue("action") == "delete" {
// 		timestampString := r.FormValue("timestamp")
// 		layout := "2006-01-02 15:04:05 -0700 MST"
// 		timestamp, err := time.Parse(layout, timestampString)
// 		if err != nil {
// 			helpers.ServerError(w, err)
// 			return
// 		}
// 		result := models.Result{
// 			CreatedAt: timestamp,
// 		}
// 		err = m.DB.DeleteResult(result)
// 		if err != nil {
// 			helpers.ServerError(w, err)
// 		}
// 	}

// 	var emptyStandings models.Result
// 	data := make(map[string]interface{})

// 	teams, err := m.DB.GetTeams()
// 	if err != nil {
// 		helpers.ServerError(w, err)
// 		return
// 	}
// 	standings, err := m.DB.GetStandings()
// 	if err != nil {
// 		helpers.ServerError(w, err)
// 		return
// 	}
// 	lastResults, err := m.DB.GetLastResults(10)
// 	if err != nil {
// 		helpers.ServerError(w, err)
// 		return
// 	}

// 	data["result"] = emptyStandings
// 	data["teams"] = teams
// 	data["standings"] = standings
// 	data["last_results"] = lastResults

// 	render.Template(w, r, "standings.page.tmpl", &models.TemplateData{
// 		Form: forms.New(nil),
// 		Data: data,
// 	})
// }

func (m *Repository) PostStandings(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	path := r.URL.Path
	fmt.Println(r.URL.RawPath)
	queries := r.URL.RawQuery
	fmt.Println(queries)
	fullPath := path
	if queries != "" {
		fullPath += "?" + queries
	}
	fmt.Println(fullPath)

	if r.FormValue("home_team") == "" || r.FormValue("away_team") == "" {
		m.App.Session.Put(r.Context(), "error", "Please choose teams")
		http.Redirect(w, r, fullPath, http.StatusSeeOther)
		return
	}
	if r.FormValue("home_score") == "" || r.FormValue("away_score") == "" {
		m.App.Session.Put(r.Context(), "error", "Please insert score")
		http.Redirect(w, r, fullPath, http.StatusSeeOther)
		return
	}

	homeTeamID, err := strconv.Atoi(r.FormValue("home_team"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	homeScore, err := strconv.Atoi(r.Form.Get("home_score"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	awayScore, err := strconv.Atoi(r.Form.Get("away_score"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	awayTeamID, err := strconv.Atoi(r.FormValue("away_team"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	seasonID, err := strconv.Atoi(r.FormValue("season_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	result := models.Result{
		HomeTeam:  models.Team{TeamID: homeTeamID},
		HomeScore: homeScore,
		AwayScore: awayScore,
		AwayTeam:  models.Team{TeamID: awayTeamID},
		SeasonID:  seasonID,
	}

	form := forms.New(r.PostForm)
	form.Required("home_team", "home_score", "away_score", "away_team", "season_id")
	form.AreDifferent("home_team", "away_team", "Teams have to be different")
	form.AreDifferent("home_score", "away_score", "Score can't be a draw")

	if !form.Valid() {
		m.App.Session.Put(r.Context(), "error", "Error adding the result")
		http.Redirect(w, r, fullPath, http.StatusSeeOther)
		// data := make(map[string]interface{})
		// data["NBAresult"] = result

		// render.Template(w, r, "standings.page.tmpl", &models.TemplateData{
		// 	Form: form,
		// 	Data: data,
		// })
		return
	}

	err = m.DB.AddResult(result)
	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "result", result)

	fmt.Println(fullPath)
	http.Redirect(w, r, fullPath, http.StatusSeeOther)
}

func (m *Repository) AllTime(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
