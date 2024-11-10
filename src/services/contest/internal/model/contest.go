package model

import "time"

type Contest struct {
	ID        int
	Name      string
	StartDate time.Time
	EndDate   time.Time
}

type Team struct {
	ID   int
	Name string
}

type ContestsTeam struct {
	ContestID int
	TeamID    int
}

type Points struct {
	ID         int
	TeamID     int
	QuestionID int
	ContestID  int
	Point      int
}
