package midwares

import (
	apiexception "CONFESSION-WALL-BACKEND/app/apiException"
	"CONFESSION-WALL-BACKEND/app/utils"
	"errors"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// 请求错误中间件
func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if err != nil {
				var apiErr *apiexception.Error

				if errors.As(err, &apiErr) {
					utils.JsonErrorResponse(c, apiErr.Code, apiErr.Msg)
				}
				return
			}
		}
	}
}

// authMiddleware 身份鉴权中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user_id := session.Get("user_id")
		if user_id == nil {
			c.JSON(401, gin.H{
				"code": 401,
				"msg":  "请先登录",
				"data": nil,
			})
			c.Abort()
			return
		}
		c.Set("user_id", user_id)
		c.Next()
	}
}