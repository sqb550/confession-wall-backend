package userController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/models"
	"confession-wall-backend/app/services/userService"
	"confession-wall-backend/app/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ShowUser struct {
	UserID int `json:"user_id"`
	Token    string `json:"token"`
	Avatar   string`json:"avatar"`
}

func Login(c *gin.Context) {
	var data LoginData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	var user *models.User
	user, err = userService.GetUserByUsername(data.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			apiException.AbortWithException(c, apiException.UserNotFound, err)
		} else {
			apiException.AbortWithException(c, apiException.ServerError, err)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		apiException.AbortWithException(c, apiException.PasswordError, err)
		return
	}
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	result := ShowUser{
		UserID: int(user.ID),
		Token:    token,
		Avatar: user.Avatar,
	}

	utils.JsonSuccessResponse(c, result)

}
