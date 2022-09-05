package models

type TeamInfo struct {
	TeamName     string
	Abbreviation string
	TeamColor    string
	DarkText     string
}

type Result struct {
	HomeTeam  string
	HomeScore int
	AwayScore int
	AwayTeam  string
}
