package utils

import (
	"confession-wall-backend/config/config"
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:   config.Config.GetString("redis.addr"),
		Password:config.Config.GetString("redis.password"),
		DB:      config.Config.GetInt("redis.db"),
	})

}

// 浏览量自增
func IncrViewCount(postID int, c *gin.Context) error {
	hashKey := "post:" + strconv.Itoa(postID)
	_, err := redisClient.HIncrBy(c, hashKey, "views", 1).Result()
	return err
}
func CheckUserLike(userID int, postID int, c *gin.Context) (bool, error) {
	likeUserKey := "post:users" + strconv.Itoa(postID)
	isLiked, err := redisClient.SIsMember(c, likeUserKey, userID).Result()
	return isLiked, err
}
func LikeHandler(userID int, postID int, c *gin.Context) error {
	pipe := redisClient.Pipeline()
	likeUserKey := "post:users" + strconv.Itoa(postID)
	hashKey := "post:" + strconv.Itoa(postID)
	pipe.SAdd(c, likeUserKey, userID)
	pipe.HIncrBy(c, hashKey, "likes", 1)
	_, err := pipe.Exec(c)
	return err
}

func CancelLikeHandler(userID int, postID int, c *gin.Context) error {
	pipe := redisClient.Pipeline()
	likeUserKey := "post:users" + strconv.Itoa(postID)
	hashKey := "post:" + strconv.Itoa(postID)
	pipe.SRem(c, likeUserKey, userID)
	pipe.HIncrBy(c, hashKey, "likes", -1)
	_, err := pipe.Exec(c)
	return err
}
func GetLikeAndViews(postID int, c *gin.Context) (string, string, bool, error) {
	hashKey := "post:" + strconv.Itoa(postID)
	stats, err := redisClient.HGetAll(c, hashKey).Result()
	if err != nil {
		if err == redis.Nil {
			return "", "", false, nil
		} else {
			return "", "", false, err
		}
	}
	return stats["likes"], stats["views"], true, nil
}
func ScanPosts(ctx context.Context)([]string,error){
	var allKeys []string
	cursor:=uint64(0)
	for{
		var keys []string
		keys,cursor,err:=redisClient.Scan(ctx,cursor,"post:*",10).Result()
		if err!=nil{
			return nil,err
		}
		allKeys=append(allKeys,keys...)
		if cursor== 0 {
			break
		}
	}
	return allKeys,nil
}



func UpdateHot(c *gin.Context, postID int, likes int, views int) error {
	hot := likes*3 + views*2
	postIDStr := strconv.Itoa(postID)
	ctx := c.Request.Context()
	err := redisClient.ZAdd(ctx, "post:hot:rank", []redis.Z{
		{
			Score:  float64(hot),
			Member: postIDStr,
		},
	}...).Err()
	return err
}
func GetTopHotRank(c *gin.Context) ([]int, error) {
	members, err := redisClient.ZRevRange(c, "post:hot:rank", 0, 9).Result()
	if err != nil {
		return nil, err
	}
	postIDs := make([]int,0)
	for _, member := range members {
		postID, err := strconv.Atoi(member)
		if err != nil {
			return nil, err
		}
		postIDs=append(postIDs,postID)
	}
	return postIDs, nil
}
func Delete(c *gin.Context,postID int)error{
	postIDStr := strconv.Itoa(postID)
	ctx := c.Request.Context()
	err:=redisClient.ZRem(ctx,"post:hot:rank",postIDStr).Err()
	return err
}


// 设置文件的hash值
func SetFileHashToCache(md5str string, path string, c *gin.Context) error {
	return redisClient.Set(c, "md5:"+md5str, path, 0).Err()
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
