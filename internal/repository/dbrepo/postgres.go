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
func (m *postgresDBRepo) UpdateTeamInfo(res models.NBATeamInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update nba_teams set name = $2, abbreviation = $3, team_color = $4, text_color = $5 where team_id = $1`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.ID,
		res.Name,
		res.Abbreviation,
		res.TeamColor,
		res.DarkText,
	)

	if err != nil {
		return err
	}

	return nil
}
