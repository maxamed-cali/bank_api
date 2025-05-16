package models

import "time"

type Credential struct {
    ID           uint      `gorm:"primaryKey"`
    UserID       uint      `gorm:"uniqueIndex"` // one-to-one
    Email       string     `gorm:"size:100;unique;not null" ` 
    PasswordHash string    `gorm:"size:255;not null"`
    CreatedAt    time.Time
    User         User      `gorm:"foreignKey:UserID"`
}