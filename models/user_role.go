package models

type UserRole struct {
    ID     uint `gorm:"primaryKey"`
    UserID uint
    RoleID uint
}
