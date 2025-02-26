package main

import (
	"go-spe/api"
	"go-spe/pkg/cache"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("error loading .env file: %v", err)
		return
	}

	cache.InitRedis()

	// Setup API
	router := gin.Default()
	api.SetupRoutes(router)

	// Start server
	log.Println("Starting server on port 8080...")
	router.Run(":8080")
}
