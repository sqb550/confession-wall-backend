package userservices

import (
	"CONFESSION-WALL-BACKEND/app/models"
	"CONFESSION-WALL-BACKEND/config/database"
)

func GetUserByUsername(Username string) (*models.User, error) {
	var user models.User
	result := database.DB.Where("username=?", Username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func Register(user *models.User) error {
	result := database.DB.Create(user)
	return result.Error
}

func SeekUser(UserID int) (*models.User, error) {
	var data models.User
	result := database.DB.Where("ID=?", UserID).First(&data)
	if result.Error != nil {
		return nil, result.Error
	}

	return &data, nil

}

func UpdateName(UserID int,Name string)error{
	result:=database.DB.Model(&models.User{}).Where("user_id=?",UserID).Update("Name",Name)
	return result.Error

}

func UpdateConfession(UserID int,Name string)error{
	result:=database.DB.Model(&models.Confession{}).Where("user_id=?",UserID).Where("anonymous=?",false).Update("Name",Name)
	return result.Error

}

func UpdatePassword(NewPassword string,UserID int)error{
	result:=database.DB.Where("user_id=?",UserID).Update("password",NewPassword)
	return result.Error
}