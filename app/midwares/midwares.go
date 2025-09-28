package midwares

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/utils"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

// 请求错误中间件
func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if err != nil {
				var apiErr *apiException.Error

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
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"code": 401,
				"msg":  "请先登录",
				"data": nil,
			})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, "", 2)
		if len(parts) != 2 || parts[0] != "bearer" {
			c.JSON(401, gin.H{
				"code": 401,
				"msg":  "登录异常，请重新登录",
				"data": nil,
			})
			c.Abort()
			return

		}
		tokenString := parts[1]
		token, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{
				"code": 401,
				"msg":  "登录异常，请重新登录",
				"data": nil,
			})
			c.Abort()
			return

		}
		claims, err := utils.ExtractClaims(token)
		if err != nil {
			c.JSON(401, gin.H{
				"code": 401,
				"msg":  "登录异常，请重新登录",
				"data": nil,
			})
			c.Abort()
			return

		}
		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}
