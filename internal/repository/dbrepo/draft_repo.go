package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/oskarpedosk/baltijas-kauss/internal/models"
)

// Add a random player to a team
func (m *postgresDBRepo) SelectRandomPlayer(random int) (models.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Remove team and position
	stmt := `
	select
		"player_id",
		"first_name",
		"last_name",
		"primary_position",
		"secondary_position"
	from 
		players
	where
		"team_id" = 1
	order by
		"overall" desc,
		"attributes/TotalAttributes" desc
	limit
		1
	offset
		$1
	`

	row := m.DB.QueryRowContext(ctx, stmt, random)

	var player models.Player
	err := row.Scan(
		&player.PlayerID,
		&player.FirstName,
		&player.LastName,
		&player.PrimaryPosition,
		&player.SecondaryPosition,
	)

	if err != nil {
		return player, err
	}

	return player, nil
}

func (m *postgresDBRepo) GetDraftID() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT COALESCE(MAX(draft_id), 0) FROM drafts`

	var draftID int
	err := m.DB.QueryRowContext(ctx, stmt).Scan(&draftID)
	if err != nil {
		fmt.Println(err)
		return draftID, err
	}
	return draftID + 1, nil
}

func (m *postgresDBRepo) AddDraftPick(draftID int, draftPick models.DraftPick) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into drafts (draft_id, pick, player_id, team_id, created_at, updated_at) 
	values ($1, $2, $3, $4, now(), now())`

	_, err := m.DB.ExecContext(ctx, stmt,
		draftID,
		draftPick.Pick,
		draftPick.PlayerID,
		draftPick.TeamID,
	)

	if err != nil {
		return err
	}

	return nil
}