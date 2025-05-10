package models

import "time"

type Credential struct {
    ID           uint      `gorm:"primaryKey"`
    UserID       uint      `gorm:"uniqueIndex"` // one-to-one
    Username     string    `gorm:"size:50;unique;not null" json:"username"`
    PasswordHash string    `gorm:"size:255;not null"`
    CreatedAt    time.Time
    User         User      `gorm:"foreignKey:UserID"`
}