package model

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Role     []Role
}

type Role struct {
	ID         int
	Name       string
	Namespace  string
	Permission []Permission
}
type Permission struct {
	ID          int
	Name        string
	Description string
}
