package router

import (
	"confession-wall-backend/app/controllers/postController"
	"confession-wall-backend/app/controllers/userController"
	"confession-wall-backend/app/midwares"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	r.POST("/api/login", userController.Login)
	r.POST("/api/reg", userController.Register)
	const pre = "/api"
	api := r.Group(pre)
	api.Use(midwares.AuthMiddleware())
	{
		api.PUT("/user/name", userController.UpdateName)
		api.PUT("/user/password", userController.UpdatePassword)
		api.PUT("/avatar",midwares.RateLimiter() ,userController.UploadAvatar)

		api.POST("/post", postController.Release)
		api.PUT("/post", postController.Update)
		api.DELETE("/post", postController.Delete)
		api.GET("/post", postController.QueryPosts)
		api.POST("/picture",midwares.RateLimiter() ,postController.UploadPicture)
		api.GET("/reply", postController.QueryComment)
		api.POST("/reply", postController.Comment)
		api.POST("/like", postController.Like)
		api.GET("/mypost", postController.QueryMyPosts)
		api.GET("/blacklist", postController.ShowBlock)
		api.POST("/block", postController.Block)
		api.GET("/hotRank", postController.GetHotRank)
	}

}
