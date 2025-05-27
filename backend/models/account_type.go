package models

type AccountType struct {
	ID          uint   `gorm:"primaryKey"`
	TypeName    string `gorm:"unique;not null" json:"type_name"`
	Description string
	Currency    string `gorm:"not null" json:"currency"`
}