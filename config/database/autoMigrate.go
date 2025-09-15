package database

import (
	"CONFESSION-WALL-BACKEND/app/models"

	"gorm.io/gorm"
)

func autoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Confession{},
		&models.Comment{},
		&models.Block{},
		&models.CommentReplies{},
	)
	return err
}