package model

type ResponsePoints struct {
	TeamID int     `json:"team_id"`
	Name   string  `json:"name"`
	Points []Point `json:"points"`
}

type ResponseIPs map[string][]string
