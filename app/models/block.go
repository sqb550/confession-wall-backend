package models


type Block struct{
	ID uint `json:"id"`
	UserID int `json:"user_id"`
	BlockedID int `json:"blocked_id"`
}