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

type NBATeamInfo struct {
	ID           int
	Name         string
	Abbreviation string
	Color1       string
	Color2       string
	DarkText     string
}

// User is the users model
type User struct {
	UserID      int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
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
	HomeTeam   int
	HomeScore  int
	AwayScore  int
	AwayTeam   int
	Time       time.Time
	TimeString string
}

// NBAPlayer is the NBA player model
type Player struct {
	PlayerID                       int
	FirstName                      string
	LastName                       string
	PrimaryPosition                *string
	SecondaryPosition              *string
	TeamID                         int
	AssignedPosition               int
	Archetype                      *string
	Height                         *int
	Weight                         *int
	NBATeam                        *string
	Nationality                    *string
	Birthdate                      *string
	Jersey                         string
	Draft                          *string
	ImgID                          string
	RatingsURL                     *string
	Overall                        *int
	AttributesOutsideScoring       *int
	AttributesAthleticism          *int
	AttributesInsideScoring        *int
	AttributesPlaymaking           *int
	AttributesDefending            *int
	AttributesRebounding           *int
	AttributesIntangibles          *int
	AttributesPotential            *int
	AttributesTotalAttributes      *int
	AttributesCloseShot            *int
	AttributesMidRangeShot         *int
	AttributesThreePointShot       *int
	AttributesFreeThrow            *int
	AttributesShotIQ               *int
	AttributesOffensiveConsistency *int
	AttributesSpeed                *int
	AttributesAcceleration         *int
	AttributesStrength             *int
	AttributesVertical             *int
	AttributesStamina              *int
	AttributesHustle               *int
	AttributesOverallDurability    *int
	AttributesLayup                *int
	AttributesStandingDunk         *int
	AttributesDrivingDunk          *int
	AttributesPostHook             *int
	AttributesPostFade             *int
	AttributesPostControl          *int
	AttributesDrawFoul             *int
	AttributesHands                *int
	AttributesPassAccuracy         *int
	AttributesBallHandle           *int
	AttributesSpeedWithBall        *int
	AttributesPassIQ               *int
	AttributesPassVision           *int
	AttributesInteriorDefense      *int
	AttributesPerimeterDefense     *int
	AttributesSteal                *int
	AttributesBlock                *int
	AttributesLateralQuickness     *int
	AttributesHelpDefenseIQ        *int
	AttributesPassPerception       *int
	AttributesDefensiveConsistency *int
	AttributesOffensiveRebound     *int
	AttributesDefensiveRebound     *int
	BronzeBadges                   int
	SilverBadges                   int
	GoldBadges                     int
	HOFBadges                      int
	TotalBadges                    int
	CreatedAt                      time.Time
	UpdatedAt                      time.Time
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
	BronzeUrl string
	SilverUrl string
	GoldUrl   string
	HOFUrl    string
}

// Players badges is the NBA players badges model
type PlayersBadges struct {
	PlayerID  int
	FirstName string
	LastName  string
	BadgeID   int
	Name      string
	Level     string
}

type NBAStandings struct {
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
