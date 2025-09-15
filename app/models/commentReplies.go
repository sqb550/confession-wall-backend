package models
type CommentReplies struct{
	ID uint `json:"id"`
	CommentID int `json:"comment_id"`
	Content string `json:"content"`
}