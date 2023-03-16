package dbrepo

import (
	"context"
	"time"

	"github.com/oskarpedosk/baltijas-kauss/internal/models"
)

func (m *postgresDBRepo) AddResult(res models.Result) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into results (home_team_id, home_score, away_score, away_team_id) 
	values ($1, $2, $3, $4)`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.HomeTeamID,
		res.HomeScore,
		res.AwayScore,
		res.AwayTeamID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) UpdateResult(res models.Result) error {
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
		created_at = $5
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.HomeTeamID,
		res.HomeScore,
		res.AwayScore,
		res.AwayTeamID,
		res.CreatedAt,
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
		created_at = $1
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.CreatedAt,
	)

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
	} else {
		query := `
			SELECT *
			FROM results
			WHERE season_id = $1
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
				&result.Season,
				&result.HomeTeamID,
				&result.HomeScore,
				&result.AwayScore,
				&result.AwayTeamID,
				&result.CreatedAt,
				&result.UpdatedAt,
			)
			if err != nil {
				return results, err
			}
			results = append(results, result)
		}
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
			&result.Season,
			&result.HomeTeamID,
			&result.HomeScore,
			&result.AwayScore,
			&result.AwayTeamID,
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
