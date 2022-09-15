package repository

import "github.com/oskarpedosk/baltijas-kauss/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	UpdateNBATeamInfo(team models.NBATeamInfo) error
	AddNBAResult(res models.Result) error
	GetNBAPlayers() ([]models.NBAPlayer, error)
	GetNBATeamInfo() ([]models.NBATeam, error)
	UpdateNBAPlayer(player models.NBAPlayer) error 
}