package models

type Comment struct{
	ID uint `json:"comment_id"`
	ConfessionID int `json:"confession_id"`
	Content string `json:"content"`
	ReplyID int `json:"reply_id"`
	Username string `json:"username"`

}