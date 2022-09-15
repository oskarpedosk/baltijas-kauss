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
func (m *postgresDBRepo) UpdateNBATeamInfo(team models.NBATeamInfo) error {
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

// Updates NBA player team
func (m *postgresDBRepo) UpdateNBAPlayer(player models.NBAPlayer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	update 
		nba_players
	set 
		team_id = $2,
		assigned = $3
	where
		player_id = $1
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		player.PlayerID,
		player.TeamID,
		player.Assigned,
	)

	if err != nil {
		return err
	}

	return nil
}

// Display all NBA players
func (m *postgresDBRepo) GetNBAPlayers() ([]models.NBAPlayer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var players []models.NBAPlayer

	query := `
	select 
		*
	from 
		nba_players
	order by
		"stats/Overall" desc,
		"stats/Total Attributes" desc
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return players, err
	}

	defer rows.Close()
	for rows.Next() {
		var player models.NBAPlayer
		err := rows.Scan(
			&player.PlayerID,
			&player.FirstName,
			&player.LastName,
			&player.PrimaryPosition,
			&player.SecondaryPosition,
			&player.Archetype,
			&player.NBATeam,
			&player.Height,
			&player.Weight,
			&player.ImgUrl,
			&player.PlayerUrl,
			&player.TeamID,
			&player.StatsOverall,
			&player.StatsOutsideScoring,
			&player.StatsAthleticism,
			&player.StatsInsideScoring,
			&player.StatsPlaymaking,
			&player.StatsDefending,
			&player.StatsRebounding,
			&player.StatsCloseShot,
			&player.StatsMidRangeShot,
			&player.StatsThreePointShot,
			&player.StatsFreeThrow,
			&player.StatsShotIQ,
			&player.StatsOffensiveConsistency,
			&player.StatsSpeed,
			&player.StatsAcceleration,
			&player.StatsStrength,
			&player.StatsVertical,
			&player.StatsStamina,
			&player.StatsHustle,
			&player.StatsOverallDurability,
			&player.StatsLayup,
			&player.StatsStandingDunk,
			&player.StatsDrivingDunk,
			&player.StatsPostHook,
			&player.StatsPostFade,
			&player.StatsPostControl,
			&player.StatsDrawFoul,
			&player.StatsHands,
			&player.StatsPassAccuracy,
			&player.StatsBallHandle,
			&player.StatsSpeedWithBall,
			&player.StatsPassIQ,
			&player.StatsPassVision,
			&player.StatsInteriorDefense,
			&player.StatsPerimeterDefense,
			&player.StatsSteal,
			&player.StatsBlock,
			&player.StatsLateralQuickness,
			&player.StatsHelpDefenseIQ,
			&player.StatsPassPerception,
			&player.StatsDefensiveConsistency,
			&player.StatsOffensiveRebound,
			&player.StatsDefensiveRebound,
			&player.StatsIntangibles,
			&player.StatsPotential,
			&player.StatsTotalAttributes,
			&player.BronzeBadges,
			&player.SilverBadges,
			&player.GoldBadges,
			&player.HOFBadges,
			&player.TotalBadges,
			&player.Assigned,
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

// Display all NBA players
func (m *postgresDBRepo) GetNBATeamInfo() ([]models.NBATeam, error) {
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
