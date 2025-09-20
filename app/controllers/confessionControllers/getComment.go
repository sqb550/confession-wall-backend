package confession

import (
	apiexception "CONFESSION-WALL-BACKEND/app/apiException"
	confessionservices "CONFESSION-WALL-BACKEND/app/services/confessionServices"
	"CONFESSION-WALL-BACKEND/app/utils"

	"github.com/gin-gonic/gin"
)

type ConfessionComment struct {
	ConfessionID int `json:"confession_id"`
}

func GetComment(c *gin.Context){
	var data ConfessionComment
	err := c.ShouldBindJSON(&data) 
	if err != nil {
		apiexception.AbortWithException(c, apiexception.ParamError, err)
		return
	}

	comments,err:=confessionservices.ShowComments(data.ConfessionID)
	if err!=nil{
		apiexception.AbortWithException(c,apiexception.ServerError,err)
		return
	}
	utils.JsonSuccessResponse(c,comments)

}