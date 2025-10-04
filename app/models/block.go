package models

import "time"

type Block struct {
	ID        uint `json:"id"`
	UserID    int  `json:"user_id"`
	BlockedID int  `json:"blocked_id"`
	DeletedAt time.Time `json:"deleted_at"`
}