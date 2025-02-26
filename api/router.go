package api

import (
	"log"
	"net/http"

	"go-spe/internal/auth"
	"go-spe/internal/domain/repository"
	"go-spe/internal/domain/service"
	"go-spe/pkg/db"
	"go-spe/pkg/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	v1 := router.Group("/api/v1")
	{
		// Initialize TransactionHandler
		transactionRepo := repository.NewTransactionRepositor(dbConn)
		transactionService := service.NewTransactionService(transactionRepo)
		transactionHandler := handler.NewTransactionHandler(transactionService)

		v1.GET("/", auth.SignatureMiddleware(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Welcome to API"})
		})
		v1.POST("/check-status", auth.SignatureMiddleware(), transactionHandler.CheckTransactionStatus)
	}
}
