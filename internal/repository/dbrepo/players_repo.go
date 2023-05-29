package dbrepo

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/oskarpedosk/baltijas-kauss/internal/helpers"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
)

// Updates NBA player team
func (m *postgresDBRepo) SwitchTeam(player models.Player) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	update 
		players
	set 
		team_id = $1,
		assigned_position = 0
	where
		player_id = $2
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		player.TeamID,
		player.PlayerID,
	)

	if err != nil {
		return err
	}

	return nil
}

// Assigns NBA player to a position
func (m *postgresDBRepo) AssignPosition(playerID, position int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	update 
		players
	set 
		assigned_position = $1
	where
		player_id = $2
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		position,
		playerID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) GetTeamPlayers(teamID int) ([]models.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select 
		player_id, team_id, first_name, last_name, img_url, assigned_position, overall
	from 
		players
	where
		"team_id" = $1
	order by
		"overall" desc
	`

	var players []models.Player

	row, err := m.DB.QueryContext(ctx, query, teamID)
	if err != nil {
		return players, err
	}

	defer row.Close()
	for row.Next() {
		var player models.Player
		err := row.Scan(
			&player.PlayerID,
			&player.TeamID,
			&player.FirstName,
			&player.LastName,
			&player.ImgURL,
			&player.AssignedPosition,
			&player.Overall,
		)
		if err != nil {
			return players, err
		}
		players = append(players, player)
	}

	return players, nil
}

// Drop NBA player from a team
func (m *postgresDBRepo) DropPlayer(playerID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	update 
		players
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

// Reset players team
func (m *postgresDBRepo) ResetPlayers() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	update 
		players
	set 
		team_id = 1,
		assigned_position = 0
	`

	_, err := m.DB.ExecContext(ctx, stmt)

	if err != nil {
		return err
	}

	return nil
}

// Add player to a team
func (m *postgresDBRepo) AddPlayer(playerID, teamID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	update 
		players
	set 
		team_id = $1,
		assigned_position = 0
	where
		player_id = $2
	`

	_, err := m.DB.ExecContext(ctx, stmt, teamID, playerID)

	if err != nil {
		return err
	}

	return nil
}

// Get player by ID
func (m *postgresDBRepo) GetPlayer(playerID int) (models.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select 
		*
	from 
		players
	where
		"player_id" = $1
	`
	var player models.Player

	row, err := m.DB.QueryContext(ctx, query, playerID)
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
			&player.TeamID,
			&player.AssignedPosition,
			&player.Archetype,
			&player.Height,
			&player.Weight,
			&player.NBATeam,
			&player.Nationality,
			&player.Birthdate,
			&player.Jersey,
			&player.Draft,
			&player.ImgURL,
			&player.RatingsURL,
			&player.Overall,
			&player.Attributes.OutsideScoring,
			&player.Attributes.Athleticism,
			&player.Attributes.InsideScoring,
			&player.Attributes.Playmaking,
			&player.Attributes.Defending,
			&player.Attributes.Rebounding,
			&player.Attributes.Intangibles,
			&player.Attributes.Potential,
			&player.Attributes.TotalAttributes,
			&player.Attributes.CloseShot,
			&player.Attributes.MidRangeShot,
			&player.Attributes.ThreePointShot,
			&player.Attributes.FreeThrow,
			&player.Attributes.ShotIQ,
			&player.Attributes.OffensiveConsistency,
			&player.Attributes.Speed,
			&player.Attributes.Acceleration,
			&player.Attributes.Strength,
			&player.Attributes.Vertical,
			&player.Attributes.Stamina,
			&player.Attributes.Hustle,
			&player.Attributes.OverallDurability,
			&player.Attributes.Layup,
			&player.Attributes.StandingDunk,
			&player.Attributes.DrivingDunk,
			&player.Attributes.PostHook,
			&player.Attributes.PostFade,
			&player.Attributes.PostControl,
			&player.Attributes.DrawFoul,
			&player.Attributes.Hands,
			&player.Attributes.PassAccuracy,
			&player.Attributes.BallHandle,
			&player.Attributes.SpeedWithBall,
			&player.Attributes.PassIQ,
			&player.Attributes.PassVision,
			&player.Attributes.InteriorDefense,
			&player.Attributes.PerimeterDefense,
			&player.Attributes.Steal,
			&player.Attributes.Block,
			&player.Attributes.LateralQuickness,
			&player.Attributes.HelpDefenseIQ,
			&player.Attributes.PassPerception,
			&player.Attributes.DefensiveConsistency,
			&player.Attributes.OffensiveRebound,
			&player.Attributes.DefensiveRebound,
			&player.BronzeBadges,
			&player.SilverBadges,
			&player.GoldBadges,
			&player.HOFBadges,
			&player.TotalBadges,
			&player.Legend,
			&player.CreatedAt,
			&player.UpdatedAt,
		)
		if player.Birthdate != "" {
			timestamp, err := time.Parse("January 2, 2006", player.Birthdate)
			if err != nil {
				log.Println(err)
			}
			player.Age = fmt.Sprintf("%d y.o.", int(time.Since(timestamp).Hours()/24/365))
		}

		if err != nil {
			return player, err
		}
	}

	return player, nil
}

