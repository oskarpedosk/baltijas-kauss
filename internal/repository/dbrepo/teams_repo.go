package dbrepo

import (
	"context"
	"time"

	"github.com/oskarpedosk/baltijas-kauss/internal/models"
)

func (m *postgresDBRepo) GetTeam(teamID int) (models.Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select 
		*
	from 
		teams
	where
		"team_id" = $1
	`

	var team models.Team

	row, err := m.DB.QueryContext(ctx, query, teamID)
	if err != nil {
		return team, err
	}

	defer row.Close()
	for row.Next() {
		err := row.Scan(
			&team.TeamID,
			&team.Name,
			&team.Abbreviation,
			&team.Color1,
			&team.Color2,
			&team.TextColor,
			&team.UserID,
			&team.CreatedAt,
			&team.UpdatedAt,
		)
		if err != nil {
			return team, err
		}
	}

	return team, nil
}

func (m *postgresDBRepo) GetTeams() ([]models.Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var teams []models.Team

	query := `
	select 
		*
	from 
		teams
	order by
		"team_id" asc
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return teams, err
	}

	defer rows.Close()
	for rows.Next() {
		var team models.Team
		err := rows.Scan(
			&team.TeamID,
			&team.Name,
			&team.Abbreviation,
			&team.Color1,
			&team.Color2,
			&team.TextColor,
			&team.UserID,
			&team.CreatedAt,
			&team.UpdatedAt,
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
func (m *postgresDBRepo) UpdateTeam(team models.Team) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update teams set name = $1, abbreviation = $2, team_color1 = $3, team_color2 = $4, text_color = $5, updated_at = now() where team_id = $6`

	_, err := m.DB.ExecContext(ctx, stmt,
		team.Name,
		team.Abbreviation,
		team.Color1,
		team.Color2,
		team.TextColor,
		team.TeamID,
	)

	if err != nil {
		return err
	}

	return nil
}
