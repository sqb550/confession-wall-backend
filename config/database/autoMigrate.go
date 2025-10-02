package database

import (
	"confession-wall-backend/app/models"

	"gorm.io/gorm"
)

func autoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
		&models.Block{},
		&models.Picture{},
	)
	return err
}
