package postController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/services/postService"
	"confession-wall-backend/app/utils"

	"github.com/gin-gonic/gin"
)
type HotPostData struct{
	PostID int `json:"post_id"`
	Content string `json:"content"`
}

func GetHotRank(c *gin.Context) {
	postIDs, err := utils.GetTopHotRank(c)
	if err != nil {
		apiException.AbortWithException(c, apiException.HotRankError, err)
		return
	}
	hotPosts:=[]HotPostData{}
	for _, postID := range postIDs {
		data,err:=postService.SeekPost(postID)
		if err!=nil{
			apiException.AbortWithException(c,apiException.ServerError,err)
			return
		}
		hotPosts=append(hotPosts, HotPostData{
			PostID: postID,
			Content: data.Content,
		})
	}
	utils.JsonSuccessResponse(c, hotPosts)

}
