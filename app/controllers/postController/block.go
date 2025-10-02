package postController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/models"
	"confession-wall-backend/app/services/postService"
	"confession-wall-backend/app/services/userService"
	"confession-wall-backend/app/utils"

	"github.com/gin-gonic/gin"
)

type BlockData struct {
	BlockID int `json:"block_id"`
}
type GetBlockedData struct{
	BlockID int `json:"block_id"`
	BlockName string `json:"block_name"`
	Avatar   string `json:"avatar"`
}

func Block(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID,ok:=val.(float64)
	if !ok{
		apiException.AbortWithException(c,apiException.ServerError,nil)
		return
	}
	userIDInt:=int(userID)
	var data BlockData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	err = postService.Block(&models.Block{
		UserID:    userIDInt,
		BlockedID: data.BlockID,
	})
	if err!=nil{
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	utils.JsonSuccessResponse(c,"拉黑成功")
}

func ShowBlock(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID,ok:=val.(float64)
	if !ok{
		apiException.AbortWithException(c,apiException.ServerError,nil)
		return
	}
	userIDInt:=int(userID)
	result, err := postService.ShowBlock(userIDInt) //result为post中的结构体数组
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	var BlockList []GetBlockedData
	for _, data := range result {
		blockedData,err:=userService.SeekUser(data.BlockedID)
		if err!=nil{
			apiException.AbortWithException(c,apiException.ServerError,err)
			return
		}
		BlockList = append(BlockList, GetBlockedData{
			BlockID: data.BlockedID,
			BlockName: blockedData.Username,
			Avatar: blockedData.Avatar,
		})
	}
	utils.JsonSuccessResponse(c, BlockList)

}
