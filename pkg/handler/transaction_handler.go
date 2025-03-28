package handler

import (
	"encoding/json"
	"go-spe/internal/domain/models"
	"go-spe/internal/domain/service"
	"go-spe/pkg/messaging"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) TransactionNotification(c *gin.Context) {
	var trx *models.Transaction
	if err := c.ShouldBindJSON(&trx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Publish to RabbitMQ
	body, _ := json.Marshal(trx)
	err := messaging.Channel.Publish(
		"", "transaction_notifications", false, false,
		amqp091.Publishing{ContentType: "application/json", Body: body},
	)
	if err != nil {
		log.Println("Failed to publish to RabbitMQ:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue transaction"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"code": "00", "message": "success"})
}

func (h *TransactionHandler) CheckTransactionStatus(c *gin.Context) {
	type RequestBody struct {
		RequestID  string `json:"request_id"`
		BillNumber string `json:"bill_number"`
	}

	var reqBody RequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.service.GetTransactionStatus(reqBody.RequestID, reqBody.BillNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if transaction == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "01",
			"status": "failed", "data": "Transaction not found"})
		return
	}

	// Marshal the transaction struct into a map
	transactionMap := make(map[string]interface{})
	transactionBytes, _ := json.Marshal(transaction)
	json.Unmarshal(transactionBytes, &transactionMap)

	// Create the response
	response := map[string]interface{}{
		"code":   "00",
		"status": "success",
	}

	// Merge the transaction map into the response
	for key, value := range transactionMap {
		response[key] = value
	}

	c.JSON(http.StatusOK, response)
}
