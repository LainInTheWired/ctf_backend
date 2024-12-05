package model

type User struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Role     []Role `json:"role,omitempty"`
}

type Role struct {
	ID         int          `json:"id,omitempty"`
	Name       string       `json:"name,omitempty"`
	Namespace  string       `json:"namespace,omitempty"`
	Contest    []string     `json:"contest,omitempty"`
	Permission []Permission `json:"permission,omitempty"`
}
type Permission struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}
