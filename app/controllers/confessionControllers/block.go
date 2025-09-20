package confession

import (
	apiexception "CONFESSION-WALL-BACKEND/app/apiException"
	"CONFESSION-WALL-BACKEND/app/models"
	confessionservices "CONFESSION-WALL-BACKEND/app/services/confessionServices"
	"CONFESSION-WALL-BACKEND/app/utils"

	"github.com/gin-gonic/gin"
)

type BlockData struct {
	BlockID int `json:"block_id"`
	UserID  int `json:"user_id"`
}

type ShowBlockData struct{
	UserID int `json:"user_id"`
}

func Block(c *gin.Context) {
	var data BlockData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiexception.AbortWithException(c, apiexception.ParamError, err)
		return
	}
	err=confessionservices.Block(&models.Block{
		UserID: data.UserID,
		BlockedID: data.BlockID,
	})

}

func ShowBlock(c *gin.Context){
	var data ShowBlockData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiexception.AbortWithException(c, apiexception.ParamError, err)
		return
	}
	result, err := confessionservices.ShowBlock(data.UserID) //result为post中的结构体数组
	if err != nil {
		apiexception.AbortWithException(c, apiexception.ServerError, err)
	}

	var BlockList []int
	for _, data := range result {
		BlockList = append(BlockList,data.BlockedID)
	}
	utils.JsonSuccessResponse(c,BlockList)


}