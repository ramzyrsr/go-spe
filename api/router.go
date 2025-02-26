package api

import (
	"net/http"

	"go-spe/internal/auth"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/", auth.SignatureMiddleware(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Welcome to API"})
		})
	}
}
