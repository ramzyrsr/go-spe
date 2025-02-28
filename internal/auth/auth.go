package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-spe/internal/domain/models"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Generate HMAC SHA512 and encode in base64
func generateSignature(payload string, key string) string {
	h := hmac.New(sha512.New, []byte(key))
	h.Write([]byte(payload))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Middleware to validate X-Signature
func SignatureMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}

		// Extract X-Signature Header
		signature := c.GetHeader("X-Signature")
		if signature == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing X-Signature header"})
			c.Abort()
			return
		}

		// Prepare the payload value
		var reqBody *models.Transaction
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Reset the request body so it can be read again
		body, _ := json.Marshal(reqBody)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Check the route path to determine the payload
		var payload string
		switch c.FullPath() {
		case "/api/v1/transaction-notification":
			payload = fmt.Sprintf("%s:%s:%s", reqBody.RequestID, reqBody.RRN, reqBody.MerchantID)
		case "/api/v1/check-status":
			payload = fmt.Sprintf("%s", reqBody.BillNumber)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route"})
			c.Abort()
			return
		}

		// Generate expected signature
		expectedSignature := generateSignature(payload, os.Getenv("SECRET_KEY"))
		fmt.Println(expectedSignature)

		// Compare expected signature with received signature
		if signature != expectedSignature {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			c.Abort()
			return
		}

		// If valid, proceed
		c.Next()
	}
}
