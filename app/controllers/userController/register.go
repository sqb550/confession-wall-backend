package user

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/models"
	"confession-wall-backend/app/services/userService"
	"confession-wall-backend/app/utils"

	"unicode"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserData struct {
	Username string `json:"user_name"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var data UserData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	flag := true
	if data.Username == "" {
		flag = false
	}
	for _, r := range data.Username {
		if !unicode.IsDigit(r) {
			flag = false
		}
	}
	if !flag {
		apiException.AbortWithException(c, apiException.UsernameError, nil)
	}

	flag = false
	length := len(data.Password)
	if length >= 8 && length <= 16 {
		flag = true
	}

	if !flag {
		apiException.AbortWithException(c, apiException.PasswordLengthError, nil)
		return
	}

	result, err := userService.GetUserByUsername(data.Username)
	if result != nil {
		apiException.AbortWithException(c, apiException.UserExist, err)
		return
	} else if err != gorm.ErrRecordNotFound {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	hashPassword, err:= bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err!=nil{
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	err = userService.Register(&models.User{
		Username: data.Username,
		Name:     data.Name,
		Password: string(hashPassword),
	})
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}
