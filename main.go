package main

import (
	"confession-wall-backend/app/utils"
	"confession-wall-backend/config/database"
	"confession-wall-backend/config/router"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	database.Init()
	r := gin.Default()
	err:=os.MkdirAll("./uploads",0755)
	if err!=nil{
		log.Fatal("Create upload directory error:", err)
	
	}
	r.Static("/uploads","./uploads")
	
	c:=cron.New()
	_,err=c.AddFunc("**/5****",utils.SyncCacheToDB)
	c.Start()
	defer c.Stop()

	router.Init(r)
	err = r.Run(":8080")
	if err != nil {
		log.Fatal("Server start error:", err)
	}
}
