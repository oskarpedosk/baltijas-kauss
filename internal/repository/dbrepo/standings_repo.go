package dbrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/oskarpedosk/baltijas-kauss/internal/models"
)

func (m *postgresDBRepo) AddResult(result models.Result) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT MAX(result_id) FROM results`

	var resultID sql.NullInt64
	err := m.DB.QueryRowContext(ctx, stmt).Scan(&resultID)
	if err != nil {
		// Handle the error
		if err == sql.ErrNoRows || !resultID.Valid {
			resultID.Int64 = 1
			resultID.Valid = true
		} else {
			return err
		}
	}

	stmt = `insert into results (result_id, season_id, home_team_id, home_score, away_score, away_team_id, created_at, updated_at) 
	values ($1, $2, $3, $4, $5, $6, now(), now())`

	_, err = m.DB.ExecContext(ctx, stmt,
		resultID.Int64+1,
		result.SeasonID,
		result.HomeTeam.TeamID,
		result.HomeScore,
		result.AwayScore,
		result.AwayTeam.TeamID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) UpdateResult(result models.Result) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	update 
		results
	set 
		home_team_id = $1,
		home_score = $2,
		away_score = $3,
		away_team_id = $4
	where
		result_id = $5
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		result.HomeTeam.TeamID,
		result.HomeScore,
		result.AwayScore,
		result.AwayTeam.TeamID,
		result.ResultID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) DeleteResult(res models.Result) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	delete from 
		results
	where
		result_id = $1
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.ResultID,
	)

	if err != nil {
		return err
	}

	return nil
}

// Create a new season
func (m *postgresDBRepo) StartNewSeason() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT MAX(season_id) FROM seasons`

	var seasonID int
	err := m.DB.QueryRowContext(ctx, stmt).Scan(&seasonID)
	if err != nil {
		return err
	}

	stmt = `
	INSERT INTO seasons (season_id, created_at, updated_at) 
	VALUES ($1, NOW(), NOW())
	`

	_, err = m.DB.ExecContext(ctx, stmt, seasonID+1)
	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) GetSeasons() ([]models.Season, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var seasons = []models.Season{}

	query := `
		SELECT *
		FROM seasons
		ORDER BY season_id DESC
		`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return seasons, err
	}
	defer rows.Close()
	for rows.Next() {
		var season models.Season
		err := rows.Scan(
			&season.SeasonID,
			&season.CreatedAt,
			&season.UpdatedAt,
		)
		if err != nil {
			return seasons, err
		}
		seasons = append(seasons, season)
	}

	return seasons, nil
}

func (m *postgresDBRepo) GetSeasonResults(seasonID int) ([]models.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var results []models.Result

	if seasonID == 0 {
		query := `
			SELECT MAX(season_id)
			FROM seasons
		`

		err := m.DB.QueryRowContext(ctx, query).Scan(&seasonID)
		if err != nil {
			return results, err
		}
	}
	query := `
		SELECT *
		FROM results
		WHERE season_id = $1
		ORDER BY result_id DESC
		`
	rows, err := m.DB.QueryContext(ctx, query, seasonID)
	if err != nil {
		return results, err
	}
	defer rows.Close()
	for rows.Next() {
		var result models.Result
		err := rows.Scan(
			&result.ResultID,
			&result.SeasonID,
			&result.HomeTeam.TeamID,
			&result.HomeScore,
			&result.AwayScore,
			&result.AwayTeam.TeamID,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return results, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (m *postgresDBRepo) GetAllResults() ([]models.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var results []models.Result
	query := `
		SELECT *
		FROM results
		ORDER BY result_id DESC
		`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return results, err
	}

	defer rows.Close()
	for rows.Next() {
		var result models.Result
		err := rows.Scan(
			&result.ResultID,
			&result.SeasonID,
			&result.HomeTeam.TeamID,
			&result.HomeScore,
			&result.AwayScore,
			&result.AwayTeam.TeamID,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return results, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (m *postgresDBRepo) GetHeadToHeadResults(team1, team2 int) ([]models.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var results []models.Result

	query := `
		SELECT
			*
		FROM
			results
		WHERE
			home_team_id = $1 AND away_team_id = $2
		OR
			home_team_id = $3 AND away_team_id = $4
		ORDER BY
			result_id DESC
		`
	rows, err := m.DB.QueryContext(ctx, query, team1, team2, team2, team1)
	if err != nil {
		return results, err
	}
	defer rows.Close()
	for rows.Next() {
		var result models.Result
		err := rows.Scan(
			&result.ResultID,
			&result.SeasonID,
			&result.HomeTeam.TeamID,
			&result.HomeScore,
			&result.AwayScore,
			&result.AwayTeam.TeamID,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return results, err
		}
		results = append(results, result)
	}

	return results, nil
}
