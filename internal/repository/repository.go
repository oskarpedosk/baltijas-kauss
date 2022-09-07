package repository

import "github.com/oskarpedosk/baltijas-kauss/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	UpdateTeamInfo(res models.NBATeamInfo) error
}