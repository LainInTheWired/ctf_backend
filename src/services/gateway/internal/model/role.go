package model

type Permission struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Role struct {
	ID         int          `json:"id"`
	Name       string       `json:"name"`
	Namespace  string       `json:"namespace"`
	Contest    []string     `json:"contest"`
	Permission []Permission `json:"permission"`
}

type User struct {
	ID int `json:"id"`
}
