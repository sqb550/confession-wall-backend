package postController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/models"
	"confession-wall-backend/app/services/postService"
	"confession-wall-backend/app/services/userService"
	"confession-wall-backend/app/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

type CommentData struct {
	PostID  int    `json:"post_id"`
	Content string `json:"content"`
	RepliedID int    `json:"replied_id"`
}
type ShowCommentsData struct{
	Total int64 `json:"total"`
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
	var name string
	if data.RepliedID!=0{
		comment,err:=postService.SeekComment(data.RepliedID)
		if err != nil {
			apiException.AbortWithException(c, apiException.ParamError, err)
			return
		}
		repliedUser,err:=userService.SeekUser(comment.UserID)
			if err!=nil{
				apiException.AbortWithException(c,apiException.ServerError,err)
				return
			}
		name=repliedUser.Name
	}else{
		postData,err:=postService.SeekPost(data.PostID)
		if err!=nil{
			apiException.AbortWithException(c,apiException.ServerError,err)
			return
		}
		name=postData.Name
	}
	err = postService.Comment(&models.Comment{
		PostID:   data.PostID,
		UserID: userIDInt,
		RepliedID:  data.RepliedID,
		RepliedTo: name,
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
	total,err:=postService.CountComments(data.RepliedID)
	if err!=nil{
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	result:=ShowCommentsData{
		Total: total,
	}
	utils.JsonSuccessResponse(c, result)

}