func (m *postgresDBRepo) UpdatePlayer(player models.Player) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	UPDATE
		players
	SET
		first_name = $1,
		last_name = $2,
		primary_position = $3,
		secondary_position = $4,
		archetype = $5,
		height = $6,
		weight = $7,
		nba_team = $8,
		nationality = $9,
		birthdate = $10,
		jersey = $11,
		draft = $12,
		img_url = $13,
		ratings_url = $14,
		overall = $15,
		"attributes/OutsideScoring" = $16,
		"attributes/Athleticism" = $17,
		"attributes/InsideScoring" = $18,
		"attributes/Playmaking" = $19,
		"attributes/Defending" = $20,
		"attributes/Rebounding" = $21,
		"attributes/Intangibles" = $22,
		"attributes/Potential" = $23,
		"attributes/TotalAttributes" = $24,
		"attributes/CloseShot" = $25,
		"attributes/MidRangeShot" = $26,
		"attributes/ThreePointShot" = $27,
		"attributes/FreeThrow" = $28,
		"attributes/ShotIQ" = $29,
		"attributes/OffensiveConsistency" = $30,
		"attributes/Speed" = $31,
		"attributes/Acceleration" = $32,
		"attributes/Strength" = $33,
		"attributes/Vertical" = $34,
		"attributes/Stamina" = $35,
		"attributes/Hustle" = $36,
		"attributes/OverallDurability" = $37,
		"attributes/Layup" = $38,
		"attributes/StandingDunk" = $39,
		"attributes/DrivingDunk" = $40,
		"attributes/PostHook" = $41,
		"attributes/PostFade" = $42,
		"attributes/PostControl" = $43,
		"attributes/DrawFoul" = $44,
		"attributes/Hands" = $45,
		"attributes/PassAccuracy" = $46,
		"attributes/BallHandle" = $47,
		"attributes/SpeedwithBall" = $48,
		"attributes/PassIQ" = $49,
		"attributes/PassVision" = $50,
		"attributes/InteriorDefense" = $51,
		"attributes/PerimeterDefense" = $52,
		"attributes/Steal" = $53,
		"attributes/Block" = $54,
		"attributes/LateralQuickness" = $55,
		"attributes/HelpDefenseIQ" = $56,
		"attributes/PassPerception" = $57,
		"attributes/DefensiveConsistency" = $58,
		"attributes/OffensiveRebound" = $59,
		"attributes/DefensiveRebound" = $60,
		bronze_badges = $61,
		silver_badges = $62,
		gold_badges = $63,
		hof_badges = $64,
		total_badges = $65,
		updated_at = now()
	WHERE
		player_id = $66
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		player.FirstName,
		player.LastName,
		player.PrimaryPosition,
		player.SecondaryPosition,
		player.Archetype,
		helpers.NewNullInt(*player.Height),
		helpers.NewNullInt(*player.Weight),
		player.NBATeam,
		player.Nationality,
		player.Birthdate,
		player.Jersey,
		player.Draft,
		player.ImgURL,
		player.RatingsURL,
		player.Overall,
		player.Attributes.OutsideScoring,
		player.Attributes.Athleticism,
		player.Attributes.InsideScoring,
		player.Attributes.Playmaking,
		player.Attributes.Defending,
		player.Attributes.Rebounding,
		player.Attributes.Intangibles,
		player.Attributes.Potential,
		player.Attributes.TotalAttributes,
		player.Attributes.CloseShot,
		player.Attributes.MidRangeShot,
		player.Attributes.ThreePointShot,
		player.Attributes.FreeThrow,
		player.Attributes.ShotIQ,
		player.Attributes.OffensiveConsistency,
		player.Attributes.Speed,
		player.Attributes.Acceleration,
		player.Attributes.Strength,
		player.Attributes.Vertical,
		player.Attributes.Stamina,
		player.Attributes.Hustle,
		player.Attributes.OverallDurability,
		player.Attributes.Layup,
		player.Attributes.StandingDunk,
		player.Attributes.DrivingDunk,
		player.Attributes.PostHook,
		player.Attributes.PostFade,
		player.Attributes.PostControl,
		player.Attributes.DrawFoul,
		player.Attributes.Hands,
		player.Attributes.PassAccuracy,
		player.Attributes.BallHandle,
		player.Attributes.SpeedWithBall,
		player.Attributes.PassIQ,
		player.Attributes.PassVision,
		player.Attributes.InteriorDefense,
		player.Attributes.PerimeterDefense,
		player.Attributes.Steal,
		player.Attributes.Block,
		player.Attributes.LateralQuickness,
		player.Attributes.HelpDefenseIQ,
		player.Attributes.PassPerception,
		player.Attributes.DefensiveConsistency,
		player.Attributes.OffensiveRebound,
		player.Attributes.DefensiveRebound,
		player.BronzeBadges,
		player.SilverBadges,
		player.GoldBadges,
		player.HOFBadges,
		player.TotalBadges,
		player.PlayerID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) UpdatePlayerBadges(player models.Player, badges []models.Badge) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	delete from 
		players_badges
	where
		player_id = $1
	`

	_, err := m.DB.ExecContext(ctx, stmt, player.PlayerID)

	if err != nil {
		return err
	}

	for _, badge := range badges {
		badgeID, err := m.GetBadgeID(badge.URL)
		if err != nil {
			return err
		}
		if badgeID == 0 {
			badgeID, err = m.CreateNewBadge(badge)
			if err != nil {
				return err
			}
		}
		stmt := `
		INSERT INTO
			players_badges
		(player_id, badge_id, first_name, last_name, name, level) 
		values ($1, $2, $3, $4, $5, $6)`

		_, err = m.DB.ExecContext(ctx, stmt,
			player.PlayerID,
			badgeID,
			player.FirstName,
			player.LastName,
			badge.Name,
			badge.Level,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (m *postgresDBRepo) GetBadgeID(url string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT badge_id
	FROM 
		badges
	WHERE
		url = $1
	`
	badgeID := 0

	row, err := m.DB.QueryContext(ctx, query, url)
	if err != nil {
		return 0, err
	}

	defer row.Close()
	for row.Next() {
		err := row.Scan(
			&badgeID,
		)
		if err != nil {
			return 0, err
		}
	}

	return badgeID, nil
}

