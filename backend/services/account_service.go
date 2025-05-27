package services

import (
	"bank/db"
	"bank/dtos"
	"bank/models"
	"errors"

	"github.com/lib/pq"
)


func CreateAccount(acc *models.Account) error {
	query := `INSERT INTO accounts (account_number, balance, user_id, account_type_id, created_at)
	          VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
	err := db.GetDB().QueryRow(query, acc.AccountNumber, acc.Balance, acc.UserID, acc.AccountTypeID).
		Scan(&acc.ID)
	if err != nil {
		// Check for PostgreSQL unique constraint violation
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			if pgErr.Constraint == "accounts_account_number_key" {
				return errors.New("account number already exists")
			}
			return errors.New("duplicate value violates a unique constraint")
		}
		return err
	}

	// Log audit for create
	_ = LogAudit(&acc.UserID, "CREATE", "accounts", acc.ID, "Account created")
	return nil
}


// Get balance for a specific account by account ID
func GetAccountBalance(id string) (float64, error) {
	var balance float64
	query := `SELECT balance FROM accounts WHERE account_number = $1`
	err := db.GetDB().QueryRow(query, id).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}
// Get all accounts with their account type names (simulating Preload)
func GetAllAccounts() ([]models.Account, error) {
	query := `SELECT a.id, a.account_number, a.balance, a.user_id, a.account_type_id,  at.type_name AS account_type_name,at.description, currency,a.created_at
	          FROM accounts a
	          JOIN account_types at ON a.account_type_id = at.id`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var acc models.Account
		var accountTypeName string
		var description string
		var currency string

		err := rows.Scan(&acc.ID, &acc.AccountNumber, &acc.Balance, &acc.UserID, &acc.AccountTypeID,
			&accountTypeName,&description,&currency, &acc.CreatedAt)
		if err != nil {
			return nil, err
		}

		acc.AccountType = models.AccountType{TypeName: accountTypeName,Description: description, Currency: currency}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

// Get a single account by ID (with account type)
func GetAccountsByUserID(userID uint) ([]dtos.AccountResponse, error) {
	query := `SELECT a.id, a.account_number, a.balance, a.user_id, a.account_type_id,
	                 at.type_name, at.description, at.currency,
	                 u.full_name, a.created_at
	          FROM accounts a
	          JOIN account_types at ON a.account_type_id = at.id
	          JOIN users u ON a.user_id = u.id
	          WHERE a.user_id = $1`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return []dtos.AccountResponse{}, err // return empty slice instead of nil
	}
	defer rows.Close()

	responses := make([]dtos.AccountResponse, 0) // Always create a slice

	for rows.Next() {
		var acc dtos.AccountResponse
		err := rows.Scan(
			&acc.ID, &acc.AccountNumber, &acc.Balance, &acc.UserID, &acc.AccountTypeID,
			&acc.TypeName, &acc.Description, &acc.Currency,
			&acc.Name, &acc.CreatedAt,
		)
		if err != nil {
			return []dtos.AccountResponse{}, err
		}
		responses = append(responses, acc)
	}

	// Optionally check if no rows were found
	if len(responses) == 0 {
		return responses, nil // or custom error/info message if needed
	}

	return responses, nil
}




// Update an account
func UpdateAccount(id uint, updated *models.Account) error {
	query := `UPDATE accounts 
	          SET account_number = $1, balance = $2, user_id = $3, account_type_id = $4
	          WHERE id = $5`
	result, err := db.GetDB().Exec(query, updated.AccountNumber, updated.Balance, updated.UserID, updated.AccountTypeID, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no record updated")
	}

	// Log audit for update
	_ = LogAudit(&updated.UserID, "UPDATE", "accounts", id, "Account updated")
	return nil
}


// Delete an account
func DeleteAccount(id uint) error {
	query := `DELETE FROM accounts WHERE id = $1`
	result, err := db.GetDB().Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no record deleted")
	}

	// Log audit for delete 
	_ = LogAudit(nil, "DELETE", "accounts", id, "Account deleted")
	return nil
}


func GetUserByAccountNumber( accountNumber string) (*dtos.AccoutResponse, error) {
    var user dtos.AccoutResponse
    query := `SELECT u.full_name, a.account_number
        FROM users u
        JOIN accounts a ON u.id = a.user_id
		where a.account_number = $1
        LIMIT 1`
    err :=db.GetDB().QueryRow(query, accountNumber).Scan( &user.FullName,&user.AccountNumber,)
    if err != nil {
        return nil, err
    }
    return &user, nil
}