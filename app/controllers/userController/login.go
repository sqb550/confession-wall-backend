package user

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
	Username string `json:"username"`
	Token    string `json:"token"`
}

func Login(c *gin.Context) {
	var data LoginData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	//获取用户信息和判断用户是否存在
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

	//判断密码是否正确

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		apiException.AbortWithException(c, apiException.PasswordError, err)
		return
	}
	user,err=userService.GetUserByUsername(data.Username)
	if err!=nil{
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	token,err:=utils.GenerateToken(user.ID)
	if err!=nil{
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	result := ShowUser{
		Username: user.Username,
		Token: token,
	}

	utils.JsonSuccessResponse(c, result)

}
