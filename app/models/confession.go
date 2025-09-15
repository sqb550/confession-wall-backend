package models

type Confession struct{
	ID uint `json:"confession_id"`
	UserID int `json:"user_id"`
	Content string `json:"content"`
	Picture string `json:"picture"`
	Likes int `json:"likes"`
	Name string `json:"name"`
}