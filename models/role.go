package models

import "time"

type Role struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"size:50;unique;not null"`
    CreatedAt time.Time
}