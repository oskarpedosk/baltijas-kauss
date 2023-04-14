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
		var err error
		seasonID, err = strconv.Atoi(r.URL.Query().Get("s"))
		if err != nil {
			helpers.ServerError(w, err)
		}
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
		for j := 0; j < len(slice)-i-1; j++ {
			if slice[j].WinPercentage < slice[j+1].WinPercentage {
				slice[j], slice[j+1] = slice[j+1], slice[j]
			} else if slice[j].WinPercentage == slice[j+1].WinPercentage {
				if slice[j].BasketsSum < slice[j+1].BasketsSum {
					slice[j], slice[j+1] = slice[j+1], slice[j]
				}
			}
		}
	}
	return slice
}

func (m *Repository) PostStandings(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if r.FormValue("home_team") == "" || r.FormValue("away_team") == "" {
		m.App.Session.Put(r.Context(), "warning", "Please choose teams")
		http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
		return
	}
	if r.FormValue("home_score") == "" || r.FormValue("away_score") == "" {
		m.App.Session.Put(r.Context(), "warning", "Please insert score")
		http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
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
		http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
		return
	}

	err = m.DB.AddResult(result)
	if err != nil {
		msg := fmt.Sprintf("error is %v", err)
		m.App.Session.Put(r.Context(), "error", msg)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "Result successfully added!")
	http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
}

func (m *Repository) AllTime(w http.ResponseWriter, r *http.Request) {
	results, err := m.DB.GetAllResults()
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

	var allTeamsHeadToHead = [][]models.HeadToHead{}
	for i := 0; i < len(teamsWithoutFA); i++ {
		team1 := teamsWithoutFA[i]
		var team1HeadToHead = []models.HeadToHead{}
		for j := 0; j < len(teamsWithoutFA); j++ {
			team2 := teamsWithoutFA[j]
			if team1.TeamID != team2.TeamID {
				headToHeadResults, err := m.DB.GetHeadToHeadResults(team1.TeamID, team2.TeamID)
				if err != nil {
					helpers.ServerError(w, err)
				}
				team1HeadToHead = append(team1HeadToHead, models.HeadToHead{
					Home:      team1,
					Away:      team2,
					Standings: CalculateHeadToHead(team1, team2, headToHeadResults),
				})
			}
		}
		allTeamsHeadToHead = append(allTeamsHeadToHead, team1HeadToHead)
	}

	data := make(map[string]interface{})
	data["teams"] = teamsWithoutFA
	data["standings"] = standings
	data["allTeamsHeadToHead"] = allTeamsHeadToHead

	render.Template(w, r, "alltime.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func CalculateHeadToHead(team1, team2 models.Team, results []models.Result) models.Standings {
	headToHead := models.Standings{}

	// Count wins and losses
	for _, result := range results {
		if result.HomeScore > result.AwayScore {
			if result.HomeTeam.TeamID == team1.TeamID {
				headToHead.HomeWins++
				headToHead.Streak += "W"
				headToHead.BasketsFor += result.HomeScore
				headToHead.BasketsAgainst += result.AwayScore
			} else {
				headToHead.AwayLosses++
				headToHead.Streak += "L"
				headToHead.BasketsFor += result.AwayScore
				headToHead.BasketsAgainst += result.HomeScore
			}
		} else {
			if result.HomeTeam.TeamID == team1.TeamID {
				headToHead.HomeLosses++
				headToHead.Streak += "L"
				headToHead.BasketsFor += result.HomeScore
				headToHead.BasketsAgainst += result.AwayScore
			} else {
				headToHead.AwayWins++
				headToHead.Streak += "W"
				headToHead.BasketsFor += result.AwayScore
				headToHead.BasketsAgainst += result.HomeScore
			}

		}
	}

	headToHead.TotalWins = headToHead.HomeWins + headToHead.AwayWins
	headToHead.TotalLosses = headToHead.HomeLosses + headToHead.AwayLosses
	headToHead.TotalLosses = headToHead.HomeLosses + headToHead.AwayLosses
	headToHead.Played = headToHead.TotalWins + headToHead.TotalLosses
	winsLosses := headToHead.Streak
	if len(headToHead.Streak) > 0 {
		headToHead.Streak = string(headToHead.Streak[0])
		streakCount := 0
		for _, char := range winsLosses {
			if string(char) == string(headToHead.Streak[0]) {
				streakCount++
			} else {
				break
			}
		}
		headToHead.StreakCount = streakCount
	} else {
		headToHead.Streak = ""
		headToHead.StreakCount = 0
	}
	headToHead.BasketsSum = headToHead.BasketsFor - headToHead.BasketsAgainst

	winPercentage := 0
	forAvg := 0.0
	againstAvg := 0.0
	if headToHead.Played != 0 {
		winPercentage = headToHead.TotalWins * 1000 / headToHead.Played
		forAvg = toFixed(float64(headToHead.BasketsFor)/float64(headToHead.Played), 1)
		againstAvg = toFixed(float64(headToHead.BasketsAgainst)/float64(headToHead.Played), 1)
	}
	headToHead.WinPercentage = winPercentage
	headToHead.ForAvg = forAvg
	headToHead.AgainstAvg = againstAvg

	lastGames := []string{"", "", "", "", "", "", "", "", "", ""}
	x := 10
	y := 0

	if len(winsLosses) < x {
		x = len(winsLosses)
	}

	for i := x; i > 0; i-- {
		lastGames[i-1] = string(winsLosses[y])
		y++
	}

	headToHead.LastFive = lastGames

	return headToHead
}
