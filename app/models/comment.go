package models

import "gorm.io/gorm"

type Comment struct {
	ID        uint           `json:"comment_id"`
	PostID    int            `json:"post_id"`
	Content   string         `json:"content"`
	RepliedID   int            `json:"replied_id"`
	RepliedTo  string`json:"replied_to"`
	UserID    int            `json:"user_id"`
	DeletedAt gorm.DeletedAt `json:"deleted_time"`
}
