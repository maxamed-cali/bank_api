package services

import (
	"bank/db"
	"bank/models"
	"bank/websocket"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)



func MoneyTransfer(tx *models.Transaction) error {
	if tx.ToAccountID == nil {
		return errors.New("missing destination account for transfer")
	}

	dbtx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			dbtx.Rollback()
			panic(p)
		}
	}()

	var sender, receiver models.Account

	// Fetch sender
	err = dbtx.QueryRow(`SELECT id, account_number, balance, is_active, user_id FROM accounts WHERE account_number = $1`, tx.AccountID).
		Scan(&sender.ID, &sender.AccountNumber, &sender.Balance, &sender.IsActive, &sender.UserID)
	if err != nil {
		dbtx.Rollback()
		return fmt.Errorf("sender account not found: %v", err)
	}

	// Fetch receiver
	err = dbtx.QueryRow(`SELECT id, account_number, balance, is_active, user_id FROM accounts WHERE account_number = $1`, *tx.ToAccountID).
		Scan(&receiver.ID, &receiver.AccountNumber, &receiver.Balance, &receiver.IsActive, &receiver.UserID)
	if err != nil {
		dbtx.Rollback()
		return fmt.Errorf("receiver account not found: %v", err)
	}

	if sender.AccountNumber == receiver.AccountNumber {
		dbtx.Rollback()
		return errors.New("cannot transfer to self")
	}
	if !sender.IsActive || !receiver.IsActive {
		dbtx.Rollback()
		return errors.New("both accounts must be active")
	}
	if sender.Balance < tx.Amount {
		dbtx.Rollback()
		return errors.New("insufficient balance")
	}

	// Update balances
	_, err = dbtx.Exec(`UPDATE accounts SET balance = balance - $1 WHERE account_number = $2`, tx.Amount, sender.AccountNumber)
	if err != nil {
		dbtx.Rollback()
		return err
	}
	_, err = dbtx.Exec(`UPDATE accounts SET balance = balance + $1 WHERE account_number = $2`, tx.Amount, receiver.AccountNumber)
	if err != nil {
		dbtx.Rollback()
		return err
	}

	// Insert sender transaction
	_, err = dbtx.Exec(`INSERT INTO transactions (account_id, to_account_id, transaction_type, amount, description, transaction_date)
	                    VALUES ($1, $2, 'DEBIT', $3, $4, NOW())`,
		sender.AccountNumber, receiver.AccountNumber, tx.Amount,
		fmt.Sprintf("Transferred to Account ID %s", receiver.AccountNumber))
	if err != nil {
		dbtx.Rollback()
		return err
	}

	// Insert receiver transaction
	_, err = dbtx.Exec(`INSERT INTO transactions (account_id, to_account_id, transaction_type, amount, description, transaction_date)
	                    VALUES ($1, $2, 'CREDIT', $3, $4, NOW())`,
		receiver.AccountNumber, sender.AccountNumber, tx.Amount,
		fmt.Sprintf("Received from Account ID %s", sender.AccountNumber))
	if err != nil {
		dbtx.Rollback()
		return err
	}

	// Insert notification
	message := fmt.Sprintf("You received %.2f from %s", tx.Amount, sender.AccountNumber)
	_, err = dbtx.Exec(`INSERT INTO notifications (user_id, message, created_at) VALUES ($1, $2, NOW())`,
		receiver.UserID, message)
	if err != nil {
		dbtx.Rollback()
		return err
	}

	// Commit transaction
	if err := dbtx.Commit(); err != nil {
		return err
	}

	// Send real-time notification (non-DB)
	websocket.NotifyChan <- websocket.NotificationMessage{
		UserID:  receiver.UserID,
		Message: message,
	}

	return nil
}


