package repository

import (
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	GetTeam(teamID int) (models.Team, error)
	GetTeams() ([]models.Team, error)
	UpdateTeam(team models.Team) error
	AddPlayer(playerID, teamID int) error
	DropPlayer(playerID int) error
	CreateNewBadge(models.Badge) (int, error)
	GetBadgeID(url string) (int, error)
	UpdatePlayer(player models.Player) error
	UpdatePlayerBadges(models.Player, []models.Badge) error
	CountPlayers() (int, error)
	ResetPlayers() error
	GetRandomPlayer(random int) (models.Player, error)
	GetPlayers(filter models.Filter) ([]models.Player, error)

	GetTeamPlayers(teamID int) ([]models.Player, error)
	GetPlayer(playerID int) (models.Player, error)
	GetPlayerBadges(playerID int) ([]models.Badge, error)
	AssignPosition(playerID, position int) error
	SwitchTeam(player models.Player) error

	GetAllResults() ([]models.Result, error)
	GetSeasonResults(seasonID int) ([]models.Result, error)
	StartNewSeason() error
	GetSeasons() ([]models.Season, error)
	AddResult(res models.Result) error
	UpdateResult(res models.Result) error
	DeleteResult(res models.Result) error

	GetUser(userID int) (models.User, error)
	Authenticate(email, testPassword string) (int, string, int, error)
}
