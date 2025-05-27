package models

import "time"

type Transaction struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint       `json:"user_id"`  // Foreign Key to Users
	AccountID       string      `json:"account_id"`
	TransactionType string    `gorm:"type:varchar(20);not null" json:"transaction_type"` 
	ToAccountID     *string  `json:"to_account_id,omitempty"`  // Destination account (only for TRANSFER) // DEPOSIT, WITHDRAW, TRANSFER
	Amount          float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	TransactionDate time.Time `gorm:"autoCreateTime" json:"transaction_date"`
	Description     string    `json:"description"`
}