package confessionController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/models"
	"confession-wall-backend/app/services/confessionService"
	"confession-wall-backend/app/services/userService"
	"strconv"

	"confession-wall-backend/app/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReleaseData struct {
	Content   string `json:"content" binding:"required"`
	Picture   []string `json:"picture"`
	Anonymous bool   `json:"anonymous" binding:"required"`
	Invisible bool   `json:"invisible" binding:"required"`
}

type PageQuery struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

type ShowConfessionData struct {
	Name      string    `json:"name"`
	Likes     int       `json:"likes"`
	Comments  int       `json:"comments"`
	Content   string    `json:"content"`
	Views     int       `json:"views"`
	Avatar    string    `json:"avatar"`
	Picture   []string    `json:"pictrue"`
	CreatedAt time.Time `json:"created_time"`
	UpdatedAt time.Time `json:"updated_time"`
}

type DeleteData struct {
	ConfessionID int `form:"confession_id" json:"id" binding:"required"`
}

type UpdateData struct {
	ID      uint   `json:"confession_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func Release(c *gin.Context) {
	userID,_:=c.Get("user_id")
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
		name = uuid.NewString()

	} else {
		name = user.Name
	}

	err = confessionService.ReleaseConfession(&models.Confession{
		UserID:    userIDInt,
		Avatar:    user.Avatar,
		Content:   data.Content,
		Picture:   data.Picture,
		Invisible: data.Invisible,
		Anonymous: data.Anonymous,
		Name:      name,
		Likes:     0,
		Views:     0,
	})
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	utils.JsonSuccessResponse(c, nil)

}

func ShowConfessions(c *gin.Context) {
	var data PageQuery
	err := c.ShouldBindQuery(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
	}
	offset := (data.Page - 1) * data.PageSize                              //计算偏移量
	result, err := confessionService.ShowConfession(offset, data.PageSize) //result为post中的结构体数组
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
	}

	var NewConfession []ShowConfessionData
	for _, data := range result {
		likesStr,viewsStr,flag,err:=utils.GetLikeAndViews(int(data.ID),c)
		likes:=0
		views:=0
		if err!=nil{
			apiException.AbortWithException(c,apiException.ServerError,err)
			return
		}
		if flag{
			likes,err=strconv.Atoi(likesStr)
			if err!=nil{
				apiException.AbortWithException(c,apiException.ServerError,err)
				return
			}
			views,err=strconv.Atoi(viewsStr)
			if err!=nil{
				apiException.AbortWithException(c,apiException.ServerError,err)
				return
			}
		}
		NewConfession = append(NewConfession,
			ShowConfessionData{
				Name:      data.Name,
				Likes:     likes,
				Comments:  data.Comments,
				Views:     views,
				Content:   data.Content,
				Avatar:    data.Avatar,
				Picture:   data.Picture,
				CreatedAt: data.CreatedAt,
				UpdatedAt: data.UpdatedAt,
			})
		err=utils.IncrViewCount(int(data.ID),c)
		if err!=nil{
			apiException.AbortWithException(c,apiException.ServerError,nil)
			return
		}
	}

	utils.JsonSuccessResponse(c, NewConfession)
}

func ShowMyConfessions(c *gin.Context) {
	userID,_:=c.Get("user_id")
	userIDInt, _ := userID.(int)
	result, err := confessionService.ShowMyConfession(userIDInt) 
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
	}
	utils.JsonSuccessResponse(c, result)

}

func Delete(c *gin.Context) {

	var data DeleteData
	err := c.ShouldBindQuery(&data) //绑定请求的参数到结构体中
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	err = confessionService.Delete(data.ConfessionID) //删除某一条帖子
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			apiException.AbortWithException(c, apiException.ConfessionNotFound, err)
		} else {
			apiException.AbortWithException(c, apiException.ServerError, err)
		}
		return

	}
	utils.JsonSuccessResponse(c, nil)

}

func Update(c *gin.Context) {

	var data UpdateData
	err := c.ShouldBindJSON(&data) //绑定
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}

	err = confessionService.Update(int(data.ID), data.Content) //更新数据
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			apiException.AbortWithException(c, apiException.ConfessionNotFound, err)
		} else {
			apiException.AbortWithException(c, apiException.ServerError, err)
		}
		return
	}
	utils.JsonSuccessResponse(c, nil)

}
