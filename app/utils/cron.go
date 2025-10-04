package utils

import (
	"confession-wall-backend/app/models"
	"confession-wall-backend/config/database"
	"context"
	"log"
	"strconv"
	"time"
)

// 定时任务同步缓存到数据库
func SyncCacheToDB() {
	ctx := context.Background()

	keys, err := ScanPosts(ctx)
	if err != nil {
		log.Printf("Error getting keys: %v\n", err)
		return
	}
	for _, key := range keys {
		post_id := key[len("post:"):]

		val, err := redisClient.HGetAll(ctx, key).Result()
		if err != nil {
			log.Printf("Error getting post data from cache: %v\n", err)
		}
		likes, _ := strconv.Atoi(val["likes"])
		views, _ := strconv.Atoi(val["views"])
		err = database.DB.Model(&models.Post{}).Where("id =?", post_id).Updates(map[string]interface{}{
			"likes": likes,
			"views": views,
		}).Error
		if err != nil {
			log.Printf("Error updating database: %v\n", err)
			
		}
	}
}

func ScheduleRelease() {
	now := time.Now()
	result := database.DB.Model(&models.Post{}).Where("release_status=? AND release_time<=?", false, now).Update("release_status", true)
	err := result.Error
	if err != nil {
		log.Printf("Error upload posts: %v\n", err)
		return
	}
}
