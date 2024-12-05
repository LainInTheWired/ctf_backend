package model

import "time"

type Contest struct {
	ID        int        `json:"id,"`
	Name      string     `json:"name"`
	StartDate time.Time  `json:"start_date"`
	EndDate   time.Time  `json:"end_date"`
	Questions []Question `json:"questions"`
}

type ContestsTeam struct {
	ContestID int
	TeamID    int
}

type Point struct {
	ID           int       `json:"id,omitempty"`
	TeamID       int       `json:"team_id,omitempty"`
	QuestionID   int       `json:"question_id,omitempty"`
	ContestID    int       `json:"contest_id,omitempty"`
	InsertDate   time.Time `json:"insert_date,omitempty"`
	Point        int       `json:"point"`
	CurrentPoint int       `json:"current_point,omitempty"`
}
