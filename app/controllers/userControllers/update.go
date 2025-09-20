package user

import (
	apiexception "CONFESSION-WALL-BACKEND/app/apiException"
	userservices "CONFESSION-WALL-BACKEND/app/services/userServices"
	"CONFESSION-WALL-BACKEND/app/utils"
	"CONFESSION-WALL-BACKEND/config/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)


type UpdateNameData struct{
	Name string `json:"name"`
	UserID int `json:"user_id"`
}

type PasswordData struct{
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	UserID int `json:"user_id"`
}

func UpdateName(c *gin.Context){
	var data UpdateNameData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiexception.AbortWithException(c, apiexception.ParamError, err)
		return
	}
	tx := database.DB.Begin()
	if tx.Error != nil {
		apiexception.AbortWithException(c, apiexception.ServerError, tx.Error)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err=userservices.UpdateName(data.UserID,data.Name)
	if err!=nil{
		tx.Rollback()
		apiexception.AbortWithException(c,apiexception.ServerError,err)
		return
	}
	err=userservices.UpdateConfession(data.UserID,data.Name)
	if err!=nil{
		tx.Rollback()
		apiexception.AbortWithException(c,apiexception.ServerError,err)
		return
	}
	_ = tx.Commit()

	utils.JsonSuccessResponse(c, nil)
}
func UpdatePassword(c *gin.Context){
	var data PasswordData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiexception.AbortWithException(c, apiexception.ParamError, err)
		return
	}
	user,err:=userservices.SeekUser(data.UserID)
	if err!=nil{
		apiexception.AbortWithException(c,apiexception.ServerError,err)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.OldPassword))
	if err != nil {
		apiexception.AbortWithException(c, apiexception.PasswordError, err)
		return
	}
	hashNewPassword, _ := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
	err=userservices.UpdatePassword(string(hashNewPassword),data.UserID)
	if err!=nil{
		apiexception.AbortWithException(c,apiexception.ServerError,err)
	}
	utils.JsonSuccessResponse(c,nil)

	
}