package apiexception

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Code int
	Msg  string
}

// Error implements error.
func (e *Error) Error() string {
	panic("unimplemented")
}

var (
	ServerError         = NewError(200500, "系统异常，请稍后重试")
	ParamError          = NewError(200501, "参数错误")
	UsernameError       = NewError(200502, "用户名必须为纯数字")
	PasswordLengthError = NewError(200503, "密码长度必须在8-16位")
	UserExist           = NewError(200504, "用户名已存在")
	UserNotFound        = NewError(200505, "用户名不存在")
	ConfessionNotFound        = NewError(200506, "该表白不存在")
	PasswordError       = NewError(200507, "密码错误")
	USerIDError         = NewError(200508, "user_id传递失败")
	LikeError           = NewError(200509, "点赞失败")

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
