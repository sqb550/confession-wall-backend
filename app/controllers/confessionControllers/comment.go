package confession

import (
	apiexception "CONFESSION-WALL-BACKEND/app/apiException"
	"CONFESSION-WALL-BACKEND/app/models"
	confessionservices "CONFESSION-WALL-BACKEND/app/services/confessionServices"
	userservices "CONFESSION-WALL-BACKEND/app/services/userServices"

	"github.com/gin-gonic/gin"
)

type CommentData struct {
	ConfessionID int    `json:"confession_id"`
	Content      string `json:"content"`
	ReplyID      int    `json:"reply_id"`
	UserID       int    `json:"user_id"`
}

func Comment(c *gin.Context){
	var data CommentData
	err := c.ShouldBindJSON(&data) 
	if err != nil {
		apiexception.AbortWithException(c, apiexception.ParamError, err)
		return
	}
	result,err:=userservices.SeekUser(data.UserID)
	if err !=nil{
		apiexception.AbortWithException(c,apiexception.ServerError,err)
	}
	err=confessionservices.Comment(&models.Comment{
		ConfessionID: data.ConfessionID,
		Username: result.Username,
		ReplyID: data.ReplyID,
		Content: data.Content,	
	})
	if err !=nil{
		apiexception.AbortWithException(c,apiexception.ServerError,err)
	}




}