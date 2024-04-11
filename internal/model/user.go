package model

type User struct {
	ID       int    `json:"-" db:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
