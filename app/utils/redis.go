package utils

import (
	"CONFESSION-WALL-BACKEND/app/models"
	"CONFESSION-WALL-BACKEND/config/database"
	"context"
	"fmt"
	"log"
	"time"

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

func GetConfessionCacheKey(ConfessionID uint) string {
	return fmt.Sprintf("confession:%d", ConfessionID)
}

// 从缓存获取点赞值
func GetConfessionFromCache(ConfessionID uint, c *gin.Context) (int, bool, error) {
	key := GetConfessionCacheKey(ConfessionID)
	value, err := redisClient.Get(c, key).Int()
	if err != nil {
		if err == redis.Nil {
			return 0, false, nil
		}
		return 0, false, err
	}

	return value, true, nil
}

// likes原子自增
func LikesIncr(ConfessionID int, c *gin.Context) error {
	key := GetConfessionCacheKey(uint(ConfessionID))
	_, err := redisClient.Incr(c, key).Result()
	return err
}

// 设置缓存
func SetConfessionToCache(confession *models.Confession, c *gin.Context) error {
	key := GetConfessionCacheKey(confession.ID)

	return redisClient.Set(c, key, confession.Likes, 30*time.Minute).Err()
}

// 定时任务同步缓存到数据库
func SyncCacheToDB() {
	ctx := context.Background()

	keys, err := redisClient.Keys(ctx, "post:*").Result()
	if err != nil {
		log.Printf("Error getting keys: %v\n", err)
		return
	}

	for _, key := range keys {
		post_id := key[len("post:"):]

		val, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			log.Printf("Error getting post data from cache: %v\n", err)
			return
		}

		err = database.DB.Model(&models.Confession{}).Where("id =?", post_id).Update("likes", val).Error
		if err != nil {
			log.Printf("Error updating database: %v\n", err)
			return
		}
	}
}
