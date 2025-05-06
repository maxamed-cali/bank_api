package models

type AccountType struct {
	ID          uint   `gorm:"primaryKey"`
	TypeName    string `gorm:"unique;not null"`
	Description string
	Currency    string `gorm:"not null"`
}