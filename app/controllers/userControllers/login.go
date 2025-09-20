package user

import (
	apiexception "CONFESSION-WALL-BACKEND/app/apiException"
	"CONFESSION-WALL-BACKEND/app/models"
	userservices "CONFESSION-WALL-BACKEND/app/services/userServices"
	"CONFESSION-WALL-BACKEND/app/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ShowUser struct {
	Username   string `json:"username"`
	
}

func Login(c *gin.Context) {
	//接收参数
	var data LoginData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiexception.AbortWithException(c, apiexception.ParamError, err)
		return
	}

	//获取用户信息和判断用户是否存在
	var user *models.User
	user, err = userservices.GetUserByUsername(data.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			apiexception.AbortWithException(c, apiexception.UserNotFound, err)
		} else {
			apiexception.AbortWithException(c, apiexception.ServerError, err)
		}
		return
	}

	//判断密码是否正确

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		apiexception.AbortWithException(c, apiexception.PasswordError, err)
		return
	}
	result := ShowUser{
		Username:user.Username,
		
	}
	
	utils.JsonSuccessResponse(c, result)

}