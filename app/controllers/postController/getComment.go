package postController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/services/postService"
	"confession-wall-backend/app/services/userService"
	"confession-wall-backend/app/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type postComment struct {
	PostID int `json:"post_id" form:"post_id"`
}
type QueryCommentsData struct{
	ID uint `json:"id"`
	Content string `json:"content"`
	ReplyTo string `json:"reply_to"`
	Avatar string `json:"avatar"`
	Author string `json:"author"`
}

func QueryComment(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID,ok:=val.(float64)
	if !ok{
		apiException.AbortWithException(c,apiException.ServerError,nil)
		return
	}
	userIDInt:=int(userID)
	var data postComment
	err := c.ShouldBindQuery(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	blocked,err:=postService.QueryBlock(userIDInt)
	if err!=nil{
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	blockIDs:=make([]int,0)
	for _,block:=range blocked{
		blockIDs=append(blockIDs, block.BlockedID)

	}
	result, err := postService.QueryComments(data.PostID,blockIDs)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	var comments []QueryCommentsData
	for _,data:=range result{
		user,err:=userService.SeekUser(data.UserID)
		if err!=nil{
			apiException.AbortWithException(c,apiException.ServerError,err)
			return
		}
		comments=append(comments,QueryCommentsData{
			Content: data.Content,
			ID: data.ID,
			ReplyTo: data.RepliedTo,
			Author: user.Name,
			Avatar: user.Avatar,
		})
	}
	err = utils.IncrViewCount(data.PostID, c)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, nil)
		return
	}
	likesStr, viewsStr, _, err := utils.GetLikeAndViews(data.PostID, c)
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

	utils.JsonSuccessResponse(c, comments)

}
