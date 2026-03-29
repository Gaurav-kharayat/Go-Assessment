package main

import (
	"inventory-service/db"
	"inventory-service/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	db.InitDB()

	r := gin.Default()

	routes.SetupRoutes(r)

	r.Run(":8080")
}
