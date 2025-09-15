package models


type Block struct{
	ID uint `json:"id"`
	UserID int `json:"user_id"`
	UserName string `json:"user_name"`
}