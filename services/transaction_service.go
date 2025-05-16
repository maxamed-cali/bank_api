package services

import (
	"bank/db"
	"bank/models"
	"bank/websocket"
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

	// Update sender balance
	_, err = dbtx.Exec(`UPDATE accounts SET balance = balance - $1 WHERE account_number = $2`, tx.Amount, sender.AccountNumber)
	if err != nil {
		dbtx.Rollback()
		return err
	}

	// Update receiver balance
	_, err = dbtx.Exec(`UPDATE accounts SET balance = balance + $1 WHERE account_number = $2`, tx.Amount, receiver.AccountNumber)
	if err != nil {
		dbtx.Rollback()
		return err
	}

	// Insert sender transaction (DEBIT)
	_, err = dbtx.Exec(`INSERT INTO transactions (account_id, to_account_id, transaction_type, amount, description, user_id, transaction_date)
	                    VALUES ($1, $2, 'DEBIT', $3, $4, $5, NOW())`,
		sender.AccountNumber, receiver.AccountNumber, tx.Amount,
		fmt.Sprintf("Transferred to Account ID %s", receiver.AccountNumber), tx.UserID)
	if err != nil {
		dbtx.Rollback()
		return err
	}

	// Log audit for sender debit action
	_ = LogAudit(&sender.UserID, "CREATE", "transactions", sender.ID, fmt.Sprintf("Debited %.2f to %s", tx.Amount, receiver.AccountNumber))

	// Insert receiver transaction (CREDIT)
	_, err = dbtx.Exec(`INSERT INTO transactions (account_id, to_account_id, transaction_type, amount, description, user_id, transaction_date)
	                    VALUES ($1, $2, 'CREDIT', $3, $4, $5, NOW())`,
		receiver.AccountNumber, sender.AccountNumber, tx.Amount,
		fmt.Sprintf("Received from Account ID %s", sender.AccountNumber), tx.UserID)
	if err != nil {
		dbtx.Rollback()
		return err
	}

	// Log audit for receiver credit action
	_ = LogAudit(&receiver.UserID, "CREATE", "transactions", receiver.ID, fmt.Sprintf("Credited %.2f from %s", tx.Amount, sender.AccountNumber))

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
	var requestID uint
	err = dbtx.QueryRow(`
		INSERT INTO money_requests 
			(requester_id, recipient_id, amount, status,user_id, expires_at, requeste_at)
		VALUES 
			($1, $2, $3, $4,$5,  NOW(), NOW())
		RETURNING id
	`, request.RequesterID, request.RecipientID, request.Amount, request.Status,request.UserID).Scan(&requestID)
	if err != nil {
		dbtx.Rollback()
		return err
	}

	// Get recipient account user ID
	var recipientUserID uint
	err = dbtx.QueryRow(`
		SELECT user_id FROM accounts WHERE account_number = $1
	`, request.RecipientID).Scan(&recipientUserID)
	if err != nil {
		dbtx.Rollback()
		return fmt.Errorf("recipient account not found: %v", err)
	}

	// Log audit for the money request creation
	_ = LogAudit(nil, "CREATE", "money_requests", requestID, fmt.Sprintf("Money request of %.2f from %s to %s", request.Amount, request.RequesterID, request.RecipientID))

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



type TransactionFilter struct {
	UserID          *uint
	AccountID       *string
	TransactionType *string
	MinAmount       *float64
	MaxAmount       *float64
	StartDate       *time.Time
	EndDate         *time.Time
	DescriptionLike *string
}

func GetTransactionHistory(filter TransactionFilter) ([]models.Transaction, error) {
	baseQuery := `SELECT id, user_id, account_id, transaction_type, to_account_id, amount, transaction_date, description FROM transactions WHERE 1=1`
	var params []interface{}
	var conditions string

	if filter.UserID != nil {
		params = append(params, *filter.UserID)
		conditions += fmt.Sprintf(" AND user_id = $%d", len(params))
	}
	if filter.AccountID != nil {
		params = append(params, *filter.AccountID)
		conditions += fmt.Sprintf(" AND account_id = $%d", len(params))
	}
	if filter.TransactionType != nil {
		params = append(params, *filter.TransactionType)
		conditions += fmt.Sprintf(" AND transaction_type = $%d", len(params))
	}
	if filter.MinAmount != nil {
		params = append(params, *filter.MinAmount)
		conditions += fmt.Sprintf(" AND amount >= $%d", len(params))
	}
	if filter.MaxAmount != nil {
		params = append(params, *filter.MaxAmount)
		conditions += fmt.Sprintf(" AND amount <= $%d", len(params))
	}
	if filter.StartDate != nil {
		params = append(params, *filter.StartDate)
		conditions += fmt.Sprintf(" AND transaction_date >= $%d", len(params))
	}
	if filter.EndDate != nil {
		params = append(params, *filter.EndDate)
		conditions += fmt.Sprintf(" AND transaction_date <= $%d", len(params))
	}
	if filter.DescriptionLike != nil {likePattern := "%" + *filter.DescriptionLike + "%"
		params = append(params, likePattern)
		conditions += fmt.Sprintf(" AND description ILIKE $%d", len(params))
	}

	fullQuery := baseQuery + conditions

	rows, err := db.GetDB().Query(fullQuery, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(&t.ID, &t.UserID, &t.AccountID, &t.TransactionType, &t.ToAccountID, &t.Amount, &t.TransactionDate, &t.Description)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}


