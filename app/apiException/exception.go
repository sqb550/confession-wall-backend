package apiException

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Code int
	Msg  string
}

// Error implements error.
func (e *Error) Error() string {
	return fmt.Sprintf("error:code=%d,message=%s",e.Code,e.Msg)
}

var (
	ServerError         = NewError(200500, "系统异常，请稍后重试")
	ParamError          = NewError(200501, "参数错误")
	UsernameError       = NewError(200502, "用户名必须为纯数字")
	PasswordLengthError = NewError(200503, "密码长度必须在8-16位")
	UserExist           = NewError(200504, "用户名已存在")
	UserNotFound        = NewError(200505, "用户名不存在")
	PostNotFound        = NewError(200506, "该表白不存在")
	PasswordError       = NewError(200507, "密码错误")
	LikeError           = NewError(200508, "点赞失败")
	FlieError           = NewError(200509, "图片上传失败")
	RedisError          = NewError(200510, "缓存失败")
	FileNotFound        = NewError(200511, "未找到上传的图片")
	FileNumberError     = NewError(200512, "图片最多上传九张")
	CancelLikeError     = NewError(200513, "取消点赞失败")
	HotRankError        = NewError(200514, "排行榜获取失败")
	FileSizeError=NewError(200515,"文件大小不能超过5MB")

	NotFound = NewError(200404, http.StatusText(http.StatusNotFound))
)

func NewError(code int, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

func AbortWithException(c *gin.Context, apiError *Error, err error) {
	_ = c.AbortWithError(200, apiError)
}
