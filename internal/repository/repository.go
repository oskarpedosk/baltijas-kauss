package repository

import "github.com/oskarpedosk/baltijas-kauss/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	UpdateNBATeamInfo(team models.NBATeamInfo) error
	AddNBAResult(res models.Result) error
	GetNBAPlayers() ([]models.NBAPlayer, error)
	GetNBATeamInfo() ([]models.NBATeam, error)
	UpdateNBAPlayer(player models.NBAPlayer) error 
	AssignNBAPlayer(player models.NBAPlayer) error
	GetNBAStandings() ([]models.NBAStandings, error)
	GetLastResults(count int) ([]models.Result, error)
	UpdateNBAResult(res models.Result) error
	DeleteNBAResult(res models.Result) error
}