package repository

import "github.com/oskarpedosk/baltijas-kauss/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	GetNBATeamInfo() ([]models.NBATeam, error)
	GetNBAPlayers() ([]models.NBAPlayer, error)
	UpdateNBATeamInfo(team models.NBATeamInfo) error
	AddNBAPlayer(teamID, playerID int) error
	DropNBAPlayer(playerID int) error

	GetNBAPlayerByID(id int) (models.NBAPlayer, error)
	AssignNBAPlayer(player models.NBAPlayer) error
	UpdateNBAPlayer(player models.NBAPlayer) error 
	GetNBAPlayersBadges() ([]models.PlayersBadges, error)
	GetNBABadges() ([]models.Badge, error)

	GetNBAStandings() ([]models.NBAStandings, error)
	GetLastResults(count int) ([]models.Result, error)
	AddNBAResult(res models.Result) error
	UpdateNBAResult(res models.Result) error
	DeleteNBAResult(res models.Result) error

	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)
}