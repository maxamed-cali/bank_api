package services

import (
	"bank/db"
	"bank/models"
	"errors"

	"github.com/lib/pq"
)

// List all users
func ListAllUsers() ([]models.User, error) {
	query := `SELECT id, full_name, email, phone_number, address FROM users`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.FullName, &user.Email, &user.PhoneNumber, &user.Address)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}



// Create a new role
func CreateRole(roleName string) (*models.Role, error) {
	query := `INSERT INTO roles (name, created_at) VALUES ($1, NOW()) RETURNING id`

	var role models.Role
	role.Name = roleName

	err := db.DB.QueryRow(query, roleName).Scan(&role.ID)
	if err != nil {
		return nil, err
	}

	return &role, nil
}


// Get all roles
func GetAllRoles() ([]models.Role, error) {
	query := `SELECT id, name FROM roles`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role

	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// Assign a role to a user
func AssignRolesToUser(userID uint, roleNames []string) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if user exists
	var userExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, userID).Scan(&userExists)
	if err != nil {
		return err
	}
	if !userExists {
		return errors.New("user not found")
	}

	// Find roles by names
	query := `SELECT id FROM roles WHERE name = ANY($1)`
	rows, err := tx.Query(query, pq.Array(roleNames))
	if err != nil {
		return err
	}
	defer rows.Close()

	var roleIDs []uint
	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			return err
		}
		roleIDs = append(roleIDs, id)
	}
	if len(roleIDs) == 0 {
		return errors.New("roles not found")
	}

	// Remove existing roles
	_, err = tx.Exec(`DELETE FROM user_roles WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}

	// Assign new roles
	stmt, err := tx.Prepare(`INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, roleID := range roleIDs {
		if _, err := stmt.Exec(userID, roleID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

