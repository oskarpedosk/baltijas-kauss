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

func (m *postgresDBRepo) GetDrafts() ([]models.DraftPick, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var drafts = []models.DraftPick{}

	query := `
		SELECT draft_id, created_at, updated_at
		FROM drafts
		ORDER BY draft_id DESC
		`
		rows, err := m.DB.QueryContext(ctx, query)
		if err != nil {
			return drafts, err
		}
		defer rows.Close()
		for rows.Next() {
			var draft models.DraftPick
			err := rows.Scan(
				&draft.DraftID,
				&draft.CreatedAt,
				&draft.UpdatedAt,
			)
			if err != nil {
				return drafts, err
			}
			drafts = append(drafts, draft)
		}

	return drafts, nil
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
	return draftID, nil
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

func (m *postgresDBRepo) GetDraft(draftID int) ([]models.DraftPick, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var draft []models.DraftPick

	query := `
	SELECT drafts.pick, drafts.team_id, CONCAT(players.first_name, ' ', players.last_name) AS name
	FROM drafts
	JOIN players ON players.player_id = drafts.player_id
	WHERE draft_id = $1;
	`

	rows, err := m.DB.QueryContext(ctx, query, draftID)
	if err != nil {
		return draft, err
	}

	defer rows.Close()
	for rows.Next() {
		var pick models.DraftPick
		err := rows.Scan(
			&pick.Pick,
			&pick.TeamID,
			&pick.Name,
		)
		if err != nil {
			return draft, err
		}
		draft = append(draft, pick)
	}

	return draft, nil
}