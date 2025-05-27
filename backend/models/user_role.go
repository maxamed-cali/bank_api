package models

type UserRole struct {
    ID     uint `gorm:"primaryKey"`
    UserID uint
    RoleID uint
    User User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
    Role Role `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

