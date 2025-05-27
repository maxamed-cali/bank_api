package models

import "time"

type Account struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	UserID        uint        `json:"user_id"`
	AccountNumber string      `gorm:"unique;not null" json:"account_number"`
	Balance       float64     `gorm:"type:decimal(15,2);default:0.00" json:"balance"`
	AccountTypeID uint        `json:"account_type_id"`
	IsActive      bool        `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time   `gorm:"autoCreateTime" json:"created_at"`
	AccountType   AccountType `gorm:"foreignKey:AccountTypeID" json:"account_type"`
}