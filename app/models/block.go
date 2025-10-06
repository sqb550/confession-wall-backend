package models

import (

	"gorm.io/gorm"
)

type Block struct {
	ID        uint `json:"id"`
	UserID    int  `json:"user_id"`
	BlockedID int  `json:"blocked_id"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}