func MoneyRequest(request *models.MoneyRequest) error {
	if request.Amount <= 0 {
		return errors.New("invalid amount")
	}

	if request.RequesterID == request.RecipientID {
		return errors.New("cannot request from self")
	}

	request.Status = "PENDING"
	request.ExpiresAt = time.Now().Add(24 * time.Hour)

	dbtx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			dbtx.Rollback()
			panic(p)
		}
	}()

	// Insert money request
	_, err = dbtx.Exec(`
		INSERT INTO money_requests 
			(requester_id, recipient_id, amount, status, expires_at, requeste_at)
		VALUES 
			($1, $2, $3, $4,  NOW(), NOW())
	`, request.RequesterID, request.RecipientID, request.Amount, request.Status)
	if err != nil {
		dbtx.Rollback()
		return err
	}

	// Get recipient account
	var recipientUserID uint
	err = dbtx.QueryRow(`
		SELECT user_id FROM accounts WHERE account_number = $1
	`, request.RecipientID).Scan(&recipientUserID)
	if err != nil {
		dbtx.Rollback()
		return fmt.Errorf("recipient account not found: %v", err)
	}

	// Commit the transaction
	if err := dbtx.Commit(); err != nil {
		return err
	}

	// Send WebSocket notification
	message := fmt.Sprintf("User %v requested %.2f from you", request.RequesterID, request.Amount)
	websocket.NotifyChan <- websocket.NotificationMessage{
		UserID:  recipientUserID,
		Message: message,
	}

	return nil
}

func AcceptMoneyRequest(requestID uint) error {
	dbtx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			dbtx.Rollback()
			panic(p)
		}
	}()

	// Fetch the money request
	var req models.MoneyRequest
	err = dbtx.QueryRow(`
		SELECT id, requester_id, recipient_id, amount, status 
		FROM money_requests 
		WHERE id = $1
	`, requestID).Scan(&req.ID, &req.RequesterID, &req.RecipientID, &req.Amount, &req.Status)
	if err != nil {
		dbtx.Rollback()
		return fmt.Errorf("money request not found: %v", err)
	}

	if req.Status != "PENDING" {
		dbtx.Rollback()
		return errors.New("request is no longer active")
	}

	// Update status to ACCEPTED
	_, err = dbtx.Exec(`
		UPDATE money_requests 
		SET status = 'ACCEPTED'
		WHERE id = $1
	`, requestID)
	if err != nil {
		dbtx.Rollback()
		return err
	}

	// Commit before calling MoneyTransfer
	if err := dbtx.Commit(); err != nil {
		return err
	}

	// Use existing transfer logic (can be raw SQL or GORM-based)
	tx := &models.Transaction{
		AccountID:   req.RecipientID,
		ToAccountID: &req.RequesterID,
		Amount:      req.Amount,
		Description: fmt.Sprintf("Accepted request ID %d", req.ID),
	}
	return MoneyTransfer(tx)
}


func DeclineMoneyRequest(requestID uint) error {
	dbtx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			dbtx.Rollback()
			panic(p)
		}
	}()

	// Fetch the money request
	var req models.MoneyRequest
	err = dbtx.QueryRow(`
		SELECT id, requester_id, recipient_id, amount, status
		FROM money_requests
		WHERE id = $1
	`, requestID).Scan(&req.ID, &req.RequesterID, &req.RecipientID, &req.Amount, &req.Status)
	if err != nil {
		dbtx.Rollback()
		return fmt.Errorf("money request not found: %v", err)
	}

	if req.Status != "PENDING" {
		dbtx.Rollback()
		return errors.New("request is no longer active")
	}

	// Update status to DECLINED
	_, err = dbtx.Exec(`
		UPDATE money_requests 
		SET status = 'DECLINED'
		WHERE id = $1
	`, requestID)
	if err != nil {
		dbtx.Rollback()
		return err
	}

	// Get requester account (for WebSocket notification)
	var userID uint
	err = dbtx.QueryRow(`
		SELECT user_id FROM accounts WHERE account_number = $1
	`, req.RequesterID).Scan(&userID)
	if err != nil {
		dbtx.Rollback()
		return fmt.Errorf("requester account not found: %v", err)
	}

	if err := dbtx.Commit(); err != nil {
		return err
	}

	// Send WebSocket notification
	message := fmt.Sprintf("Your money request (ID %d) was declined", req.ID)
	websocket.NotifyChan <- websocket.NotificationMessage{
		UserID:  userID,
		Message: message,
	}

	return nil
}
	

