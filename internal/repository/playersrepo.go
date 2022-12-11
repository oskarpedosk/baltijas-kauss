package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/oskarpedosk/baltijas-kauss/internal/config"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
)


type PlayersRepo interface {
	GetNBAPlayerByID(id int) (models.NBAPlayer, error)
	UpdateNBAPlayer(player models.NBAPlayer) error 
	GetNBAPlayersBadges() ([]models.PlayersBadges, error)
	GetNBABadges() ([]models.Badge, error)
	GetRandomNBAPlayer(random int) (models.NBAPlayer, error)
	GetNBAPlayersWithBadges() ([]models.NBAPlayer, error)
	GetNBAPlayersWithoutBadges() ([]models.NBAPlayer, error)
}

type playersRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPlayersRepo(conn *sql.DB, a *config.AppConfig) PlayersRepo {
	return &playersRepo{
		App: a,
		DB:  conn,
	}
}

// Updates NBA player team
func (m *playersRepo) UpdateNBAPlayer(player models.NBAPlayer) error {
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

// Display NBA players with player limit
func (m *playersRepo) GetNBAPlayersWithBadges() ([]models.NBAPlayer, error) {
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
	limit
		$1
	`

	rows, err := m.DB.QueryContext(ctx, query, playerCount)
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
			&player.BronzeBadgesCount,
			&player.SilverBadgesCount,
			&player.GoldBadgesCount,
			&player.HOFBadgesCount,
			&player.TotalBadgesCount,
			&player.Assigned,
		)
		if err != nil {
			return players, err
		}

		query2 := `
		select 
			"badge_id", "level"
		from 
			nba_players_badges
		where
			"player_id" = $1
		`

		rows, err := m.DB.QueryContext(ctx, query2, player.PlayerID)
		if err != nil {
			return players, err
		}

		bronzeBadges := player.BronzeBadges
		silverBadges := player.SilverBadges
		goldBadges := player.GoldBadges
		hofBadges := player.HOFBadges

		defer rows.Close()
		for rows.Next() {
			var badge models.PlayersBadges
			err := rows.Scan(
				&badge.BadgeID,
				&badge.Level,
			)
			if err != nil {
				return players, err
			}

			query3 := `
			select 
				*
			from 
				nba_badges
			where
				"badge_id" = $1
			`

			rows, err := m.DB.QueryContext(ctx, query3, badge.BadgeID)
			if err != nil {
				return players, err
			}

			defer rows.Close()
			for rows.Next() {
				var playerBadge models.Badge
				err := rows.Scan(
					&playerBadge.BadgeID,
					&playerBadge.Name,
					&playerBadge.Type,
					&playerBadge.Info,
					&playerBadge.BronzeUrl,
					&playerBadge.SilverUrl,
					&playerBadge.GoldUrl,
					&playerBadge.HOFUrl,
				)
				if err != nil {
					return players, err
				}

				if badge.Level == "Bronze" {
					bronzeBadges = append(bronzeBadges, playerBadge)
				}
				if badge.Level == "Silver" {
					silverBadges = append(silverBadges, playerBadge)
				}
				if badge.Level == "Gold" {
					goldBadges = append(goldBadges, playerBadge)
				}
				if badge.Level == "HOF" {
					hofBadges = append(hofBadges, playerBadge)
				}

			}

			player.BronzeBadges = bronzeBadges
			player.SilverBadges = silverBadges
			player.GoldBadges = goldBadges
			player.HOFBadges = hofBadges

			if err = rows.Err(); err != nil {
				return players, err
			}

		}

		if err = rows.Err(); err != nil {
			return players, err
		}
		players = append(players, player)
	}

	return players, nil
}

// Display all NBA players without badges
func (m *playersRepo) GetNBAPlayersWithoutBadges() ([]models.NBAPlayer, error) {
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
			&player.BronzeBadgesCount,
			&player.SilverBadgesCount,
			&player.GoldBadgesCount,
			&player.HOFBadgesCount,
			&player.TotalBadgesCount,
			&player.Assigned,
		)
		if err != nil {
			return players, err
		}

		players = append(players, player)
	}

	return players, nil
}

// Add a random NBA player directly to a team
func (m *playersRepo) GetRandomNBAPlayer(random int) (models.NBAPlayer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	fmt.Println("-------------------")
	fmt.Println("random=", random)
	// Remove team and position
	stmt := `
	select
		"player_id",
		"first_name",
		"last_name",
		"primary_position",
		"secondary_position"
	from 
		nba_players
	where
		"team_id" is null
	order by
		"stats/Overall" desc,
		"stats/Total Attributes" desc
	limit
		1
	offset
		$1
	`

	row := m.DB.QueryRowContext(ctx, stmt, random)

	var player models.NBAPlayer
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

// Get all NBA badges
func (m *playersRepo) GetNBABadges() ([]models.Badge, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var badges []models.Badge

	query := `
	select 
		*
	from 
		nba_badges
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return badges, err
	}

	defer rows.Close()
	for rows.Next() {
		var badge models.Badge
		err := rows.Scan(
			&badge.BadgeID,
			&badge.Name,
			&badge.Type,
			&badge.Info,
			&badge.BronzeUrl,
			&badge.SilverUrl,
			&badge.GoldUrl,
			&badge.HOFUrl,
		)
		if err != nil {
			return badges, err
		}
		badges = append(badges, badge)
	}

	if err = rows.Err(); err != nil {
		return badges, err
	}

	return badges, nil
}

// Get all NBA badges
func (m *playersRepo) GetNBAPlayersBadges() ([]models.PlayersBadges, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var playersBadges []models.PlayersBadges

	query := `
	select 
		*
	from 
		nba_players_badges
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return playersBadges, err
	}

	defer rows.Close()
	for rows.Next() {
		var playersBadge models.PlayersBadges
		err := rows.Scan(
			&playersBadge.PlayerID,
			&playersBadge.FirstName,
			&playersBadge.LastName,
			&playersBadge.BadgeID,
			&playersBadge.Name,
			&playersBadge.Level,
		)
		if err != nil {
			return playersBadges, err
		}
		playersBadges = append(playersBadges, playersBadge)
	}

	if err = rows.Err(); err != nil {
		return playersBadges, err
	}

	return playersBadges, nil
}


// Display NBA player by ID
func (m *playersRepo) GetNBAPlayerByID(id int) (models.NBAPlayer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select 
		*
	from 
		nba_players
	where
		"player_id" = $1
	`
	var player models.NBAPlayer

	row, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return player, err
	}

	defer row.Close()
	for row.Next() {
		err := row.Scan(
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
			&player.BronzeBadgesCount,
			&player.SilverBadgesCount,
			&player.GoldBadgesCount,
			&player.HOFBadgesCount,
			&player.TotalBadgesCount,
			&player.Assigned,
		)
		if err != nil {
			return player, err
		}
	}

	return player, nil
}
