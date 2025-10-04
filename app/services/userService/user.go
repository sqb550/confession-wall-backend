package userService

import (
	"confession-wall-backend/app/models"
	"confession-wall-backend/config/database"

	"gorm.io/gorm"
)

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := database.DB.Where("username=?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func Register(user *models.User) error {
	result := database.DB.Create(user)
	return result.Error
}

func SeekUser(userID int) (*models.User, error) {
	var data models.User
	result := database.DB.Where("id=?", userID).First(&data)
	if result.Error != nil {
		return nil, result.Error
	}

	return &data, nil

}

func UpdateName(tx *gorm.DB,userID int, name string) error {
	result := tx.Model(&models.User{}).Where("id=?", userID).Update("Name", name)
	return result.Error

}

func Updatepost(tx *gorm.DB,userID int, name string) error {
	result := tx.Model(&models.Post{}).Where("id=?", userID).Where("anonymous=?", false).Update("Name", name)
	return result.Error

}

func UpdatePassword(newPassword string, userID int) error {
	result := database.DB.Model(&models.User{}).Where("id=?", userID).Update("password", newPassword)
	return result.Error
}

func UploadAvatar(userID int, url string) error {
	result := database.DB.Model(&models.User{}).Where("id=?", userID).Update("Avatar", url)
	return result.Error
}
func UpdateAvatar(userID int, url string) error {
	result := database.DB.Model(&models.Post{}).Where("user_id=?", userID).Update("Avatar", url)
	return result.Error
}


