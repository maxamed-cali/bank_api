package models

import "time"

type User struct {
    ID          uint       `gorm:"primaryKey"`
    FullName    string     `gorm:"size:100;not null" `
    PhoneNumber string     `gorm:"size:15"  json:"phone_number"`
    Address     string     `json:"address"`
    IsActive    bool       `gorm:"default:true"`
    CreatedAt   time.Time
    Credential  *Credential `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
    Roles       []Role     `gorm:"many2many:user_roles;"`
}