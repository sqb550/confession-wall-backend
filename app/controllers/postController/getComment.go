package postController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/services/postService"
	"confession-wall-backend/app/utils"

	"github.com/gin-gonic/gin"
)

type postComment struct {
	PostID int `json:"post_id"`
}

func GetComment(c *gin.Context) {
	var data postComment
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	comments, err := postService.ShowComments(data.PostID)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.JsonSuccessResponse(c, comments)

}
