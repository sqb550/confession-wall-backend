package postController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/models"
	"confession-wall-backend/app/services/postService"
	"confession-wall-backend/app/services/userService"
	"confession-wall-backend/app/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReleaseData struct {
	Content       string    `json:"content" binding:"required"`
	Picture       []string  `json:"picture"`
	Anonymous     bool      `json:"anonymous" binding:"required"`
	Invisible     bool      `json:"invisible" binding:"required"`
	ReleaseTime   time.Time `json:"release_time" binding:"required"`
	ReleaseStatus bool      `json:"release_status" binding:"required"`
}

func Release(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt, _ := userID.(int)
	var data ReleaseData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	user, err := userService.SeekUser(userIDInt)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	var name string
	if data.Anonymous {
		name = uuid.New().String()

	} else {
		name = user.Name
	}
	var releaseTime time.Time
	if data.ReleaseStatus {
		releaseTime = time.Now()
	} else {
		releaseTime = data.ReleaseTime
	}

	err = postService.Releasepost(&models.Post{
		UserID:        userIDInt,
		Avatar:        user.Avatar,
		Content:       data.Content,
		Picture:       data.Picture,
		Invisible:     data.Invisible,
		Anonymous:     data.Anonymous,
		Name:          name,
		ReleaseTime:   releaseTime,
		ReleaseStatus: true,
		Likes:         0,
		Views:         0,
	})
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.JsonSuccessResponse(c, nil)

}