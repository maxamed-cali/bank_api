package services

import (
	"bank/db"
	"bank/models"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Register user, credentials, and role
func RegisterUser(input models.User, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := db.DB.Create(&input).Error; err != nil {
		return err
	}

	cred := models.Credential{
		UserID:       input.ID,
		Username:     username,
		PasswordHash: string(hashedPassword),
	}
	if err := db.DB.Create(&cred).Error; err != nil {
		return err
	}

	var role models.Role
	db.DB.Where("name = ?", "User").FirstOrCreate(&role, models.Role{Name: "User"})

	return db.DB.Create(&models.UserRole{
		UserID: input.ID,
		RoleID: role.ID,
	}).Error
}

// Authenticate username & password and return userID + role
func Authenticate(username, password string) (uint, string, error) {
	var cred models.Credential
	if err := db.DB.Where("username = ?", username).First(&cred).Error; err != nil {
		return 0, "", errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(cred.PasswordHash), []byte(password)); err != nil {
		return 0, "", errors.New("invalid username or password")
	}

	var userRole models.UserRole
	if err := db.DB.Where("user_id = ?", cred.UserID).First(&userRole).Error; err != nil {
		return 0, "", errors.New("user role not found")
	}

	var role models.Role
	if err := db.DB.First(&role, userRole.RoleID).Error; err != nil {
		return 0, "", errors.New("role not found")
	}

	return cred.UserID, role.Name, nil
}

// Reset password
func ResetUserPassword(userID uint, oldPwd, newPwd string) error {
	var cred models.Credential
	if err := db.DB.Where("user_id = ?", userID).First(&cred).Error; err != nil {
		return gorm.ErrRecordNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(cred.PasswordHash), []byte(oldPwd)); err != nil {
		return errors.New("old password incorrect")
	}

	hashedNew, _ := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	cred.PasswordHash = string(hashedNew)

	return db.DB.Save(&cred).Error
}

// Get user profile
func GetUserByID(userID uint) (models.User, error) {
	var user models.User
	err := db.DB.First(&user, userID).Error
	return user, err
}