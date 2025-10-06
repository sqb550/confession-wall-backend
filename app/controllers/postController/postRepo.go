package postController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/services/postService"
	"confession-wall-backend/config/database"
	"strconv"

	"confession-wall-backend/app/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)



type PageData struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

type PostData struct {
	PostID      int      `json:"post_id"`
	Name        string    `json:"name"`
	UserID      int `json:"user_id"`
	Likes       int       `json:"likes"`
	Comments    int       `json:"comments"`
	Content     string    `json:"content"`
	Views       int       `json:"views"`
	Avatar      string    `json:"avatar"`
	Picture     []string  `json:"picture"`
	ReleaseTime time.Time `json:"release_time"`
	UpdatedAt   time.Time `json:"updated_time"`
}
type ShowPosts struct{
	Total int64 `json:"total"`
	PostList []PostData `json:"post_list"`
}
type MyPostData struct{
	Content string `json:"content"`
	PostID int `json:"post_id"`
	Urls []string `json:"urls"`
}

type DeleteData struct {
	PostID int `form:"post_id" json:"id" binding:"required"`
}

type UpdateData struct {
	ID      uint   `json:"post_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}



func QueryPosts(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID,ok:=val.(float64)
	if !ok{
		apiException.AbortWithException(c,apiException.ServerError,nil)
		return
	}
	userIDInt:=int(userID)
	var data PageData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	offset := (data.Page - 1) * data.PageSize 
	var blockedID []int
	blocks,err:=postService.QueryBlock(userIDInt) 
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	for _,block:=range blocks{
		blockedID=append(blockedID, block.BlockedID)
	}              
	result, err := postService.QueryPost(offset, data.PageSize,blockedID) //result为post中的结构体数组
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}

	var Newpost []PostData
	for _, data := range result {
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
				UserID: data.UserID,
				Name:      data.Name,
				Likes:     likes,
				Comments:  data.Comments,
				Views:     views,
				Content:   data.Content,
				Avatar:    data.Avatar,
				Picture:   urls,
				UpdatedAt: data.UpdatedAt,
			})
	}
	total,err:=postService.CountPosts()
	if err!=nil{
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	showResult:=ShowPosts{
		Total: total,
		PostList: Newpost,
	}
	utils.JsonSuccessResponse(c, showResult)
}

func QueryMyPosts(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID,ok:=val.(float64)
	if !ok{
		apiException.AbortWithException(c,apiException.ServerError,nil)
		return
	}
	userIDInt:=int(userID)
	result, err := postService.QueryMyPost(userIDInt)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
	}
	var myPosts []MyPostData
	for _ ,data:=range result{
		pictures,err:=postService.GetPictures(int(data.ID))
		if err!=nil{
			apiException.AbortWithException(c,apiException.ServerError,err)
			return
		}
		urls:=make([]string,0)
		for _,picture:=range pictures{
			urls = append(urls, picture.URL)
		}
		myPosts=append(myPosts, MyPostData{
			Content: data.Content,
			PostID: int(data.ID),
			Urls: urls,
		})

	}
	utils.JsonSuccessResponse(c, myPosts)

}

func Delete(c *gin.Context) {

	var data DeleteData
	err := c.ShouldBindQuery(&data) //绑定请求的参数到结构体中
	if err != nil {
		apiException.AbortWithException(c, apiException.ParamError, err)
		return
	}
	tx := database.DB.Begin()
	if tx.Error != nil {
		apiException.AbortWithException(c, apiException.ServerError, tx.Error)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err = postService.Delete(tx,data.PostID) //删除某一条帖子
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			tx.Rollback()
			apiException.AbortWithException(c, apiException.PostNotFound, err)
		} else {
			tx.Rollback()
			apiException.AbortWithException(c, apiException.ServerError, err)
		}
		return
	}
	err=postService.DeletePicture(tx,data.PostID)
	if err!=nil{
		tx.Rollback()
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	err=postService.DeleteComment(tx,data.PostID)
	if err!=nil{
		tx.Rollback()
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	_ = tx.Commit()

	err=utils.Delete(c,data.PostID)
	if err!=nil{
		apiException.AbortWithException(c,apiException.ServerError,err)
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
