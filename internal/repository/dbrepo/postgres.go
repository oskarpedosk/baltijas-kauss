package dbrepo

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
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
func (m *postgresDBRepo) SwitchTeam(player models.Player) error {
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
			&player.CreatedAt,
			&player.UpdatedAt,
		)
		if player.Birthdate != "" {
			timestamp, err := time.Parse("January 2, 2006", player.Birthdate)
			if err != nil {
				fmt.Println(err)
			}
			player.Age = fmt.Sprintf("%dy.o.", int(time.Since(timestamp).Hours()/24/365))
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
		*player.Height,
		*player.Weight,
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

// Filter players with pagination limit
func (m *postgresDBRepo) FilterPlayers(perPage int, offset int, queries url.Values) ([]models.Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := models.Filter{
		OverallMin: 1,
		OverallMax: 99,
	}
	var players []models.Player

	for key, value := range queries {
		fmt.Println(key, " => ", value)
		queryInt, err := strconv.Atoi(value[0])
		if err != nil {
			continue
		}
		switch key {
		case "ovrl":
			filter.OverallMin = queryInt
		case "ovrh":
			filter.OverallMax = queryInt
		} 
	}

	query := `
	SELECT *
	FROM players
	WHERE $1 <= overall AND overall <= $2
	ORDER BY "overall" DESC, "attributes/TotalAttributes" DESC
	LIMIT $3
	OFFSET $4
	`

	rows, err := m.DB.QueryContext(ctx, query,
		filter.OverallMin,
		filter.OverallMax,
		perPage,
		offset,
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
