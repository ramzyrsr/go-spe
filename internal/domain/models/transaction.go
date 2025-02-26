package models

import "time"

// Transaction represents a payment transaction entity
type Transaction struct {
	RequestID          string    `json:"request_id"`
	CustomerPAN        string    `json:"customer_pan"`
	Amount             float64   `json:"amount"`
	TransactionDate    time.Time `json:"transaction_date"`
	MerchantID         string    `json:"merchant_id"`
	MerchantName       string    `json:"merchant_name"`
	MerchantCity       string    `json:"merchant_city"`
	CurrencyCode       string    `json:"currency_code"`
	PaymentStatus      string    `json:"payment_status"`
	PaymentDescription string    `json:"payment_description"`
}
