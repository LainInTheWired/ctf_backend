package model

type Team struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Users []User `json:"users"`
}

type Contest struct {
	ID   int
	Name string
}

type ContestTeams struct {
	ContestID int
	TeamID    int
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string
}

type TeamUsers struct {
	TeamID int
	UserID int
}

type Question struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	CategoryId   int    `json:"category_id"`
	Description  string `json:"description"`
	VMID         int    `json:"vmid"`
	Env          string `json:"env"`
	CategoryName string `json:"category_name"`
}
