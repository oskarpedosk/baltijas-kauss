package repository

import (
	"context"
	"database/sql"
	"math"
	"time"

	"github.com/oskarpedosk/baltijas-kauss/internal/config"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
)

type StandingsRepo interface {
	GetNBAStandings() ([]models.NBAStandings, error)
	GetLastResults(count int) ([]models.Result, error)
	AddNBAResult(res models.Result) error
	UpdateNBAResult(res models.Result) error
	DeleteNBAResult(res models.Result) error
}

type standingsRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewStandingsRepo(conn *sql.DB, a *config.AppConfig) StandingsRepo {
	return &standingsRepo{
		App: a,
		DB:  conn,
	}
}

// Adds a result to NBA results table
func (m *standingsRepo) AddNBAResult(res models.Result) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into nba_results (home_team_id, home_score, away_score, away_team_id, timestamp) 
	values ($1, $2, $3, $4, CURRENT_TIMESTAMP)`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.HomeTeam,
		res.HomeScore,
		res.AwayScore,
		res.AwayTeam,
	)

	if err != nil {
		return err
	}

	return nil
}

// Update NBA result
func (m *standingsRepo) UpdateNBAResult(res models.Result) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	update 
		nba_results
	set 
		home_team_id = $1,
		home_score = $2,
		away_score = $3,
		away_team_id = $4
	where
		timestamp = $5
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.HomeTeam,
		res.HomeScore,
		res.AwayScore,
		res.AwayTeam,
		res.Time,
	)

	if err != nil {
		return err
	}

	return nil
}

// Delete NBA result from database
func (m *standingsRepo) DeleteNBAResult(res models.Result) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	delete from 
		nba_results
	where
		timestamp = $1
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.Time,
	)

	if err != nil {
		return err
	}

	return nil
}

// Display all NBA standings
func (m *standingsRepo) GetNBAStandings() ([]models.NBAStandings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var standingsSlice []models.NBAStandings

	for i := 1; i < 5; i++ {
		homeWins := 0
		homeLosses := 0
		awayWins := 0
		awayLosses := 0
		basketsFor := 0
		basketsAgainst := 0
		games := ""
		query := `
			select 
				*
			from 
				nba_results
			where
				"home_team_id" = $1 or "away_team_id" = $1
			order by
				timestamp asc
			`

		rows, err := m.DB.QueryContext(ctx, query, i)
		if err != nil {
			return standingsSlice, err
		}

		var singleGame models.Result
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(
				&singleGame.HomeTeam,
				&singleGame.HomeScore,
				&singleGame.AwayScore,
				&singleGame.AwayTeam,
				&singleGame.Time,
			)
			if err != nil {
				return standingsSlice, err
			}

			if singleGame.HomeTeam == i {
				basketsFor += singleGame.HomeScore
				basketsAgainst += singleGame.AwayScore
				if singleGame.HomeScore > singleGame.AwayScore {
					homeWins += 1
					games = "W" + games
				} else {
					homeLosses += 1
					games = "L" + games
				}
			} else {
				basketsFor += singleGame.AwayScore
				basketsAgainst += singleGame.HomeScore
				if singleGame.AwayScore > singleGame.HomeScore {
					awayWins += 1
					games = "W" + games
				} else {
					awayLosses += 1
					games = "L" + games
				}
			}
		}

		x := 5
		chars := []rune(games)

		if len(chars) < x {
			x = len(chars)
		}

		lastGames := []string{"", "", "", "", ""}
		y := 0

		for i := x; i > 0; i-- {
			lastGames[i-1] = string(chars[y])
			y++
		}

		streak := ""
		streakCount := 0

		if games != "" {
			streak = string(chars[0])
			if len(chars) > 1 {
				for i := 0; i < len(chars); i++ {
					if string(chars[i]) != streak {
						break
					} else {
						streakCount += 1
					}
				}
			} else {
				streakCount = 1
			}
		}

		if err = rows.Err(); err != nil {
			return standingsSlice, err
		}

		totalWins := (homeWins + awayWins)
		totalGames := (homeWins + homeLosses + awayWins + awayLosses)
		winPercentage := 0

		forAvg := 0.0
		againstAvg := 0.0

		if totalGames != 0 {
			winPercentage = totalWins * 1000 / totalGames
			forAvg = toFixed(float64(basketsFor)/float64(totalGames), 1)
			againstAvg = toFixed(float64(basketsAgainst)/float64(totalGames), 1)
		}

		teamStandings := models.NBAStandings{
			TeamID:         i,
			WinPercentage:  winPercentage,
			Played:         totalGames,
			TotalWins:      homeWins + awayWins,
			TotalLosses:    homeLosses + awayLosses,
			HomeWins:       homeWins,
			HomeLosses:     homeLosses,
			AwayWins:       awayWins,
			AwayLosses:     awayLosses,
			Streak:         streak,
			StreakCount:    streakCount,
			BasketsFor:     basketsFor,
			BasketsAgainst: basketsAgainst,
			BasketsSum:     basketsFor - basketsAgainst,
			ForAvg:         float64(forAvg),
			AgainstAvg:     float64(againstAvg),
			LastFive:       lastGames,
		}
		standingsSlice = append(standingsSlice, teamStandings)
	}

	orderedSlice := order(standingsSlice)
	return orderedSlice, nil
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func order(slice []models.NBAStandings) []models.NBAStandings {
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

// Display however many last games
func (m *standingsRepo) GetLastResults(count int) ([]models.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var resultsSlice []models.Result

	query := `
			select 
				*
			from 
				nba_results
			order by
				timestamp asc
			`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return resultsSlice, err
	}

	var singleGame models.Result
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&singleGame.HomeTeam,
			&singleGame.HomeScore,
			&singleGame.AwayScore,
			&singleGame.AwayTeam,
			&singleGame.Time,
		)
		if err != nil {
			return resultsSlice, err
		}

		layout := "02/01/2006 15:04"

		result := models.Result{
			HomeTeam:   singleGame.HomeTeam,
			HomeScore:  singleGame.HomeScore,
			AwayScore:  singleGame.AwayScore,
			AwayTeam:   singleGame.AwayTeam,
			Time:       singleGame.Time,
			TimeString: singleGame.Time.Round(15 * time.Minute).Format(layout),
		}

		resultsSlice = append([]models.Result{result}, resultsSlice...)
	}

	if err = rows.Err(); err != nil {
		return resultsSlice, err
	}

	return resultsSlice, nil
}
