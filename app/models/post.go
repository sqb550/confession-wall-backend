package models

import "time"

type Post struct {
	ID            uint      `json:"post_id"`
	UserID        int       `json:"user_id"`
	Avatar        string    `json:"avatar"`
	Content       string    `json:"content"`
	Picture       []string  `json:"picture"`
	Likes         int       `json:"likes"`
	Name          string    `json:"name"`
	Views         int       `json:"views"`
	Comments      int       `json:"comments"`
	Invisible     bool      `json:"invisible"`
	Anonymous     bool      `json:"anonymous"`
	ReleaseTime   time.Time `json:"release_time"`
	ReleaseStatus bool      `json:"release_status"`
	UpdatedAt     time.Time `json:"updated_time"`
	DeletedAt     time.Time `json:"deleted_at"`
}
