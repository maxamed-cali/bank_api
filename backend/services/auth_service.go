package services

import (
	"bank/db"
	"bank/models"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Register user, credentials, and role
func RegisterUser(input models.User, email, password string) error {
    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    // Begin a transaction
    tx, err := db.DB.Begin()
    if err != nil {
        return err
    }

    defer func() {
        if err != nil {
            tx.Rollback()
        } else {
            tx.Commit()
        }
    }()

    // Insert into users table and get ID
    var userID int
    query := `INSERT INTO users (full_name,  phone_number, address,created_at) VALUES ($1, $2, $3,$4) RETURNING id`
    err = tx.QueryRow(query, input.FullName,  input.PhoneNumber,input.Address,time.Now()).Scan(&userID)
    if err != nil {
        return fmt.Errorf("failed to insert user: %v", err)
    }

    // Insert into credentials table
    _, err = tx.Exec(
        `INSERT INTO credentials (user_id, email, password_hash,created_at) VALUES ($1, $2, $3,$4)`,
        userID, email, string(hashedPassword),time.Now(),
    )
    if err != nil {
     if pqErr, ok := err.(*pq.Error); ok {
            if pqErr.Code == "23505" && strings.Contains(pqErr.Constraint, "credentials_email_key") {
                return errors.New("an account with this email already exists")
            }
        }
        return fmt.Errorf("failed to insert credentials: %w", err)
    }

    // Check if "User" role exists, if not create it
    var roleID int
    err = tx.QueryRow(`SELECT id FROM roles WHERE name = 'User'`).Scan(&roleID)
    if err == sql.ErrNoRows {
        err = tx.QueryRow(
    `INSERT INTO roles (name, created_at) VALUES ($1, $2) RETURNING id`,
    "User",
    time.Now(),
).Scan(&roleID)
        if err != nil {
            return fmt.Errorf("failed to insert default role: %v", err)
        }
    } else if err != nil {
        return fmt.Errorf("failed to query role: %v", err)
    }

    // Insert into user_roles table
    _, err = tx.Exec(`INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`, userID, roleID)
    if err != nil {
        return fmt.Errorf("failed to assign user role: %v", err)
    }

    return nil
}

// Authenticate email & password and return userID + role
func Authenticate(email, password string) (uint, string, error) {
    var userID uint
    var passwordHash string

    // Step 1: Find user credentials
    err := db.DB.QueryRow(`
        SELECT user_id, password_hash FROM credentials WHERE email = $1
    `, email).Scan(&userID, &passwordHash)
    if err == sql.ErrNoRows {
        return 0, "", errors.New("invalid credentials Of email")
    } else if err != nil {
        return 0, "", err
    }

    // Step 2: Compare password
    if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
        return 0, "", errors.New("invalid credentials of password")
    }

    // Step 3: Find user role
    var roleID uint
    err = db.DB.QueryRow(`
        SELECT role_id FROM user_roles WHERE user_id = $1
    `, userID).Scan(&roleID)
    if err == sql.ErrNoRows {
        return 0, "", errors.New("user role not found")
    } else if err != nil {
        return 0, "", err
    }

    // Step 4: Get role name
    var roleName string
    err = db.DB.QueryRow(`
        SELECT name FROM roles WHERE id = $1
    `, roleID).Scan(&roleName)
    if err == sql.ErrNoRows {
        return 0, "", errors.New("role not found")
    } else if err != nil {
        return 0, "", err
    }

    return userID, roleName, nil
}

// Reset password
func ResetUserPassword(userID uint, oldPwd, newPwd string) error {
    // Step 1: Retrieve the current password hash for the user
    var passwordHash string
    err := db.DB.QueryRow(`
        SELECT password_hash FROM credentials WHERE user_id = $1
    `, userID).Scan(&passwordHash)
    if err == sql.ErrNoRows {
        return errors.New("user not found")
    } else if err != nil {
        return err
    }

    // Step 2: Compare the provided old password with the stored password hash
    if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(oldPwd)); err != nil {
        return errors.New("old password incorrect")
    }

    // Step 3: Hash the new password
    hashedNew, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    // Step 4: Update the password hash in the database
    _, err = db.DB.Exec(`
        UPDATE credentials
        SET password_hash = $1
        WHERE user_id = $2
    `, string(hashedNew), userID)
    if err != nil {
        return err
    }

    return nil
}
// Get user profile
func GetUserByID(userID uint) (models.User, error) {
    var user models.User

    row := db.DB.QueryRow(`
        SELECT id, full_name,  phone_number, created_at 
        FROM users
        WHERE id = $1`, userID)

    err := row.Scan(
        &user.ID,
        &user.FullName,
      
        &user.PhoneNumber,
        &user.CreatedAt,
       
    )

    if err == sql.ErrNoRows {
        return user, errors.New("user not found")
    } else if err != nil {
        return user, err
    }

    return user, nil
}