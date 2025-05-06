package services

import (
	"bank/models"
	"bank/db"
)

func CreateAccountType(data *models.AccountType) error {
	return db.DB.Create(data).Error
}

func GetAllAccountTypes() ([]models.AccountType, error) {
	var list []models.AccountType
	err := db.DB.Find(&list).Error
	return list, err
}

func GetAccountTypeByID(id uint) (*models.AccountType, error) {
	var a models.AccountType
	err := db.DB.First(&a, id).Error
	return &a, err
}

func UpdateAccountType(id uint, update models.AccountType) error {
	return db.DB.Model(&models.AccountType{}).Where("id = ?", id).Updates(update).Error
}

func DeleteAccountType(id uint) error {
	return db.DB.Delete(&models.AccountType{}, id).Error
}
