package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID            uint      `json:"post_id"`
	UserID        int       `json:"user_id"`
	Name string `json:"name"`
	Avatar string `json:"avatar"`
	Content       string    `json:"content"`
	Likes         int       `json:"likes"`
	Views         int       `json:"views"`
	Comments      int       `json:"comments"`
	Invisible     bool      `json:"invisible"`
	Anonymous     bool      `json:"anonymous"`
	ReleaseTime   time.Time `json:"release_time"`
	ReleaseStatus bool      `json:"release_status"`
	UpdatedAt     time.Time `json:"updated_time"`
	DeletedAt    gorm.DeletedAt `json:"deleted_time"`
}
