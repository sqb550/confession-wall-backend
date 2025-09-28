package postController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/services/postService"
	"strconv"

	"confession-wall-backend/app/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)



type PageQuery struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

type ShowPostData struct {
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
	var data PageQuery
	err := c.ShouldBindQuery(&data)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
	}
	offset := (data.Page - 1) * data.PageSize                  //计算偏移量
	result, err := postService.Showpost(offset, data.PageSize) //result为post中的结构体数组
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
	}

	var Newpost []ShowPostData
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
		Newpost = append(Newpost,
			ShowPostData{
				Name:      data.Name,
				Likes:     likes,
				Comments:  data.Comments,
				Views:     views,
				Content:   data.Content,
				Avatar:    data.Avatar,
				Picture:   data.Picture,
				UpdatedAt: data.UpdatedAt,
			})
		err = utils.IncrViewCount(int(data.ID), c)
		if err != nil {
			apiException.AbortWithException(c, apiException.ServerError, nil)
			return
		}
	}

	utils.JsonSuccessResponse(c, Newpost)
}

func ShowMyPosts(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDInt, _ := userID.(int)
	result, err := postService.ShowMypost(userIDInt)
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
