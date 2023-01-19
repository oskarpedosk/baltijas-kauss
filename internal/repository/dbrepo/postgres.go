package dbrepo

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var playerCount = 150

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

	stmt := `insert into nba_results (home_team_id, home_score, away_score, away_team_id, timestamp) 
	values ($1, $2, $3, $4, CURRENT_TIMESTAMP)`

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

// Update NBA result
func (m *postgresDBRepo) UpdateNBAResult(res models.Result) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	update 
		nba_results
	set 
		home_team_id = $1,
		home_score = $2,
		away_score = $3,
		away_team_id = $4
	where
		timestamp = $5
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.HomeTeam,
		res.HomeScore,
		res.AwayScore,
		res.AwayTeam,
		res.Time,
	)

	if err != nil {
		return err
	}

	return nil
}

// Delete NBA result from database
func (m *postgresDBRepo) DeleteNBAResult(res models.Result) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	delete from 
		nba_results
	where
		timestamp = $1
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.Time,
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

// Assigns NBA player to a position
func (m *postgresDBRepo) AssignNBAPlayer(player models.NBAPlayer) error {
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

// Display NBA players with player limit
func (m *postgresDBRepo) GetNBAPlayersWithBadges() ([]models.NBAPlayer, error) {
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

func (m *postgresDBRepo) CountPlayers() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select count(*)
	from
		"nba_players"
	`
	var count int
	rows := m.DB.QueryRowContext(ctx, query)
	err := rows.Scan(&count)
	if err != nil {
		return count, err
	}

	return count, nil
}

// Display all NBA players without badges
func (m *postgresDBRepo) GetNBAPlayersWithoutBadges(offset int) ([]models.NBAPlayer, error) {
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
	limit 30
	offset $1
	`

	rows, err := m.DB.QueryContext(ctx, query, offset)
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

// Get team info
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

// Display all NBA standings
func (m *postgresDBRepo) GetNBAStandings() ([]models.NBAStandings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var standingsSlice []models.NBAStandings

	for i := 1; i < 5; i++ {
		homeWins := 0
		homeLosses := 0
		awayWins := 0
		awayLosses := 0
		basketsFor := 0
		basketsAgainst := 0
		games := ""
		query := `
			select 
				*
			from 
				nba_results
			where
				"home_team_id" = $1 or "away_team_id" = $1
			order by
				timestamp asc
			`

		rows, err := m.DB.QueryContext(ctx, query, i)
		if err != nil {
			return standingsSlice, err
		}

		var singleGame models.Result
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(
				&singleGame.HomeTeam,
				&singleGame.HomeScore,
				&singleGame.AwayScore,
				&singleGame.AwayTeam,
				&singleGame.Time,
			)
			if err != nil {
				return standingsSlice, err
			}

			if singleGame.HomeTeam == i {
				basketsFor += singleGame.HomeScore
				basketsAgainst += singleGame.AwayScore
				if singleGame.HomeScore > singleGame.AwayScore {
					homeWins += 1
					games = "W" + games
				} else {
					homeLosses += 1
					games = "L" + games
				}
			} else {
				basketsFor += singleGame.AwayScore
				basketsAgainst += singleGame.HomeScore
				if singleGame.AwayScore > singleGame.HomeScore {
					awayWins += 1
					games = "W" + games
				} else {
					awayLosses += 1
					games = "L" + games
				}
			}
		}

		x := 5
		chars := []rune(games)

		if len(chars) < x {
			x = len(chars)
		}

		lastGames := []string{"", "", "", "", ""}
		y := 0

		for i := x; i > 0; i-- {
			lastGames[i-1] = string(chars[y])
			y++
		}

		streak := ""
		streakCount := 0

		if games != "" {
			streak = string(chars[0])
			if len(chars) > 1 {
				for i := 0; i < len(chars); i++ {
					if string(chars[i]) != streak {
						break
					} else {
						streakCount += 1
					}
				}
			} else {
				streakCount = 1
			}
		}

		if err = rows.Err(); err != nil {
			return standingsSlice, err
		}

		totalWins := (homeWins + awayWins)
		totalGames := (homeWins + homeLosses + awayWins + awayLosses)
		winPercentage := 0

		forAvg := 0.0
		againstAvg := 0.0

		if totalGames != 0 {
			winPercentage = totalWins * 1000 / totalGames
			forAvg = toFixed(float64(basketsFor)/float64(totalGames), 1)
			againstAvg = toFixed(float64(basketsAgainst)/float64(totalGames), 1)
		}

		teamStandings := models.NBAStandings{
			TeamID:         i,
			WinPercentage:  winPercentage,
			Played:         totalGames,
			TotalWins:      homeWins + awayWins,
			TotalLosses:    homeLosses + awayLosses,
			HomeWins:       homeWins,
			HomeLosses:     homeLosses,
			AwayWins:       awayWins,
			AwayLosses:     awayLosses,
			Streak:         streak,
			StreakCount:    streakCount,
			BasketsFor:     basketsFor,
			BasketsAgainst: basketsAgainst,
			BasketsSum:     basketsFor - basketsAgainst,
			ForAvg:         float64(forAvg),
			AgainstAvg:     float64(againstAvg),
			LastFive:       lastGames,
		}
		standingsSlice = append(standingsSlice, teamStandings)
	}

	orderedSlice := order(standingsSlice)
	return orderedSlice, nil
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func order(slice []models.NBAStandings) []models.NBAStandings {
	for i := 0; i < len(slice)-1; i++ {
		if slice[i].WinPercentage < slice[i+1].WinPercentage {
			slice[i], slice[i+1] = slice[i+1], slice[i]
			order(slice)
		}
		if slice[i].WinPercentage == slice[i+1].WinPercentage {
			if slice[i].BasketsSum < slice[i+1].BasketsSum {
				slice[i], slice[i+1] = slice[i+1], slice[i]
			}
		}

	}
	return slice
}

// Display however many last games
func (m *postgresDBRepo) GetLastResults(count int) ([]models.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var resultsSlice []models.Result

	query := `
			select 
				*
			from 
				nba_results
			order by
				timestamp asc
			`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return resultsSlice, err
	}

	var singleGame models.Result
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&singleGame.HomeTeam,
			&singleGame.HomeScore,
			&singleGame.AwayScore,
			&singleGame.AwayTeam,
			&singleGame.Time,
		)
		if err != nil {
			return resultsSlice, err
		}

		layout := "02/01/2006 15:04"

		result := models.Result{
			HomeTeam:   singleGame.HomeTeam,
			HomeScore:  singleGame.HomeScore,
			AwayScore:  singleGame.AwayScore,
			AwayTeam:   singleGame.AwayTeam,
			Time:       singleGame.Time,
			TimeString: singleGame.Time.Round(15 * time.Minute).Format(layout),
		}

		resultsSlice = append([]models.Result{result}, resultsSlice...)
	}

	if err = rows.Err(); err != nil {
		return resultsSlice, err
	}

	return resultsSlice, nil
}

