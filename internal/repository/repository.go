package repository

import (
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
)

type DatabaseRepo interface {
	// Users
	AllUsers() bool
	GetUser(userID int) (models.User, error)
	Authenticate(email, password string) (int, string, int, error)
	
	ChangePassword(userID int, password string) error
	UpdateUserImage(userID int, img string) error
	UpdateUserInfo(userID int, firstName, lastName, email string) error

	// Teams
	GetTeam(teamID int) (models.Team, error)
	GetTeams() ([]models.Team, error)
	UpdateTeam(team models.Team) error

	// Players
	AddPlayer(playerID, teamID int) error
	DropPlayer(playerID int) error
	SwitchTeam(player models.Player) error
	AssignPosition(playerID, position int) error

	GetBadgeID(url string) (int, error)
	CreateNewBadge(models.Badge) (int, error)
	GetPlayerBadges(playerID int) ([]models.Badge, error)

	ResetPlayers() error
	CountPlayers() (int, error)
	GetADP(playerID int) (float64, error)
	UpdatePlayer(player models.Player) error
	UpdatePlayerBadges(models.Player, []models.Badge) error
	GetPlayer(playerID int) (models.Player, error)
	GetPlayers(filter models.Filter) ([]models.Player, error)
	GetTeamPlayers(teamID int) ([]models.Player, error)

	// Draft
	GetDraftID() (int, error)
	GetDrafts() ([]models.DraftPick, error)
	GetDraft(draftID int) ([]models.DraftPick, error)
	AddDraftPick(draftID int, draftPick models.DraftPick) error
	SelectRandomPlayer(random int) (models.Player, error)

	// Standings
	AddResult(res models.Result) error
	UpdateResult(res models.Result) error
	DeleteResult(res models.Result) error

	StartNewSeason() error
	GetSeasons() ([]models.Season, error)
	GetSeasonResults(seasonID int) ([]models.Result, error)
	GetAllResults() ([]models.Result, error)
}
