package dbrepo

import (
	"context"
	"time"

	"github.com/oskarpedosk/baltijas-kauss/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// Updates NBA team info
func (m *postgresDBRepo) UpdateNBATeamInfo(res models.NBATeamInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update nba_teams set name = $2, abbreviation = $3, team_color = $4, dark_text = $5 where team_id = $1`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.ID,
		res.Name,
		res.Abbreviation,
		res.Color,
		res.DarkText,
	)

	if err != nil {
		return err
	}

	return nil
}

// Adds a result to NBA results table
func (m *postgresDBRepo) AddNBAResult(res models.Result) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into nba_results (home_team_id, home_score, away_score, away_team_id) 
	values ($1, $2, $3, $4)`

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

// Display all NBA players
func (m *postgresDBRepo) DisplayNBAPlayers() ([]models.NBAPlayer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var players []models.NBAPlayer

	query := `
	select 
		player_id, first_name, last_name, nba_team 
	from 
		nba23_players`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return players, err
	}

	for rows.Next() {
		var player models.NBAPlayer
		err := rows.Scan(
			&player.PlayerID,
			&player.FirstName,
			&player.LastName,
			&player.NBATeam,
		)
		if err != nil {
			return players, err
		}
		players = append(players, player)
	}

	if err = rows.Err(); err != nil {
		return players, err
	}

	return players, nil
}
