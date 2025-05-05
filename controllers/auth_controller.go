package controllers
import (
    "net/http"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "bank/db"
    "bank/models"
    "bank/utils"
)


type RegisterInput struct {
    FullName string `json:"full_name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Phone    string `json:"phone_number"`
    Address  string `json:"address"`
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}


func Register(c *gin.Context) {
    var input RegisterInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Hash password
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

    user := models.User{
        FullName:    input.FullName,
        Email:       input.Email,
        PhoneNumber: input.Phone,
        Address:     input.Address,
    }

    result := db.DB.Create(&user)
    if result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
        return
    }

    credential := models.Credential{
        UserID:       user.ID,
        Username:     input.Username,
        PasswordHash: string(hashedPassword),
    }
    db.DB.Create(&credential)

    // Assign default role 'User'
    var userRole models.Role
    db.DB.Where("name = ?", "User").FirstOrCreate(&userRole, models.Role{Name: "User"})
    db.DB.Create(&models.UserRole{
        UserID: user.ID,
        RoleID: userRole.ID,
    })

    c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}


type LoginInput struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
    var input LoginInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var credential models.Credential
    if err := db.DB.Where("username = ?", input.Username).First(&credential).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(credential.PasswordHash), []byte(input.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

	var userRole models.UserRole
	db.DB.Where("user_id = ?", credential.ID).First(&userRole)

	var role models.Role
	db.DB.First(&role, userRole.RoleID)

    token, err := utils.GenerateJWT(credential.UserID, role.Name)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}
