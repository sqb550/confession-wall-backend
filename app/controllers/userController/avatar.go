package userController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/services/userService"
	"confession-wall-backend/app/utils"
	"crypto/md5"
	"encoding/hex"
	"io"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadAvatar(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID,ok:=val.(float64)
	if !ok{
		apiException.AbortWithException(c,apiException.ServerError,nil)
		return
	}
	userIDInt:=int(userID)
	file, err := c.FormFile("file")
	if err != nil {
		apiException.AbortWithException(c, apiException.FlieError, err)
		return
	}
	maxSize:=int64(5*1024*1024)
	if file.Size>maxSize{
		apiException.AbortWithException(c,apiException.FileSizeError,err)
		return
	}
	src, err := file.Open()
	if err != nil {
		apiException.AbortWithException(c, apiException.FlieError, err)
		return
	}
	defer src.Close()
	hash := md5.New()
	_, err = io.Copy(hash, src)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	md5Str := hex.EncodeToString(hash.Sum(nil))
	path, flag, err := utils.GetFileHashFromCache(md5Str, c)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	if flag {
		url := "http://127.0.0.1:8080" + path
		err = userService.UploadAvatar(userIDInt, url)
		if err != nil {
			apiException.AbortWithException(c, apiException.ServerError, err)
			return
		}
		utils.JsonSuccessResponse(c, nil)
		return
	}
	ext := filepath.Ext(file.Filename)
	uniqueName := md5Str + ext
	savePath := filepath.Join("./uploads", uniqueName)
	err = utils.SetFileHashToCache(md5Str, savePath, c)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	err = c.SaveUploadedFile(file, savePath)
	if err != nil {
		apiException.AbortWithException(c, apiException.FlieError, err)
		return
	}
	url := "http://127.0.0.1:8080/uploads/" + uniqueName
	err = userService.UploadAvatar(userIDInt, url)
	if err != nil {
		apiException.AbortWithException(c, apiException.ServerError, err)
		return
	}
	err=userService.UpdateAvatar(userIDInt,url)
	if err!=nil{
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	utils.JsonSuccessResponse(c, url)
}
