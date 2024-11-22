package model

type Team struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Users []User `json:"users#`
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

type Question struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	CategoryId   int    `json:"category_id"`
	Description  string `json:"description"`
	VMID         int    `json:"vmid"`
	Env          string `json:"env"`
	CategoryName string `json:"category_name"`
}

type QuesionRequest struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	CategoryID  int      `json:"category_id,omitempty"`
	Env         string   `json:"env,omitempty"`
	Sshkeys     []string `json:"sshkeys,omitempty"`
	Memory      int      `json"memory,omitempty"`
	CPUs        int      `json"cpu,omitempty"`
	Disk        int      `json"disk,omitempty"`
	IP          string   `json:"ip,omitempty" validate:"cidr"`
	Gateway     string   `json:"gateway,omitempty" validate:"ip"`
}

type Cloudinit struct {
	ContestQuestionsID int
	TeamID             int
	Filename           string
}
