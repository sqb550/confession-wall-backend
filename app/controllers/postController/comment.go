package postController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/models"
	"confession-wall-backend/app/services/postService"
	"confession-wall-backend/app/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

type CommentData struct {
	PostID  int    `json:"post_id"`
	Content string `json:"content"`
	ReplyTo int    `json:"reply_to"`
}

func Comment(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID,ok:=val.(float64)
	if !ok{
		apiException.AbortWithException(c,apiException.ServerError,nil)
		return
	}
	userIDInt:=int(userID)
	var data CommentData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	err = postService.Comment(&models.Comment{
		PostID:   data.PostID,
		UserID: userIDInt,
		ReplyTo:  data.ReplyTo,
		Content:  data.Content,
	})
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	err = postService.IncrComments(data.PostID)
	if err != nil {
		fmt.Println(err)
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.JsonSuccessResponse(c, "评论成功")

}
