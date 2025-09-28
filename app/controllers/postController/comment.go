package postController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/models"
	"confession-wall-backend/app/services/postService"
	"confession-wall-backend/app/services/userService"
	"confession-wall-backend/app/utils"

	"github.com/gin-gonic/gin"
)

type CommentData struct {
	PostID  int    `json:"post_id"`
	Content string `json:"content"`
	ReplyID int    `json:"reply_id"`
}

func Comment(c *gin.Context) {
	userID, _ := c.Get("user_id")
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
	err = postService.Comment(&models.Comment{
		PostID:   data.PostID,
		Name: result.Name,
		ReplyID:  data.ReplyID,
		Content:  data.Content,
		Avatar: result.Avatar,
	})
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
	}
	err = postService.IncrComments(data.PostID)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
	}
	utils.JsonSuccessResponse(c, nil)

}
