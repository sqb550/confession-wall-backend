package main

import (
	"confession-wall-backend/app/midwares"
	"confession-wall-backend/app/utils"
	"confession-wall-backend/config/database"
	"confession-wall-backend/config/router"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/gin-contrib/cors"
)

func main() {
	database.Init()
	utils.InitRedis()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://124.220.30.150/"}, 
    AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
    }))
	err := os.MkdirAll("./uploads", 0755)
	if err != nil {
		log.Fatal("Create upload directory error:", err)

	}
	r.Static("/uploads", "./uploads")
	r.Use(midwares.ErrHandler())

	c := cron.New()
	_, err1 := c.AddFunc("*/5 * * * *", utils.SyncCacheToDB)
	_, err2 := c.AddFunc("* * * * *", utils.ScheduleRelease)
	if err1 != nil || err2 != nil {
		log.Fatal("Add tasks error:", err)
	}

	c.Start()
	defer c.Stop()

	router.Init(r)
	err = r.Run(":8080")
	if err != nil {
		log.Fatal("Server start error:", err)
	}
}
