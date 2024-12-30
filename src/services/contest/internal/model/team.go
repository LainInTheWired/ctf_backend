package model

type Team struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Users []User `json:"users"`
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
	ID           int                 `json:"id"`
	Name         string              `json:"name"`
	CategoryId   int                 `json:"category_id"`
	Description  string              `json:"description"`
	VMID         int                 `json:"vmid"`
	Env          string              `json:"env"`
	Answer       string              `json:"answer"`
	Point        int                 `json:"point"`
	CategoryName string              `json:"category_name"`
	CurrentPoint int                 `json:"current_point,omitempty"`
	IPs          map[string][]string `json:"ips"`
}

type QuesionRequest struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	CategoryID  int      `json:"category_id,omitempty"`
	Env         string   `json:"env,omitempty"`
	Sshkeys     []string `json:"sshkeys,omitempty"`
	Memory      int      `json:"memory,omitempty"`
	CPUs        int      `json:"cpu,omitempty"`
	Disk        int      `json:"disk,omitempty"`
	IP          string   `json:"ip,omitempty" validate:"cidr"`
	Gateway     string   `json:"gateway,omitempty" validate:"ip"`
	Password    string   `json:"password,omitempty"`
}

type Cloudinit struct {
	QuestionID int                 `json:"question_id"`
	ContestID  int                 `json:"contest_id"`
	TeamID     int                 `json:"team_id"`
	Filename   string              `json:"filename"`
	Access     string              `json:"access"`
	VMID       int                 `json:"vmid"`
	IPs        map[string][]string `json:"ips,omitempty"`
}

type QuesionResponse[T any] struct {
	Data  T      `json:"data"`
	Error string `json:"error"`
}