func (m *postgresDBRepo) CreateNewBadge(badge models.Badge) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO
			badges
		(name, type, info, img_id, url) 
		values ($1, $2, $3, $4, $5)
		RETURNING badge_id`

	var badgeID int
	err := m.DB.QueryRowContext(ctx, stmt,
		badge.Name,
		badge.Type,
		badge.Info,
		badge.ImgID,
		badge.URL,
	).Scan(&badgeID)

	if err != nil {
		return 0, err
	}

	return badgeID, nil
}

// Filter players
func (m *postgresDBRepo) GetPlayers(filter models.Filter) ([]models.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var players []models.Player

	query := `
	SELECT *
	FROM players
	WHERE ($1 = 0 OR team_id = $1)
	AND $2 <= overall AND overall <= $3
	AND $4 <= height AND height <= $5
	AND $6 <= weight AND weight <= $7
	AND $8 <= "attributes/ThreePointShot" AND "attributes/ThreePointShot" <= $9
	AND $10 <= "attributes/DrivingDunk" AND "attributes/DrivingDunk" <= $11
	AND $12 <= "attributes/Athleticism" AND "attributes/Athleticism" <= $13
	AND $14 <= "attributes/PerimeterDefense" AND "attributes/PerimeterDefense" <= $15
	AND $16 <= "attributes/InteriorDefense" AND "attributes/InteriorDefense" <= $17
	AND $18 <= "attributes/Rebounding" AND "attributes/Rebounding" <= $19
	AND (($20 = 1 AND (primary_position = 'PG' OR secondary_position = 'PG'))
		OR ($21 = 1 AND (primary_position = 'SG' OR secondary_position = 'SG'))
		OR ($22 = 1 AND (primary_position = 'SF' OR secondary_position = 'SF'))
		OR ($23 = 1 AND (primary_position = 'PF' OR secondary_position = 'PF'))
		OR ($24 = 1 AND (primary_position = 'C' OR secondary_position = 'C')))
	AND lower(first_name || '+' || last_name) LIKE '%' || lower($25) || '%'
	AND legend != $26
	ORDER BY ` + filter.Col1 + ` ` + filter.Order + `, ` + filter.Col2 + ` DESC
	LIMIT $27
	OFFSET $28
	`

	rows, err := m.DB.QueryContext(ctx, query,
		filter.TeamID,
		filter.OverallMin,
		filter.OverallMax,
		filter.HeightMin,
		filter.HeightMax,
		filter.WeightMin,
		filter.WeightMax,
		filter.ThreePointShotMin,
		filter.ThreePointShotMax,
		filter.DrivingDunkMin,
		filter.DrivingDunkMax,
		filter.AthleticismMin,
		filter.AthleticismMax,
		filter.PerimeterDefenseMin,
		filter.PerimeterDefenseMax,
		filter.InteriorDefenseMin,
		filter.InteriorDefenseMax,
		filter.ReboundingMin,
		filter.ReboundingMax,
		filter.Position1,
		filter.Position2,
		filter.Position3,
		filter.Position4,
		filter.Position5,
		filter.Search,
		filter.Era,
		filter.Limit,
		filter.Offset,
	)
	if err != nil {
		return players, err
	}

	defer rows.Close()
	for rows.Next() {
		var player models.Player
		err := rows.Scan(
			&player.PlayerID,
			&player.FirstName,
			&player.LastName,
			&player.PrimaryPosition,
			&player.SecondaryPosition,
			&player.TeamID,
			&player.AssignedPosition,
			&player.Archetype,
			&player.Height,
			&player.Weight,
			&player.NBATeam,
			&player.Nationality,
			&player.Birthdate,
			&player.Jersey,
			&player.Draft,
			&player.ImgURL,
			&player.RatingsURL,
			&player.Overall,
			&player.Attributes.OutsideScoring,
			&player.Attributes.Athleticism,
			&player.Attributes.InsideScoring,
			&player.Attributes.Playmaking,
			&player.Attributes.Defending,
			&player.Attributes.Rebounding,
			&player.Attributes.Intangibles,
			&player.Attributes.Potential,
			&player.Attributes.TotalAttributes,
			&player.Attributes.CloseShot,
			&player.Attributes.MidRangeShot,
			&player.Attributes.ThreePointShot,
			&player.Attributes.FreeThrow,
			&player.Attributes.ShotIQ,
			&player.Attributes.OffensiveConsistency,
			&player.Attributes.Speed,
			&player.Attributes.Acceleration,
			&player.Attributes.Strength,
			&player.Attributes.Vertical,
			&player.Attributes.Stamina,
			&player.Attributes.Hustle,
			&player.Attributes.OverallDurability,
			&player.Attributes.Layup,
			&player.Attributes.StandingDunk,
			&player.Attributes.DrivingDunk,
			&player.Attributes.PostHook,
			&player.Attributes.PostFade,
			&player.Attributes.PostControl,
			&player.Attributes.DrawFoul,
			&player.Attributes.Hands,
			&player.Attributes.PassAccuracy,
			&player.Attributes.BallHandle,
			&player.Attributes.SpeedWithBall,
			&player.Attributes.PassIQ,
			&player.Attributes.PassVision,
			&player.Attributes.InteriorDefense,
			&player.Attributes.PerimeterDefense,
			&player.Attributes.Steal,
			&player.Attributes.Block,
			&player.Attributes.LateralQuickness,
			&player.Attributes.HelpDefenseIQ,
			&player.Attributes.PassPerception,
			&player.Attributes.DefensiveConsistency,
			&player.Attributes.OffensiveRebound,
			&player.Attributes.DefensiveRebound,
			&player.BronzeBadges,
			&player.SilverBadges,
			&player.GoldBadges,
			&player.HOFBadges,
			&player.TotalBadges,
			&player.Legend,
			&player.CreatedAt,
			&player.UpdatedAt,
		)
		if err != nil {
			return players, err
		}

		players = append(players, player)
	}

	return players, nil
}

func (m *postgresDBRepo) CountPlayers() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT COUNT(*) FROM players;
	`

	var count int

	err := m.DB.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *postgresDBRepo) GetPlayerBadges(playerID int) ([]models.Badge, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var badges []models.Badge

	query := `
	SELECT badges.badge_id, badges.name, badges.type, badges.info, badges.img_id, badges.url 
	FROM players_badges 
	JOIN badges ON players_badges.badge_id = badges.badge_id 
	WHERE players_badges.player_id = $1;
	`

	rows, err := m.DB.QueryContext(ctx, query, playerID)
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
			&badge.ImgID,
			&badge.URL,
		)
		if err != nil {
			return badges, err
		}
		badges = append(badges, badge)
	}

	return badges, nil
}

