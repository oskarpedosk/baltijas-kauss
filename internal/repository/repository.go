package repository

import (
	"net/url"

	"github.com/oskarpedosk/baltijas-kauss/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	GetTeam(teamID int) (models.Team, error)
	GetTeams() ([]models.Team, error)
	GetPlayers(perPage int, offset int) ([]models.Player, error)
	UpdateTeam(team models.Team) error
	AddPlayer(playerID, teamID int) error
	DropPlayer(playerID int) error
	CreateNewBadge(models.Badge) (int, error)
	GetBadgeID(url string) (int, error)
	UpdatePlayer(player models.Player) error
	UpdatePlayerBadges(models.Player, []models.Badge) error
	ResetPlayers() error
	GetRandomPlayer(random int) (models.Player, error)
	FilterPlayers(perPage int, offset int, queries url.Values) ([]models.Player, error)

	GetPlayer(playerID int) (models.Player, error)
	AssignPosition(player models.Player) error
	SwitchTeam(player models.Player) error

	GetStandings() ([]models.Standings, error)
	GetLastResults(count int) ([]models.Result, error)
	AddResult(res models.Result) error
	UpdateResult(res models.Result) error
	DeleteResult(res models.Result) error

	GetUser(userID int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, int, error)
}
