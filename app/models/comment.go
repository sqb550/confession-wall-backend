package models

type Comment struct {
	ID       uint   `json:"comment_id"`
	PostID   int    `json:"post_id"`
	Content  string `json:"content"`
	ReplyID  int    `json:"reply_id"`
	Name string `json:"name"`
	Avatar string `json:"avatar"`
}
