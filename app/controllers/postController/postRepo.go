package postController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/services/postService"
	"confession-wall-backend/app/services/userService"
	"strconv"

	"confession-wall-backend/app/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)



type PageData struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type PostData struct {
	PostID      int      `json:"post_id"`
	Name        string    `json:"name"`
	Likes       int       `json:"likes"`
	Comments    int       `json:"comments"`
	Content     string    `json:"content"`
	Views       int       `json:"views"`
	Avatar      string    `json:"avatar"`
	Picture     []string  `json:"pictrue"`
	ReleaseTime time.Time `json:"release_time"`
	UpdatedAt   time.Time `json:"updated_time"`
}

type DeleteData struct {
	PostID int `form:"post_id" json:"id" binding:"required"`
}

type UpdateData struct {
	ID      uint   `json:"post_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}



func ShowPosts(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID,ok:=val.(float64)
	if !ok{
		apiException.AbortWithException(c,apiException.ServerError,nil)
		return
	}
	userIDInt:=int(userID)
	var data PageData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	offset := (data.Page - 1) * data.PageSize 
	var blockedID []int
	blocks,err:=postService.ShowBlock(userIDInt) 
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	for _,block:=range blocks{
		blockedID=append(blockedID, block.BlockedID)
	}              
	result, err := postService.ShowPost(offset, data.PageSize,blockedID) //result为post中的结构体数组
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	var Newpost []PostData
	for _, data := range result {
		user, err := userService.SeekUser(data.UserID)
		if err != nil {
			apiException.AbortWithException(c, apiException.ServerError, err)
			return
		}
		likesStr, viewsStr, flag, err := utils.GetLikeAndViews(int(data.ID), c)
		likes := 0
		views := 0
		if err != nil {
			apiException.AbortWithException(c, apiException.ServerError, err)
			return
		}
		if flag {
			likes, _ = strconv.Atoi(likesStr)
			views, _ = strconv.Atoi(viewsStr)
		}
		pictures,err:=postService.GetPictures(int(data.ID))
		if err!=nil{
			apiException.AbortWithException(c,apiException.ServerError,err)
			return
		}
		urls:=make([]string,0)
		for _,picture:=range pictures{
			urls = append(urls, picture.URL)
		}
		Newpost = append(Newpost,
			PostData{
				PostID: int(data.ID),
				Name:      data.Name,
				Likes:     likes,
				Comments:  data.Comments,
				Views:     views,
				Content:   data.Content,
				Avatar:    user.Avatar,
				Picture:   urls,
				UpdatedAt: data.UpdatedAt,
			})
		err = utils.IncrViewCount(int(data.ID), c)
		if err != nil {
			apiException.AbortWithException(c, apiException.ServerError, nil)
			return
		}
		err=utils.UpdateHot(c,int(data.ID),likes,views+1)
		if err != nil {
			apiException.AbortWithException(c, apiException.ServerError, nil)
			return
		}
	}
	utils.JsonSuccessResponse(c, Newpost)
}

func ShowMyPosts(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID,ok:=val.(float64)
	if !ok{
		apiException.AbortWithException(c,apiException.ServerError,nil)
		return
	}
	userIDInt:=int(userID)
	result, err := postService.ShowMyPost(userIDInt)
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

	err = postService.Delete(data.PostID) //删除某一条帖子
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			apiException.AbortWithException(c, apiException.PostNotFound, err)
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

	err = postService.Update(int(data.ID), data.Content) //更新数据
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			apiException.AbortWithException(c, apiException.PostNotFound, err)
		} else {
			apiException.AbortWithException(c, apiException.ServerError, err)
		}
		return
	}
	utils.JsonSuccessResponse(c, nil)

}
