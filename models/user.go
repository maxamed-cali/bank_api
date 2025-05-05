package models

import "time"

type User struct {
    ID          uint      `gorm:"primaryKey"`
    FullName    string    `gorm:"size:100;not null"`
    Email       string    `gorm:"size:100;unique;not null"`
    PhoneNumber string    `gorm:"size:15"`
    Address     string
    CreatedAt   time.Time
    Credential  Credential
    Roles       []Role `gorm:"many2many:user_roles;"`
}