// Drop NBA player from a team
func (m *postgresDBRepo) DropNBAPlayer(playerID int) error {
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
func (m *postgresDBRepo) DropAllNBAPlayers() error {
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

// Add NBA player directly to a team
func (m *postgresDBRepo) AddNBAPlayer(playerID, teamID int) error {
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

// Add a random NBA player directly to a team
func (m *postgresDBRepo) GetRandomNBAPlayer(random int) (models.NBAPlayer, error) {
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
func (m *postgresDBRepo) GetNBABadges() ([]models.Badge, error) {
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
func (m *postgresDBRepo) GetNBAPlayersBadges() ([]models.PlayersBadges, error) {
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

// Returns user by ID
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select 
		user_id, first_name, last_name, email, password, access_level
	from
		users
	where
		user_id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(
		&u.UserID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

// Updates a user in a database
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	update
		users
	set
		first_name = $1,
		last_name = $2,
		email = $3,
		access_level = $4 
	`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
	)
	if err != nil {
		return err
	}

	return nil
}

// Authenticate authenticates a user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string
	var accessLevel int

	row := m.DB.QueryRowContext(ctx, "select user_id, password, access_level from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword, &accessLevel)
	if err != nil {
		return id, "", 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", 0, errors.New("incorrect password")
	} else if err != nil {
		return 0, "", 0, err
	}

	return id, hashedPassword, accessLevel, nil
}

// Display NBA player by ID
func (m *postgresDBRepo) GetNBAPlayerByID(id int) (models.NBAPlayer, error) {
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
