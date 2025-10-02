package userController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/services/userService"
	"confession-wall-backend/app/utils"
	"confession-wall-backend/config/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UpdateNameData struct {
	Name string `json:"name"`
}

type PasswordData struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func UpdateName(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID,ok:=val.(float64)
	if !ok{
		apiException.AbortWithException(c,apiException.ServerError,nil)
		return
	}
	userIDInt:=int(userID)
	var data UpdateNameData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	tx := database.DB.Begin()
	if tx.Error != nil {
		apiException.AbortWithException(c, apiException.ServerError, tx.Error)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err = userService.UpdateName(userIDInt, data.Name)
	if err != nil{
		tx.Rollback()
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	err = userService.Updatepost(userIDInt, data.Name)
	if err != nil {
		tx.Rollback()
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	_ = tx.Commit()
	utils.JsonSuccessResponse(c, nil)
}
func UpdatePassword(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID,ok:=val.(float64)
	if !ok{
		apiException.AbortWithException(c,apiException.ServerError,nil)
		return
	}
	userIDInt:=int(userID)
	var data PasswordData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	user, err := userService.SeekUser(userIDInt)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.OldPassword))
	if err != nil {
		apiException.AbortWithException(c, apiException.PasswordError, err)
		return
	}
	hashNewPassword, err:= bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	err = userService.UpdatePassword(string(hashNewPassword), userIDInt)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.JsonSuccessResponse(c, nil)

}
