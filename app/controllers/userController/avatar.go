package user

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



func UploadAvatar(c *gin.Context){
	userID,_:=c.Get("user_id")
	userIDInt, _ := userID.(int)
	file,err:=c.FormFile("file")
	if err!=nil{
		apiException.AbortWithException(c,apiException.FlieError,err)
		return
	}
	src,err:=file.Open()
	if err!=nil{
		apiException.AbortWithException(c,apiException.FlieError,err)
	}
	hash:=md5.New()
	_,err=io.Copy(hash,src)
	if err!=nil{
		apiException.AbortWithException(c,apiException.ServerError,err)
	}
	md5Str:=hex.EncodeToString(hash.Sum(nil))
	path,flag,err:=utils.GetFileHashFromCache(md5Str,c)
	if err!=nil{
		apiException.AbortWithException(c,apiException.ServerError,err)
	}
	if flag{
		url := "http://localhost:8080" + path
		err=userService.UploadAvatar(userIDInt,url)
		if err!=nil{
			apiException.AbortWithException(c,apiException.ServerError,err)
			return
		}

	utils.JsonSuccessResponse(c,nil)
	return
	}

	ext:=filepath.Ext(file.Filename)
	uniqueName := md5Str+ext

	
	savePath := filepath.Join("./uploads", uniqueName)
	err=utils.SetFileHashToCache(md5Str,savePath,c)
	if err!=nil{
		apiException.AbortWithException(c,apiException.ServerError,err)
	}

	err = c.SaveUploadedFile(file, savePath)
	if err != nil {
		apiException.AbortWithException(c,apiException.FlieError,err)
		return
	}

	url := "http://localhost:8080/uploads/" + uniqueName
	err=userService.UploadAvatar(userIDInt,url)
	if err!=nil{
		apiException.AbortWithException(c,apiException.ServerError,err)
	}

	utils.JsonSuccessResponse(c,nil)

	
}