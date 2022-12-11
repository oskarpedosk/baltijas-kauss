package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/oskarpedosk/baltijas-kauss/internal/config"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
)

type TeamsRepo interface {
	GetNBATeamInfo() ([]models.NBATeam, error)
	UpdateNBATeamInfo(team models.NBATeamInfo) error
	AddNBAPlayer(playerID, teamID int) error
	DropNBAPlayer(playerID int) error
	DropAllNBAPlayers() error
	AssignNBAPlayer(player models.NBAPlayer) error
}



type teamsRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewTeamsRepo(conn *sql.DB, a *config.AppConfig) TeamsRepo {
	return &teamsRepo{
		App: a,
		DB:  conn,
	}
}

// Add NBA player directly to a team
func (m *teamsRepo) AddNBAPlayer(playerID, teamID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Remove team and position
	stmt := `
	update 
		nba_players
	set 
		team_id = $1,
		assigned = 0
	where
		player_id = $2
	`

	_, err := m.DB.ExecContext(ctx, stmt, teamID, playerID)

	if err != nil {
		return err
	}

	return nil
}

var playerCount = 150

// Drop NBA player from a team
func (m *teamsRepo) DropNBAPlayer(playerID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Remove team and position
	stmt := `
	update 
		nba_players
	set 
		team_id = null,
		assigned = 0
	where
		player_id = $1
	`

	_, err := m.DB.ExecContext(ctx, stmt, playerID)

	if err != nil {
		return err
	}

	return nil
}

// Drop all NBA player from a team
func (m *teamsRepo) DropAllNBAPlayers() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	update 
		nba_players
	set 
		team_id = null,
		assigned = 0
	`

	_, err := m.DB.ExecContext(ctx, stmt)

	if err != nil {
		return err
	}

	return nil
}

// Get team info
func (m *teamsRepo) GetNBATeamInfo() ([]models.NBATeam, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var teams []models.NBATeam

	query := `
	select 
		*
	from 
		nba_teams
	order by
		"team_id" asc
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return teams, err
	}

	defer rows.Close()
	for rows.Next() {
		var team models.NBATeam
		err := rows.Scan(
			&team.TeamID,
			&team.Name,
			&team.Abbreviation,
			&team.Color1,
			&team.Color2,
			&team.DarkText,
			&team.OwnerID,
		)
		if err != nil {
			return teams, err
		}
		teams = append(teams, team)
	}

	if err = rows.Err(); err != nil {
		return teams, err
	}

	return teams, nil
}

// Updates NBA team info
func (m *teamsRepo) UpdateNBATeamInfo(team models.NBATeamInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update nba_teams set name = $2, abbreviation = $3, team_color1 = $4, team_color2 = $5, dark_text = $6 where team_id = $1`

	_, err := m.DB.ExecContext(ctx, stmt,
		team.ID,
		team.Name,
		team.Abbreviation,
		team.Color1,
		team.Color2,
		team.DarkText,
	)

	if err != nil {
		return err
	}

	return nil
}


// Assigns NBA player to a position
func (m *teamsRepo) AssignNBAPlayer(player models.NBAPlayer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Remove previous player from position
	stmt := `
	update 
		nba_players
	set 
		assigned = 0
	where
		team_id = $1
	and
		assigned = $2
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		player.TeamID,
		player.Assigned,
	)

	if err != nil {
		return err
	}

	// Add new player to position
	stmt = `
	update 
		nba_players
	set 
		assigned = $3
	where
		player_id = $1
	and
		team_id = $2
	`

	_, err = m.DB.ExecContext(ctx, stmt,
		player.PlayerID,
		player.TeamID,
		player.Assigned,
	)

	if err != nil {
		return err
	}

	return nil
}