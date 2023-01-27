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

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// Updates NBA team info
func (m *postgresDBRepo) UpdateTeam(team models.Team) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update teams set name = $2, abbreviation = $3, team_color1 = $4, team_color2 = $5, dark_text = $6 where team_id = $1`

	_, err := m.DB.ExecContext(ctx, stmt,
		team.TeamID,
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

// Update NBA result
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

// Delete NBA result from database
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

// Updates NBA player team
func (m *postgresDBRepo) UpdatePlayer(player models.Player) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	update 
		players
	set 
		team_id = $2,
		assigned_position = 0
	where
		player_id = $1
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		player.PlayerID,
		player.TeamID,
	)

	if err != nil {
		return err
	}

	return nil
}

// Assigns NBA player to a position
func (m *postgresDBRepo) AssignPosition(player models.Player) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Remove previous player from position
	stmt := `
	update 
		players
	set 
		assigned_position = 0
	where
		team_id = $1
	and
		assigned_position = $2
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		player.TeamID,
		player.AssignedPosition,
	)

	if err != nil {
		return err
	}

	// Add new player to position
	stmt = `
	update 
		players
	set 
		assigned_position = $3
	where
		player_id = $1
	and
		team_id = $2
	`

	_, err = m.DB.ExecContext(ctx, stmt,
		player.PlayerID,
		player.TeamID,
		player.AssignedPosition,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) CountRows(tableName string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`select count(*) from %s`, tableName)

	var count int
	rows := m.DB.QueryRowContext(ctx, query)
	err := rows.Scan(&count)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (m *postgresDBRepo) GetPaginationData(page int, perPage int, tableName string, baseURL string) (models.PaginationData, error) {
	// Calculate total pages
	totalRows, err := m.CountRows(tableName)
	if err != nil {
		return models.PaginationData{}, err
	}
	totalPages := math.Ceil(float64(totalRows) / float64(perPage))

	// Calculate offset
	offset := (page - 1) * perPage

	pagination := models.PaginationData{
		NextPage:     page + 1,
		PreviousPage: page - 1,
		CurrentPage:  page,
		TotalPages:   int(totalPages),
		TwoBefore:    page - 2,
		TwoAfter:     page + 2,
		ThreeAfter:   page + 3,
		Offset:       offset,
		BaseURL:      baseURL,
	}
	return pagination, nil
}

// Get all players pagination limit
func (m *postgresDBRepo) GetPlayers(perPage int, offset int) ([]models.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var players []models.Player

	query := `
	select
		*
	from 
		players
	order by
		"overall" desc,
		"attributes/TotalAttributes" desc
	limit $1
	offset $2
	`

	rows, err := m.DB.QueryContext(ctx, query, perPage, offset)
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
			&player.ImgID,
			&player.RatingsURL,
			&player.Overall,
			&player.AttributesOutsideScoring,
			&player.AttributesAthleticism,
			&player.AttributesInsideScoring,
			&player.AttributesPlaymaking,
			&player.AttributesDefending,
			&player.AttributesRebounding,
			&player.AttributesIntangibles,
			&player.AttributesPotential,
			&player.AttributesTotalAttributes,
			&player.AttributesCloseShot,
			&player.AttributesMidRangeShot,
			&player.AttributesThreePointShot,
			&player.AttributesFreeThrow,
			&player.AttributesShotIQ,
			&player.AttributesOffensiveConsistency,
			&player.AttributesSpeed,
			&player.AttributesAcceleration,
			&player.AttributesStrength,
			&player.AttributesVertical,
			&player.AttributesStamina,
			&player.AttributesHustle,
			&player.AttributesOverallDurability,
			&player.AttributesLayup,
			&player.AttributesStandingDunk,
			&player.AttributesDrivingDunk,
			&player.AttributesPostHook,
			&player.AttributesPostFade,
			&player.AttributesPostControl,
			&player.AttributesDrawFoul,
			&player.AttributesHands,
			&player.AttributesPassAccuracy,
			&player.AttributesBallHandle,
			&player.AttributesSpeedWithBall,
			&player.AttributesPassIQ,
			&player.AttributesPassVision,
			&player.AttributesInteriorDefense,
			&player.AttributesPerimeterDefense,
			&player.AttributesSteal,
			&player.AttributesBlock,
			&player.AttributesLateralQuickness,
			&player.AttributesHelpDefenseIQ,
			&player.AttributesPassPerception,
			&player.AttributesDefensiveConsistency,
			&player.AttributesOffensiveRebound,
			&player.AttributesDefensiveRebound,
			&player.BronzeBadges,
			&player.SilverBadges,
			&player.GoldBadges,
			&player.HOFBadges,
			&player.TotalBadges,
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

// Get all teams
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
			&team.DarkText,
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

// Get all teams
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
			&team.DarkText,
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

// Get results and calculate standings
func (m *postgresDBRepo) GetStandings() ([]models.Standings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var standingsSlice []models.Standings

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
				results
			where
				"home_team_id" = $1 or "away_team_id" = $1
			order by
				created_at asc
			`

		rows, err := m.DB.QueryContext(ctx, query, i)
		if err != nil {
			return standingsSlice, err
		}

		var singleGame models.Result
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(
				&singleGame.ResultID,
				&singleGame.Season,
				&singleGame.HomeTeamID,
				&singleGame.HomeScore,
				&singleGame.AwayScore,
				&singleGame.AwayTeamID,
				&singleGame.CreatedAt,
				&singleGame.UpdatedAt,
			)
			if err != nil {
				return standingsSlice, err
			}

			if singleGame.HomeTeamID == i {
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

		teamStandings := models.Standings{
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

func order(slice []models.Standings) []models.Standings {
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
				results
			order by
				created_at asc
			`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return resultsSlice, err
	}

	var singleGame models.Result
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&singleGame.ResultID,
			&singleGame.Season,
			&singleGame.HomeTeamID,
			&singleGame.HomeScore,
			&singleGame.AwayScore,
			&singleGame.AwayTeamID,
			&singleGame.CreatedAt,
			&singleGame.UpdatedAt,
		)
		if err != nil {
			return resultsSlice, err
		}

		result := models.Result{
			HomeTeamID: singleGame.HomeTeamID,
			HomeScore:  singleGame.HomeScore,
			AwayScore:  singleGame.AwayScore,
			AwayTeamID: singleGame.AwayTeamID,
			CreatedAt:  singleGame.CreatedAt.Round(15 * time.Minute),
		}

		resultsSlice = append([]models.Result{result}, resultsSlice...)
	}

	if err = rows.Err(); err != nil {
		return resultsSlice, err
	}

	return resultsSlice, nil
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

