package model

type PveapiResponse[T any] struct {
	Data  string `json:"data"`
	Error string `json:"error"`
}

type Question struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	CategoryId   int    `json:"category_id"`
	Description  string `json:"description"`
	VMID         int    `json:"vmid"`
	Env          string `json:"env"`
	Answer       string `json:"answer"`
	CategoryName string `json:"category_name"`
	Point        int    `json:"point"`
}
type Category struct {
	ID   int
	Name string
}
type Template struct {
}
type CreateQuestion struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CategoryID  int      `json:"category_id"`
	Env         string   `json:"env"`
	Sshkeys     []string `json:"sshkeys"`
	Memory      int      `json"memory"`
	CPUs        int      `json"cpu" validate:"required"`
	Disk        int      `json"disk" validate:"required"`
	IP          string   `json:"ip" validate:"cidr"`
	Gateway     string   `json:"gateway" validate:"ip"`
	CategoryId  int      `json:"category_id"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
}

type CreateVM struct {
	Cloneid  int    `json:"cloneid"`
	Name     string `json:"name"`
	Memory   int    `json:"memory"`
	IP       string `json:"ip,omitempty"`
	Gateway  string `json:"gateway"`
	Disk     int    `json:"disk"`
	Cicustom string `json:"cicustom"`
	CPU      int    `json:"cpu"`
}

type CloudinitResponse struct {
	Filename  string   `json:"filename"`
	Hostname  string   `json:"hostname"`
	Sshkeys   []string `json:"sshkeys"`
	Username  string   `json:"username"`
	Password  string   `json:"passwd"`
	SshPwauth string   `json:"ssh_pwauth"`
}

type Point struct {
	ID        int `json:"id"`
	TeamID    int `json:"team_id"`
	QuesionID int `json:"question_id"`
	ContestID int `json:"contest_id"`
	Point     int `json:"point"`
}

type ResponseIPs map[string][]string
