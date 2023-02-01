package models

import (
	"time"
)

type PaginationData struct {
	NextPage     int
	PreviousPage int
	CurrentPage  int
	TotalPages   int
	TwoBefore    int
	TwoAfter     int
	ThreeAfter   int
	Offset       int
	BaseURL      string
}

// User is the users model
type User struct {
	UserID      int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	ImgID       string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Team struct {
	TeamID       int
	Name         string
	Abbreviation string
	Color1       string
	Color2       string
	DarkText     string
	UserID       int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Result struct {
	ResultID   int
	Season     int
	HomeTeamID int
	HomeScore  int
	AwayScore  int
	AwayTeamID int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type PlayerWithTeamInfo struct {
	Player Player
	Team   Team
}

// NBAPlayer is the NBA player model
type Player struct {
	PlayerID          int        `json:"player_id"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	PrimaryPosition   string     `json:"primary_position"`
	SecondaryPosition string     `json:"secondary_position"`
	TeamID            int        `json:"team_id"`
	AssignedPosition  int        `json:"assigned_position"`
	Archetype         string     `json:"archetype"`
	Height            *int       `json:"height"`
	Weight            *int       `json:"weight"`
	NBATeam           string     `json:"nba_team"`
	Nationality       string     `json:"nationality"`
	Birthdate         string     `json:"birthdate"`
	Jersey            string     `json:"jersey"`
	Draft             string     `json:"draft"`
	ImgURL            string     `json:"img_url"`
	RatingsURL        string     `json:"ratings_url"`
	Overall           int        `json:"overall"`
	Attributes        Attributes `json:"attributes"`
	BronzeBadges      int        `json:"bronze_badges"`
	SilverBadges      int        `json:"silver_badges"`
	GoldBadges        int        `json:"gold_badges"`
	HOFBadges         int        `json:"hof_badhes"`
	TotalBadges       int        `json:"total_badges"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Attributes struct {
	OutsideScoring       int
	Athleticism          int
	InsideScoring        int
	Playmaking           int
	Defending            int
	Rebounding           int
	Intangibles          int
	Potential            int
	TotalAttributes      int
	CloseShot            int
	MidRangeShot         int
	ThreePointShot       int
	FreeThrow            int
	ShotIQ               int
	OffensiveConsistency int
	Speed                int
	Acceleration         int
	Strength             int
	Vertical             int
	Stamina              int
	Hustle               int
	OverallDurability    int
	Layup                int
	StandingDunk         int
	DrivingDunk          int
	PostHook             int
	PostFade             int
	PostControl          int
	DrawFoul             int
	Hands                int
	PassAccuracy         int
	BallHandle           int
	SpeedWithBall        int
	PassIQ               int
	PassVision           int
	InteriorDefense      int
	PerimeterDefense     int
	Steal                int
	Block                int
	LateralQuickness     int
	HelpDefenseIQ        int
	PassPerception       int
	DefensiveConsistency int
	OffensiveRebound     int
	DefensiveRebound     int
}

type NBAPosition struct {
	Name   string
	Number int
}

// Badge is the NBA badge model
type Badge struct {
	BadgeID   int
	Name      string
	Type      string
	Info      string
	ImgID     string
	URL       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Players badges is the NBA players badges model
type PlayersBadges struct {
	PlayerID  int
	BadgeID   int
	FirstName string
	LastName  string
	Name      string
	Level     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Standings struct {
	TeamID         int
	WinPercentage  int
	Played         int
	TotalWins      int
	TotalLosses    int
	HomeWins       int
	HomeLosses     int
	AwayWins       int
	AwayLosses     int
	Streak         string
	StreakCount    int
	BasketsFor     int
	BasketsAgainst int
	BasketsSum     int
	ForAvg         float64
	AgainstAvg     float64
	LastFive       []string
}
