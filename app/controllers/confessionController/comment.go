package confessionController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/models"
	"confession-wall-backend/app/services/confessionService"
	"confession-wall-backend/app/services/userService"
	"confession-wall-backend/app/utils"

	"github.com/gin-gonic/gin"
)

type CommentData struct {
	ConfessionID int    `json:"confession_id"`
	Content      string `json:"content"`
	ReplyID      int    `json:"reply_id"`
}

func Comment(c *gin.Context) {
	userID,_:=c.Get("user_id")
	userIDInt, _ := userID.(int)
	var data CommentData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	result, err := userService.SeekUser(userIDInt)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
	}
	err = confessionService.Comment(&models.Comment{
		ConfessionID: data.ConfessionID,
		Username:     result.Username,
		ReplyID:      data.ReplyID,
		Content:      data.Content,
	})
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
	}
	err=confessionService.IncrComments(data.ConfessionID)
	if err!=nil{
		apiException.AbortWithException(c,apiException.ServerError,err)
	}
	utils.JsonSuccessResponse(c,nil)

}
