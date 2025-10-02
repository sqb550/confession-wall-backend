package userService

import (
	"confession-wall-backend/app/models"
	"confession-wall-backend/config/database"
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

func UpdateName(userID int, name string) error {
	result := database.DB.Model(&models.User{}).Where("id=?", userID).Update("Name", name)
	return result.Error

}

func Updatepost(userID int, name string) error {
	result := database.DB.Model(&models.Post{}).Where("id=?", userID).Where("anonymous=?", 0).Update("Name", name)
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