func AutoExpireRequests() {
	dbConn := db.DB

	fmt.Println("Checking for expired requests...", time.Now())

	// Step 1: Get all PENDING requests using raw SQL
	rows, err := dbConn.Query(`
		SELECT id, requester_id, recipient_id, amount, status
		FROM money_requests 
		WHERE status = 'PENDING' 
	`)
	if err != nil {
		log.Println("Failed to fetch expired requests:", err)
		return
	}
	defer rows.Close()

	// Step 2: Process each request
	for rows.Next() {
		var req models.MoneyRequest
		if err := rows.Scan(&req.ID, &req.RequesterID, &req.RecipientID, &req.Amount, &req.Status); err != nil {
			log.Printf("Failed to scan request: %v\n", err)
			continue
		}

		// Step 2a: Mark as EXPIRED
		_, err := dbConn.Exec(`
			UPDATE money_requests 
			SET status = 'EXPIRED'
			WHERE id = $1
		`, req.ID)
		if err != nil {
			fmt.Printf("Failed to expire request ID %d: %v\n", req.ID, err)
			continue
		}

		// Step 2b: Get requester account by account_number (RequesterID)
		var requesterAccount models.Account
		err = dbConn.QueryRow(`
			SELECT id, user_id, account_number, balance
			FROM accounts 
			WHERE account_number = $1
		`, req.RequesterID).Scan(&requesterAccount.ID, &requesterAccount.UserID, &requesterAccount.AccountNumber, &requesterAccount.Balance)
		if err != nil {
			fmt.Printf("Requester account not found for request ID %d: %v\n", req.ID, err)
			continue
		}

		// Step 2c: Send WebSocket notification
		message := fmt.Sprintf("Your money request (Account %v) has expired", requesterAccount.AccountNumber)
		websocket.NotifyChan <- websocket.NotificationMessage{
			UserID:  requesterAccount.UserID,
			Message: message,
		}
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
	}

	fmt.Println("Expired request processing completed.")
}



func GetTransactionByID(id uint) (*models.Transaction, error) {
	// Prepare the query to fetch the transaction
	query := `SELECT id, account_id, to_account_id, amount, transaction_type, description, created_at 
			  FROM transactions 
			  WHERE id = $1 
			  LIMIT 1`

	// Create an instance to hold the result
	var tx models.Transaction

	// Execute the query
	err := db.DB.QueryRow(query, id).Scan(
		&tx.ID,
		&tx.AccountID,
		&tx.ToAccountID,
		&tx.Amount,
		&tx.TransactionType,
		&tx.Description,
		
	)

	if err != nil {
		// Check if no rows were found (no transaction)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction with id %d not found", id)
		}
		// Handle any other error
		log.Println("Error executing query:", err)
		return nil, err
	}

	return &tx, nil
}

func GetAllTransactions() ([]models.Transaction, error) {
	var txs []models.Transaction

	// Prepare the SQL query to fetch all transactions
	query := `SELECT id, account_id, to_account_id, amount, transaction_type, description, created_at, updated_at FROM transactions`

	// Execute the query
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %v", err)
	}
	defer rows.Close()

	// Loop through the rows and map them to the Transaction struct
	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(&tx.ID, &tx.AccountID, &tx.ToAccountID, &tx.Amount, &tx.TransactionType, &tx.Description); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %v", err)
		}
		txs = append(txs, tx)
	}

	// Check for any error encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %v", err)
	}

	return txs, nil
}
func DeleteTransaction(id uint) error {
    // Execute the DELETE query using raw SQL
    result, err := db.DB.Exec(`
        DELETE FROM transactions WHERE id = $1
    `, id)

    if err != nil {
        return fmt.Errorf("failed to delete transaction: %v", err)
    }

    // Check if any rows were affected
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get affected rows: %v", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("transaction with ID %d not found", id)
    }

    return nil
}


