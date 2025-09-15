package models

type User struct {
	ID       uint   `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Password string `json:"password"`
	Name     string `json:"name"`
}