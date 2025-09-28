package postController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/utils"

	"github.com/gin-gonic/gin"
)

func GetHotRank(c *gin.Context) {
	postIDs, err := utils.GetTopHotRank(c)
	if err != nil {
		apiException.AbortWithException(c, apiException.HotRankError, err)
		return
	}
	m := make(map[int]float64, 10)
	for _, postID := range postIDs {
		hot, err := utils.GetHot(c, postID)
		if err != nil {
			apiException.AbortWithException(c, apiException.ServerError, err)
			return
		}
		m[postID] = hot
	}
	utils.JsonSuccessResponse(c, m)

}
