package services

import (
	"bank/db"
	"bank/models"
	"errors"
)

// Create a new account
func CreateAccount(acc *models.Account) error {
	query := `INSERT INTO accounts (account_number, balance, user_id, account_type_id, created_at)
	          VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
	return db.DB.QueryRow(query, acc.AccountNumber, acc.Balance, acc.UserID, acc.AccountTypeID).
		Scan(&acc.ID)
}

// Get all accounts with their account type names (simulating Preload)
func GetAllAccounts() ([]models.Account, error) {
	query := `SELECT a.id, a.account_number, a.balance, a.user_id, a.account_type_id,  at.type_name AS account_type_name, a.created_at
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

		err := rows.Scan(&acc.ID, &acc.AccountNumber, &acc.Balance, &acc.UserID, &acc.AccountTypeID,
			&accountTypeName, &acc.CreatedAt)
		if err != nil {
			return nil, err
		}

		acc.AccountType = models.AccountType{TypeName: accountTypeName}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

// Get a single account by ID (with account type)
func GetAccountByID(id uint) (*models.Account, error) {
	query := `SELECT a.id, a.account_number, a.balance, a.user_id, a.account_type_id,  at.type_name AS account_type_name, a.created_at
	          FROM accounts a
	          JOIN account_types at ON a.account_type_id = at.id
	          WHERE a.id = $1`

	var acc models.Account
	var accountTypeName string
	err := db.DB.QueryRow(query, id).Scan(&acc.ID, &acc.AccountNumber, &acc.Balance, &acc.UserID, &acc.AccountTypeID,
		&accountTypeName, &acc.CreatedAt)

	if err != nil {
		return nil, err
	}

	acc.AccountType = models.AccountType{TypeName: accountTypeName}
	return &acc, nil
}



// Update an account
func UpdateAccount(id uint, updated *models.Account) error {
	query := `UPDATE accounts 
	          SET account_number = $1, balance = $2, user_id = $3, account_type_id = $4
	          WHERE id = $5`
	result, err := db.DB.Exec(query, updated.AccountNumber, updated.Balance, updated.UserID, updated.AccountTypeID, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no record updated")
	}
	return nil
}

// Delete an account
func DeleteAccount(id uint) error {
	query := `DELETE FROM accounts WHERE id = $1`
	result, err := db.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no record deleted")
	}
	return nil
}
