package confessionController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/services/confessionService"
	"confession-wall-backend/app/utils"

	"github.com/gin-gonic/gin"
)

type ConfessionComment struct {
	ConfessionID int `json:"confession_id"`
}

func GetComment(c *gin.Context) {
	var data ConfessionComment
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	comments, err := confessionService.ShowComments(data.ConfessionID)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.JsonSuccessResponse(c, comments)

}
