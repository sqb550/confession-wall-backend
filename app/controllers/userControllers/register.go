package user

import (
	apiexception "CONFESSION-WALL-BACKEND/app/apiException"
	"CONFESSION-WALL-BACKEND/app/models"
	"CONFESSION-WALL-BACKEND/app/utils"
	"CONFESSION-WALL-BACKEND/app/services/userServices"
	"unicode"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserData struct {
	Username string `json:"user_name"`
	Name string `json:"name"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var data UserData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiexception.AbortWithException(c, apiexception.ParamError, err)
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
		apiexception.AbortWithException(c, apiexception.UsernameError, err)
	}


	flag = false
	length := len(data.Password)
	if length >= 8 && length <= 16 {
		flag = true
	}

	if !flag {
		apiexception.AbortWithException(c, apiexception.PasswordLengthError, err)
		return
	}


	result, err :=userservices.GetUserByUsername(data.Username)
	if result != nil {
		apiexception.AbortWithException(c, apiexception.UserExist, err)
		return
	} else if err != gorm.ErrRecordNotFound {
		apiexception.AbortWithException(c, apiexception.ServerError, err)
		return
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

	err = userservices.Register(&models.User{
		Username: data.Username,
		Name:     data.Name,
		Password: string(hashPassword),
	})
	if err != nil {
		apiexception.AbortWithException(c, apiexception.ServerError, err)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}