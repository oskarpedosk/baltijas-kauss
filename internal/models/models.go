package models

type NBATeamInfo struct {
	ID           int
	Name         string
	Abbreviation string
	Color        string
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

// NBATeam is the NBA teams model
type NBATeam struct {
	TeamID       int
	Name         string
	Abbreviation string
	Color        string
	Text         string
	OwnerID      int
	Player       NBAPlayer
}

// NBAResults is the NBA results model
type NBAResults struct {
	HomeTeam  int
	HomeScore int
	AwayScore int
	AwayTeam  int
}

type Result struct {
	HomeTeam  int
	HomeScore int
	AwayScore int
	AwayTeam  int
}

// NBAPlayer is the NBA player model
type NBAPlayer struct {
	PlayerID                  int
	FirstName                 string
	LastName                  string
	PrimaryPosition           string
	SecondaryPosition         string
	Archetype                 string
	NBATeam                   string
	Height                    int
	Weight                    int
	ImgUrl                    string
	PlayerUrl                 string
	TeamID                    int
	StatsOverall              int
	StatsOutsideScoring       int
	StatsAtheliticism         int
	StatsInsideScoring        int
	StatsPlaymaking           int
	StatsDefending            int
	StatsRebounding           int
	StatsCloseShot            int
	StatsMidRangeShot         int
	StatsThreePointShot       int
	StatsFreeThrow            int
	StatsShotIQ               int
	StatsOffensiveConsistency int
	StatsSpeed                int
	StatsAcceleration         int
	StatsStrength             int
	StatsVertical             int
	StatsStamina              int
	StatsHustle               int
	StatsOverallDurability    int
	StatsLayup                int
	StatsStandingDunk         int
	StatsDrivingDunk          int
	StatsPostHook             int
	StatsPostFade             int
	StatsPostControl          int
	StatsDrawFoul             int
	StatsHands                int
	StatsPassAccuracy         int
	StatsBallHandle           int
	StatsSpeedWithBall        int
	StatsPassIQ               int
	StatsPassVision           int
	StatsInteriorDefense      int
	StatsPerimeterDefense     int
	StatsSteal                int
	StatsBlock                int
	StatsLateralQuickness     int
	StatsHelpDefenseIQ        int
	StatsPassPerception       int
	StatsDefensiveConsistency int
	StatsOffensiveRebound     int
	StatsDefensiveRebound     int
	StatsIntangibles          int
	StatsPotential            int
	StatsTotalAttributes      int
	BronzeBadges              int
	SilverBadges              int
	GoldBadges                int
	HOFBadges                 int
	TotalBadges               int
	Assigned                  bool
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
	TotalWins      int
	TotalLosses    int
	HomeWins       int
	HomeLosses     int
	RoadWins       int
	RoadLosses     int
	LastTen        string
	Streak         int
	BasketsFor     int
	BasketsAgainst int
}