// Add a random player to a team
func (m *postgresDBRepo) GetRandomPlayer(random int) (models.Player, error) {
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

// Get user by ID
func (m *postgresDBRepo) GetUser(userID int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select 
		*
	from
		users
	where
		user_id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, userID)

	var u models.User
	err := row.Scan(
		&u.UserID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.ImgID,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
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

// Authenticates a user
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
			&player.ImgID,
			&player.RatingsURL,
			&player.Overall,
			&player.AttributesOutsideScoring,
			&player.AttributesAthleticism,
			&player.AttributesInsideScoring,
			&player.AttributesPlaymaking,
			&player.AttributesDefending,
			&player.AttributesRebounding,
			&player.AttributesIntangibles,
			&player.AttributesPotential,
			&player.AttributesTotalAttributes,
			&player.AttributesCloseShot,
			&player.AttributesMidRangeShot,
			&player.AttributesThreePointShot,
			&player.AttributesFreeThrow,
			&player.AttributesShotIQ,
			&player.AttributesOffensiveConsistency,
			&player.AttributesSpeed,
			&player.AttributesAcceleration,
			&player.AttributesStrength,
			&player.AttributesVertical,
			&player.AttributesStamina,
			&player.AttributesHustle,
			&player.AttributesOverallDurability,
			&player.AttributesLayup,
			&player.AttributesStandingDunk,
			&player.AttributesDrivingDunk,
			&player.AttributesPostHook,
			&player.AttributesPostFade,
			&player.AttributesPostControl,
			&player.AttributesDrawFoul,
			&player.AttributesHands,
			&player.AttributesPassAccuracy,
			&player.AttributesBallHandle,
			&player.AttributesSpeedWithBall,
			&player.AttributesPassIQ,
			&player.AttributesPassVision,
			&player.AttributesInteriorDefense,
			&player.AttributesPerimeterDefense,
			&player.AttributesSteal,
			&player.AttributesBlock,
			&player.AttributesLateralQuickness,
			&player.AttributesHelpDefenseIQ,
			&player.AttributesPassPerception,
			&player.AttributesDefensiveConsistency,
			&player.AttributesOffensiveRebound,
			&player.AttributesDefensiveRebound,
			&player.BronzeBadges,
			&player.SilverBadges,
			&player.GoldBadges,
			&player.HOFBadges,
			&player.TotalBadges,
			&player.CreatedAt,
			&player.UpdatedAt,
		)
		if err != nil {
			return player, err
		}
	}

	return player, nil
}
