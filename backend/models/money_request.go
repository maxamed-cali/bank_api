package models

import "time"

type MoneyRequest struct {
	ID          uint       `gorm:"primaryKey"`
    UserID       uint      `gorm:"uniqueIndex" json:"user_id"` // one-to-one
	RequesterID string      `gorm:"not null" json:"requester_id"` // who is requesting the money
	RecipientID string       `gorm:"not null" json:"recipient_id"`  // who is being asked to send money
	Amount      float64    `gorm:"not null"`
	Status      string     // PENDING, ACCEPTED, DECLINED, EXPIRED
	ExpiresAt   time.Time  // auto-expiry time       // retry count
	RequesteAt  time.Time // last retry attempt
	RecipientUserID uint	   `gorm:"uniqueIndex" json:"recipient_user_id"` // ID of the user who is being asked to send money
	//  User         User      `gorm:"foreignKey:UserID"`


}
