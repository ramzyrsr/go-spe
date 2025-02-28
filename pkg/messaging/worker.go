package messaging

import (
	"encoding/json"
	"fmt"
	"go-spe/internal/domain/models"
	"go-spe/pkg/db"
	"log"
	"strconv"
	"time"
)

func ProcessTransactions() {
	msgs, err := Channel.Consume(
		Queue.Name, "", true, false, false, false, nil,
	)
	if err != nil {
		log.Fatal("Failed to consume from queue:", err)
	}

	for msg := range msgs {
		var trx *models.Transaction
		err := json.Unmarshal(msg.Body, &trx)
		if err != nil {
			log.Println("Failed to parse transaction:", err)
			continue
		}

		// Save to database
		trx.TransactionDate = time.Now()
		floatAmount, _ := strconv.ParseFloat(trx.Amount, 64)

		query := `INSERT INTO transactions (
			request_id, customer_pan, amount, transaction_datetime, rrn, 
			bill_number, customer_name, merchant_id, merchant_name, merchant_city, 
			currency_code, payment_status, payment_description
		)
		VALUES(
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, 
			$11, $12, $13
		);`

		// Connect to the database
		dbConn, err := db.Connect()
		if err != nil {
			log.Println("Failed to connect to database:", err)
			continue
		}
		defer dbConn.Close()

		// Execute the query
		_, err = dbConn.Exec(query, trx.RequestID, trx.CustomerPAN, floatAmount, trx.TransactionDate, trx.RRN,
			trx.BillNumber, trx.CustomerName, trx.MerchantID, trx.MerchantName, trx.MerchantCity,
			trx.CurrencyCode, trx.PaymentStatus, trx.PaymentDescription)

		if err != nil {
			log.Println("Failed to save transaction:", err)
		} else {
			fmt.Println("Transaction saved:", trx.RequestID)
		}
	}
}
