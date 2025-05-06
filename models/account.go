package models

import "time"

type Account struct {
	ID            uint        `gorm:"primaryKey"`
	UserID        uint        
	AccountNumber string      `gorm:"unique;not null"`
	Balance       float64     `gorm:"default:0"`
	AccountTypeID uint        `gorm:"not null"`
	IsActive      bool        `gorm:"default:true"`
	CreatedAt     time.Time
	AccountType   AccountType `gorm:"foreignKey:AccountTypeID"`
}
