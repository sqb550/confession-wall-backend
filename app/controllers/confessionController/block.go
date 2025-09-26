package confessionController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/models"
	"confession-wall-backend/app/services/confessionService"
	"confession-wall-backend/app/utils"

	"github.com/gin-gonic/gin"
)

type BlockData struct {
	BlockID int `json:"block_id"`
}



func Block(c *gin.Context) {
	userID,_:=c.Get("user_id")
	userIDInt, _ := userID.(int)
	var data BlockData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	err = confessionService.Block(&models.Block{
		UserID: userIDInt,
		BlockedID: data.BlockID,
	})

}

func ShowBlock(c *gin.Context) {
	userID,_:=c.Get("user_id")
	userIDInt, _:= userID.(int)
	result, err := confessionService.ShowBlock(userIDInt) //result为post中的结构体数组
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	var BlockList []int
	for _, data := range result {
		BlockList = append(BlockList, data.BlockedID)
	}
	utils.JsonSuccessResponse(c, BlockList)

}
