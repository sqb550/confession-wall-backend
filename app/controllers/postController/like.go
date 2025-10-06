package postController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LikeData struct {
	PostID int `json:"post_id" form:"post_id"`
}

func Like(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID,ok:=val.(float64)
	if !ok{
		apiException.AbortWithException(c,apiException.ServerError,nil)
		return
	}
	userIDInt:=int(userID)
	var data LikeData
	err := c.ShouldBindQuery(&data)
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
	likesStr, viewsStr, _, err := utils.GetLikeAndViews(int(data.PostID), c)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	likes, _ := strconv.Atoi(likesStr)
	views, _ := strconv.Atoi(viewsStr)
	err=utils.UpdateHot(c,data.PostID,likes,views)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, nil)
		return
	}

	utils.JsonSuccessResponse(c, likes)

}
