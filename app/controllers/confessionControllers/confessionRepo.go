package confession

import (
	apiexception "CONFESSION-WALL-BACKEND/app/apiException"
	"CONFESSION-WALL-BACKEND/app/models"
	confessionservices "CONFESSION-WALL-BACKEND/app/services/confessionServices"
	userservices "CONFESSION-WALL-BACKEND/app/services/userServices"
	"CONFESSION-WALL-BACKEND/app/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReleaseData struct {
	Content   string `json:"content" binding:"required"`
	Picture   string `json:"picture"`
	Anonymous bool   `json:"anonymous" binding:"required"`
	Invisible bool   `json:"invisible" binding:"required"`
	UserID    int    `json:"user_id" binding:"required"`
}


type PageQuery struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

type ShowConfessionData struct{
	Name string `json:"name"`
	Likes int `json:"likes"`
	Comments int `json:"comments"`
	Content string `json:"content"`
	Views int `json:"views"`
	Avatar string `json:"avatar"`
	Picture string `json:"pictrue"`
	CreatedAt time.Time `json:"created_time"`
	UpdatedAt time.Time `json:"updated_time"`
}

type MyConfessionData struct{
	UserID int `json:"user_id"`
}
type DeleteData struct {
	ConfessionID int `form:"confession_id" json:"id" binding:"required"`
}

type UpdateData struct {
	ID      uint   `json:"post_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func Release(c *gin.Context) {

	var data ReleaseData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiexception.AbortWithException(c, apiexception.ParamError, err)
		return
	}
	user,err:=userservices.SeekUser(data.UserID)
		if err!=nil{
			apiexception.AbortWithException(c, apiexception.ServerError, err)
			return
		}
	var name string
	if data.Anonymous{
		name=uuid.NewString()
		

	}else{
		name=user.Name
	}
	
	
	
	err = confessionservices.ReleaseConfession(&models.Confession{
		UserID: data.UserID,
		Avatar: user.Avatar,
		Content: data.Content,
		Picture: data.Picture,
		Invisible: data.Invisible,
		Anonymous: data.Anonymous,
		Name: name,
		Likes: 0,
		Views: 0,


	})
	if err != nil {
		apiexception.AbortWithException(c, apiexception.ServerError, err)
		return
	}
	utils.JsonSuccessResponse(c, nil)

}

func ShowConfessions(c *gin.Context) {
	var data PageQuery
	err := c.ShouldBindQuery(&data)
	if err != nil {
		apiexception.AbortWithException(c, apiexception.ServerError, err)
	}
	offset := (data.Page - 1) * data.PageSize                      //计算偏移量
	result, err := confessionservices.ShowConfession(offset, data.PageSize) //result为post中的结构体数组
	if err != nil {
		apiexception.AbortWithException(c, apiexception.ServerError, err)
	}

	var NewConfession []ShowConfessionData
	for _, data := range result {
		NewConfession = append(NewConfession,
			ShowConfessionData{
				Name: data.Name,
				Likes: data.Likes,
				Comments: data.Comments,
				Views: data.Views,
				Content: data.Content,
				Avatar:data.Avatar ,
				Picture: data.Picture,
				CreatedAt: data.CreatedAt,
				UpdatedAt: data.UpdatedAt,
			
			})
	}
	
	utils.JsonSuccessResponse(c, NewConfession)
}

func ShowMyConfessions(c *gin.Context){
	var data MyConfessionData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		apiexception.AbortWithException(c, apiexception.ParamError, err)
		return
	}
	result, err := confessionservices.ShowMyConfession(data.UserID) //result为post中的结构体数组
	if err != nil {
		apiexception.AbortWithException(c, apiexception.ServerError, err)
	}
	utils.JsonSuccessResponse(c,result)

}

func Delete(c *gin.Context) {

	var data DeleteData
	err := c.ShouldBindQuery(&data) //绑定请求的参数到结构体中
	if err != nil {
		apiexception.AbortWithException(c, apiexception.ParamError, err)
		return
	}

	err = confessionservices.Delete(data.ConfessionID) //删除某一条帖子
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			apiexception.AbortWithException(c, apiexception.ConfessionNotFound, err)
		} else {
			apiexception.AbortWithException(c, apiexception.ServerError, err)
		}
		return

	}
	utils.JsonSuccessResponse(c, nil)

}

func Update(c *gin.Context) {
	
	var data UpdateData
	err := c.ShouldBindJSON(&data) //绑定
	if err != nil {
		apiexception.AbortWithException(c, apiexception.ParamError, err)
		return
	}

	err = confessionservices.Update(int(data.ID),data.Content) //更新数据
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			apiexception.AbortWithException(c, apiexception.ConfessionNotFound, err)
		} else {
			apiexception.AbortWithException(c, apiexception.ServerError, err)
		}
		return
	}
	utils.JsonSuccessResponse(c, nil)

}