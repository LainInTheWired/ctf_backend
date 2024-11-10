package model

type Team struct {
	ID   int
	Name string
}

type Contest struct {
	ID   int
	Name string
}

type ContestTeams struct {
	ContestID int
	TeamID    int
}
