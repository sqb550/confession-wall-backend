package postController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/services/postService"
	"confession-wall-backend/app/services/userService"
	"confession-wall-backend/app/utils"

	"github.com/gin-gonic/gin"
)

type postComment struct {
	PostID int `json:"post_id" form:"post_id"`
}
type GetCommentsData struct{
	ID uint `json:"id"`
	Content string `json:"content"`
	ReplyTo int `json:"reply_to"`
	Avatar string `json:"avatar"`
	Author string `json:"author"`
}

func GetComment(c *gin.Context) {
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
	blocked,err:=postService.ShowBlock(userIDInt)
	if err!=nil{
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	blockIDs:=make([]int,0)
	for _,block:=range blocked{
		blockIDs=append(blockIDs, block.BlockedID)

	}

	result, err := postService.ShowComments(data.PostID,blockIDs)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	var comments []GetCommentsData
	for _,data:=range result{
		user,err:=userService.SeekUser(data.UserID)
		if err!=nil{
			apiException.AbortWithException(c,apiException.ServerError,err)
			return
		}
		comments=append(comments,GetCommentsData{
			Content: data.Content,
			ID: data.ID,
			ReplyTo: data.ReplyTo,
			Author: user.Name,
			Avatar: user.Avatar,
		})
	}
	utils.JsonSuccessResponse(c, comments)

}
