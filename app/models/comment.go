package models

type Comment struct{
	ID uint `json:"id"`
	ConfessionID int `json:"confession_id"`
	Content string `json:"content"`

}