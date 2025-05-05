package models

import "time"

type Credential struct {
    ID           uint      `gorm:"primaryKey"`
    UserID       uint      `gorm:"uniqueIndex"`
    Username     string    `gorm:"size:50;unique;not null"`
    PasswordHash string    `gorm:"size:255;not null"`
    CreatedAt    time.Time
}
