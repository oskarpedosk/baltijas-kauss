package repository

import (
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
	UpdatePlayer(player models.Player) error
	ResetPlayers() error
	GetRandomPlayer(random int) (models.Player, error)
	FilterPlayers(perPage int, offset int, filter models.Filter) ([]models.Player, error)

	GetPlayer(playerID int) (models.Player, error)
	AssignPosition(player models.Player) error
	SwitchTeam(player models.Player) error

	GetStandings() ([]models.Standings, error)
	GetLastResults(count int) ([]models.Result, error)
	AddResult(res models.Result) error
	UpdateResult(res models.Result) error
	DeleteResult(res models.Result) error

	CountRows(tableName string) (count int, err error)
	GetPaginationData(page int, perPage int, tableName string, baseURL string) (models.PaginationData, error)
	GetUser(userID int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, int, error)
}
