package main

import (
	"go-spe/api"
	"go-spe/pkg/cache"
	"go-spe/pkg/messaging"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
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
	messaging.InitRabbitMQ()

	// Setup API
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Start ProcessTransactions in a goroutine
	go messaging.ProcessTransactions()

	// Apply rate limiting middleware globally
	router.Use(RateLimitMiddleware)

	api.SetupRoutes(router)

	// Start server
	log.Println("Starting server on port 8080...")
	router.Run(":8080")

	select {} // Blocks forever, keeping the program alive for the goroutine to run
}

// RateLimitMiddleware checks if the IP address is rate-limited
func RateLimitMiddleware(c *gin.Context) {
	ipAddress := c.ClientIP()

	// Check if the request is allowed based on rate limiting logic
	if !cache.IsRateLimited(cache.RedisClient, ipAddress) {
		// If rate limit exceeded, reject the request with a 429 status code (Too Many Requests)
		c.JSON(http.StatusTooManyRequests, gin.H{
			"message": "Rate limit exceeded. Please try again later.",
		})
		c.Abort()
		return
	}

	// Proceed to the next handler if rate limit is not exceeded
	c.Next()
}
