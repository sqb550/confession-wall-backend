package postController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/models"
	"confession-wall-backend/app/services/postService"
	"confession-wall-backend/app/services/userService"
	"confession-wall-backend/app/utils"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ReleaseData struct {
	Content       string    `json:"content"`
	Picture       []string  `json:"picture"`
	Anonymous     bool      `json:"anonymous"`
	Invisible     bool      `json:"invisible"`
	ReleaseTime   time.Time `json:"release_time"`
	ReleaseStatus bool      `json:"release_status"`
}

func Release(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID,ok:=val.(float64)
	if !ok{
		apiException.AbortWithException(c,apiException.ServerError,nil)
		return
	}
	userIDInt:=int(userID)
	var data ReleaseData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		fmt.Println(err)
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	user, err := userService.SeekUser(userIDInt)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	
	var releaseTime time.Time
	if data.ReleaseStatus {
		releaseTime = time.Now()
	} else {
		releaseTime = data.ReleaseTime
	}
	var name string
	if data.Anonymous{
		number:=1000+rand.Intn(1000)
		numberStr:=strconv.Itoa(number)
		name="用户"+numberStr
	}else{
		name=user.Name
	}
	postID,err := postService.ReleasePost(&models.Post{
		UserID:        userIDInt,
		Name:name,
		Avatar: user.Avatar,
		Content:       data.Content,
		Invisible:     data.Invisible,
		Anonymous:     data.Anonymous,
		ReleaseTime:   releaseTime,
		UpdatedAt: releaseTime,
		ReleaseStatus: data.ReleaseStatus,
		Likes:         0,
		Views:         0,
	})
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	err=utils.UpdateHot(c,postID,0,0)
	if err!=nil{
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	for _,url:=range data.Picture{
		err=postService.ReleasePicture(&models.Picture{
			URL: url,
			PostID: postID,
		})
		if err!=nil{
			apiException.AbortWithException(c,apiException.ServerError,err)
			return
		}
	}
	
	utils.JsonSuccessResponse(c, nil)

}