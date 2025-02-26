package auth

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	RequestID  string `json:"request_id"`
	RRN        string `json:"rrn"`
	BillNumber string `json:"bill_number"`
	MerchantID string `json:"merchant_id"`
}

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
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
			return
		}

		var reqBody RequestBody
		err = json.Unmarshal(body, &reqBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON"})
			return
		}
		payload := fmt.Sprintf("%s:%s:%s", reqBody.RequestID, reqBody.RRN, reqBody.MerchantID)

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
