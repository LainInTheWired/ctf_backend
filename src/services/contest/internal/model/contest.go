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
	ID         int       `json:"id,omitempty"`
	TeamID     int       `json:"team_id,omitempty"`
	QuestionID int       `json:"question_id,omitempty"`
	ContestID  int       `json:"contest_id,omitempty"`
	InsertDate time.Time `json:"insert_date,omitempty"`
	Point      int       `json:"point"`
}

type ContestQuestions struct {
	ContestID  int
	QuestionID int
	Point      int
}

type ClusterResources struct {
	Uptime    int     `json:"uptime"`
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Maxmem    int64   `json:"maxmem"`
	Node      string  `json:"node"`
	Status    string  `json:"status"`
	Maxcpu    int     `json:"maxcpu"`
	Netin     int     `json:"netin"`
	Mem       int     `json:"mem"`
	Template  int     `json:"template"`
	Diskread  int     `json:"diskread"`
	Type      string  `json:"type"`
	Diskwrite int     `json:"diskwrite"`
	Maxdisk   int64   `json:"maxdisk"`
	CPU       float64 `json:"cpu"`
	Disk      int     `json:"disk"`
	Netout    int     `json:"netout"`
	Vmid      int     `json:"vmid"`
}