func (m *postgresDBRepo) GetADP(playerID int) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sum := 0.0
	counter := 0.0

	query := `
	SELECT
		pick
	FROM
		drafts
	WHERE
		player_id = $1
	`

	rows, err := m.DB.QueryContext(ctx, query, playerID)
	if err != nil {
		return sum, err
	}

	defer rows.Close()
	for rows.Next() {
		var pick int
		err := rows.Scan(&pick)
		if err != nil {
			return sum, err
		}
		sum += float64(pick)
		counter++
	}

	if sum == 0.0 {
		return sum, nil
	}

	adp := math.Round(sum/counter*10) / 10

	return adp, nil
}

// Create new player to database
func (m *postgresDBRepo) CreateNewPlayer(ratingsURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT MAX(player_id) FROM players`

	var playerID int
	err := m.DB.QueryRowContext(ctx, stmt).Scan(&playerID)
	if err != nil {
		return err
	}

	stmt = `
	INSERT INTO
		players (player_id, first_name, last_name, ratings_url, overall, height, weight, created_at, updated_at) 
	VALUES 
		($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`

	_, err = m.DB.ExecContext(ctx, stmt, playerID+1, "New", "Player", ratingsURL, 1, helpers.NewNullInt(0), helpers.NewNullInt(0))
	if err != nil {
		return err
	}

	return nil
}

// Delete player from db
func (m *postgresDBRepo) DeletePlayer(playerID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	delete from 
		players
	where
		player_id = $1
	`

	_, err := m.DB.ExecContext(ctx, stmt, playerID)

	if err != nil {
		return err
	}

	return nil
}
