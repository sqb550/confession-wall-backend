package confessionController

import (
	"confession-wall-backend/app/apiException"
	"confession-wall-backend/app/utils"
	"crypto/md5"
	"encoding/hex"
	"io"
	"path/filepath"

	"github.com/gin-gonic/gin"
)


func UploadPicture(c *gin.Context){
	form,err:=c.MultipartForm()
	if err!=nil{
		apiException.AbortWithException(c,apiException.ServerError,err)
		return
	}
	files:=form.File["files"]
	if len(files)==0{
		apiException.AbortWithException(c,apiException.FileNotFound,nil)
	}
	if len(files)>9{
		apiException.AbortWithException(c,apiException.FileNumberError,nil)
	}
	var result []string
	for _,file:=range files{
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
			result=append(result,url)
			continue
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
		result=append(result,url)
	}
	utils.JsonSuccessResponse(c,result)
}
