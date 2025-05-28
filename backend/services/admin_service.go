package services

import (
	"bank/db"
	"bank/dtos"
	"bank/models"
	"errors"
	"fmt"

	"github.com/lib/pq"
)



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


// ToggleUserStatus updates the active status of a user.
func ActivateDeactivateUser(userID uint, isActive bool) error {
	query := `UPDATE users SET is_active = $1 WHERE id = $2`
	result, err := db.GetDB().Exec(query, isActive, userID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", userID)
	}

	// Log the action
	status := "deactivated"
	if isActive {
		status = "activated"
	}
	_ = LogAudit(&userID, "UPDATE", "users", userID, fmt.Sprintf("User %d %s", userID, status))

	return nil
}

func GetUsersWithRoles() ([]dtos.UserWithRoleDTO, error) {
    query := `
        SELECT u.id, u.full_name, u.phone_number, r.name, u.is_active, u.created_at
        FROM users u
        JOIN user_roles ur ON u.id = ur.user_id
        JOIN roles r ON ur.role_id = r.id
    `
    rows, err :=db.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []dtos.UserWithRoleDTO
    for rows.Next() {
        var user dtos.UserWithRoleDTO
        if err := rows.Scan(&user.ID,&user.Name, &user.Phone, &user.Role, &user.Status, &user.CreatedAt); err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    return users, nil
}

func GetAdminDashboardSummary() (*dtos.DashboardSummary, error) {
	var summary dtos.DashboardSummary

	// 1. Total Wallet Balance
	err := db.DB.QueryRow(`
		SELECT COALESCE(SUM(balance), 0)
		FROM accounts
	`).Scan(&summary.WalletBalance)
	if err != nil {
		return nil, err
	}

	// 2. Total Transactions
	err = db.DB.QueryRow(`
		SELECT COUNT(*)
		FROM transactions
	`).Scan(&summary.TotalTransactions)
	if err != nil {
		return nil, err
	}

	// 3. Pending Money Requests
	err = db.DB.QueryRow(`
		SELECT COUNT(*)
		FROM money_requests
		WHERE status = 'pending'
	`).Scan(&summary.PendingRequests)
	if err != nil {
		return nil, err
	}

	// 4. Total Transfers (DEBIT)
	err = db.DB.QueryRow(`
		SELECT COUNT(*)
		FROM transactions
		WHERE transaction_type = 'DEBIT'
	`).Scan(&summary.TotalTransfers)
	if err != nil {
		return nil, err
	}

	// 5. Total Sent Amount
	err = db.DB.QueryRow(`
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE transaction_type = 'DEBIT'
	`).Scan(&summary.TotalSentAmount)
	if err != nil {
		return nil, err
	}

	// 6. Total Received Amount
	err = db.DB.QueryRow(`
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE transaction_type = 'CREDIT'
	`).Scan(&summary.TotalReceivedAmount)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}
