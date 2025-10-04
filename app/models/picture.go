package models

import "gorm.io/gorm"

type Picture struct {
	URL       string `json:"url"`
	PostID    int    `json:"post_id"`
	DeletedAt gorm.DeletedAt `json:"deleted_time"`

}