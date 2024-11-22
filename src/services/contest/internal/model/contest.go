package model

import "time"

type Contest struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
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
