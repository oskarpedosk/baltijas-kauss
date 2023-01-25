package repository

import (
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	GetTeams() ([]models.Team, error)
	GetPlayers(perPage int, offset int) ([]models.Player, error)
	UpdateNBATeamInfo(team models.NBATeamInfo) error
	AddNBAPlayer(playerID, teamID int) error
	DropNBAPlayer(playerID int) error
	DropAllNBAPlayers() error
	GetRandomNBAPlayer(random int) (models.Player, error)

	GetNBAPlayerByID(id int) (models.Player, error)
	AssignNBAPlayer(player models.Player) error
	UpdateNBAPlayer(player models.Player) error

	GetStandings() ([]models.NBAStandings, error)
	GetLastResults(count int) ([]models.Result, error)
	AddNBAResult(res models.Result) error
	UpdateNBAResult(res models.Result) error
	DeleteNBAResult(res models.Result) error

	CountRows(tableName string) (count int, err error)
	GetPaginationData(page int, perPage int, tableName string, baseURL string) (models.PaginationData, error)
	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, int, error)
}
