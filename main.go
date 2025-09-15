package main

import (
	"CONFESSION-WALL-BACKEND/config/database"
	"CONFESSION-WALL-BACKEND/config/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()
	r := gin.Default()


	
	router.Init(r)
	err := r.Run(":8080")
	if err != nil {
		log.Fatal("Server start error:", err)
	}
}