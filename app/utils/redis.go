package utils

import (
	"confession-wall-backend/app/models"
	"confession-wall-backend/config/database"
	"context"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
}
//浏览量自增
func IncrViewCount(confessionID int,c *gin.Context)error{
	hashKey:="confession:"+strconv.Itoa(confessionID)
	_,err:=redisClient.HIncrBy(c,hashKey,"views",1).Result()
	return err
}
func CheckUserLike(userID int,confessionID int,c *gin.Context)(bool,error){
	likeUserKey:="confession:users"+strconv.Itoa(confessionID)
	isLiked,err:=redisClient.SIsMember(c,likeUserKey,userID).Result()
	return isLiked,err
}
func LikeHandler(userID int,confessionID int,c*gin.Context)error{
	pipe:=redisClient.Pipeline()
	likeUserKey:="confession:users"+strconv.Itoa(confessionID)
	hashKey:="confession:"+strconv.Itoa(confessionID)
	pipe.SAdd(c,likeUserKey,userID)
	pipe.HIncrBy(c,hashKey,"likes",1)
	_,err:=pipe.Exec(c)
	return err
}

func CancelLikeHandler(userID int,confessionID int,c *gin.Context)error{
	pipe:=redisClient.Pipeline()
	likeUserKey:="confession:users"+strconv.Itoa(confessionID)
	hashKey:="confession:"+strconv.Itoa(confessionID)
	pipe.SRem(c,likeUserKey,userID)
	pipe.HIncrBy(c,hashKey,"likes",-1)
	_,err:=pipe.Exec(c)
	return err
}
func GetLikeAndViews(confessionID int, c*gin.Context)(string,string,bool,error){
	hashKey:="confession:"+strconv.Itoa(confessionID)
	stats,err:=redisClient.HGetAll(c,hashKey).Result()
	if err !=nil{
		if err==redis.Nil{
			return "","",false,nil 
		}else{
			return "","",false,err
		}
	}
	return stats["likes"],stats["views"],true,nil	
}
func SetHashToCache(confessionID int,c *gin.Context)error{
	hashKey:="confession:"+strconv.Itoa(confessionID)
	err:=redisClient.HSet(c,hashKey,map[string]interface{}{
		"likes":0,
		"views":0,
	}).Err()
	return err
}



//设置文件的hash值
func SetFileHashToCache(md5str string,path string,c *gin.Context)error{
	return redisClient.Set(c,"md5:"+md5str,path,0).Err()
}
func GetFileHashFromCache(md5Str string, c *gin.Context) (string, bool, error) {
	path, err := redisClient.Get(c, "md5:"+md5Str).Result()
	if err != nil {
		if err == redis.Nil {
			return "", false, nil
		}
		return "", false, err
	}

	return path, true, nil
}

// 定时任务同步缓存到数据库
func SyncCacheToDB() {
	ctx := context.Background()

	keys, err := redisClient.Keys(ctx, "confession:*").Result()
	if err != nil {
		log.Printf("Error getting keys: %v\n", err)
		return
	}

	for _, key := range keys {
		confession_id := key[len("confession:"):]

		val, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			log.Printf("Error getting post data from cache: %v\n", err)
			return
		}

		err = database.DB.Model(&models.Confession{}).Where("id =?",confession_id).Update("likes", val).Error
		if err != nil {
			log.Printf("Error updating database: %v\n", err)
			return
		}
	}
}
