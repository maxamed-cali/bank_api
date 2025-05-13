package services

import (
	"bank/db"
	"bank/models"
	"database/sql"
	"errors"
	"fmt"
)

func CreateAccountType(at *models.AccountType) error {
	query := `INSERT INTO account_types (type_name, description, currency)
	          VALUES ($1, $2, $3) RETURNING id`
	err := db.DB.QueryRow(query, at.TypeName, at.Description, at.Currency).Scan(&at.ID)
	if err != nil {
		return err
	}

	// Log audit for create action
	_ = LogAudit(&at.ID, "CREATE", "account_types", at.ID, "Account type created")

	return nil
}
func GetAllAccountTypes() ([]*models.AccountType, error) {
	query := `SELECT id, type_name, description, currency FROM account_types`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accountTypes []*models.AccountType

	for rows.Next() {
		var at models.AccountType
		err := rows.Scan(&at.ID, &at.TypeName, &at.Description, &at.Currency)
		if err != nil {
			return nil, err
		}
		accountTypes = append(accountTypes, &at)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accountTypes, nil
}

func GetAccountTypeByID(id uint) (*models.AccountType, error) {
	query := `SELECT id, type_name, description, currency FROM account_types WHERE id = $1`

	var at models.AccountType
	err := db.DB.QueryRow(query, id).Scan(&at.ID, &at.TypeName, &at.Description, &at.Currency)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account type with ID %d not found", id)
		}
		return nil, err
	}

	return &at, nil
}

func UpdateAccountType(id uint, updated *models.AccountType) error {
	query := `UPDATE account_types SET type_name = $1, description = $2, currency = $3 WHERE id = $4`
	result, err := db.DB.Exec(query, updated.TypeName, updated.Description, updated.Currency, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no record updated")
	}

	// Log audit for update action
	_ = LogAudit(nil, "UPDATE", "account_types", id, "Account type updated")

	return nil
}

func DeleteAccountType(id uint) error {
	query := `DELETE FROM account_types WHERE id = $1`
	result, err := db.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no record deleted")
	}

	// Log audit for delete action
	_ = LogAudit(nil, "DELETE", "account_types", id, "Account type deleted")

	return nil
}
