package repository

import "github.com/oskarpedosk/baltijas-kauss/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	UpdateNBATeamInfo(res models.NBATeamInfo) error
	AddNBAResult(res models.Result) error
	DisplayNBAPlayers() ([]models.NBAPlayer, error)
}