package utils

import "github.com/gin-gonic/gin"

func JsonResponse(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func JsonSuccessResponse(c *gin.Context, data interface{}) {
	JsonResponse(c, 200, "success", data)
}

func JsonErrorResponse(c *gin.Context, code int, msg string) {
	JsonResponse(c, code, msg, nil)
}