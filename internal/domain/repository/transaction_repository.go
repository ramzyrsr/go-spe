package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"go-spe/internal/domain/models"
	"go-spe/pkg/cache"
	"time"
)

type TransactionRepository interface {
	GetTransactionStatus(request_id, bill_number string) (*models.Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepositor(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// GetTransactionStatus implements TransactionRepository.
func (r *transactionRepository) GetTransactionStatus(request_id string, bill_number string) (*models.Transaction, error) {
	trx := &models.Transaction{}
	ctx := context.Background()

	// Check Redis cache first
	status, err := cache.RedisClient.Get(ctx, bill_number).Result()
	if status != "" {
		// Unmarshal the status from JSON to struct
		err = json.Unmarshal([]byte(status), trx)
		if err != nil {
			return nil, err
		}
		// Return cached transaction
		return trx, nil
	}

	// Query database if not found in cache
	row := r.db.QueryRow(`SELECT
			request_id,
			customer_pan,
			amount,
			transaction_datetime,
			rrn,
			bill_number,
			customer_name,
			merchant_id,
			merchant_name,
			merchant_city,
			currency_code,
			payment_status,
			payment_description 
		FROM transactions WHERE bill_number = $1`, bill_number)

	err = row.Scan(&trx.RequestID, &trx.CustomerPAN, &trx.Amount, &trx.TransactionDate, &trx.RRN,
		&trx.BillNumber, &trx.CustomerName, &trx.MerchantID, &trx.MerchantName, &trx.MerchantCity,
		&trx.CurrencyCode, &trx.PaymentStatus, &trx.PaymentDescription)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Cache the transaction struct as JSON for future requests
	trxJSON, err := json.Marshal(trx)
	if err != nil {
		return nil, err
	}

	// Set the transaction in Redis with a 5-minute expiration
	cache.RedisClient.Set(ctx, bill_number, trxJSON, time.Minute*5)

	return trx, nil
}
