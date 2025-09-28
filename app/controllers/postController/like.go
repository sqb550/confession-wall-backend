package postController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/utils"

	"github.com/gin-gonic/gin"
)

type LikeData struct {
	PostID int `json:"post_id" binding:"required"`
}

func Like(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt, _ := userID.(int)
	var data LikeData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	isLiked, err := utils.CheckUserLike(userIDInt, data.PostID, c)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	if !isLiked {
		err = utils.LikeHandler(userIDInt, data.PostID, c)
		if err != nil {
			apiException.AbortWithException(c, apiException.LikeError, err)
			return
		}
	} else {
		err = utils.CancelLikeHandler(userIDInt, data.PostID, c)
		if err != nil {
			apiException.AbortWithException(c, apiException.CancelLikeError, err)
			return
		}
	}
	utils.JsonSuccessResponse(c, nil)

}
