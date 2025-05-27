package services

import (
	"bank/db"
	"bank/dtos"
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
		return errors.New("receiver account not found")
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
		fmt.Println("Requester and recipient cannot be the same", request.RequesterID, request.RecipientID)
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

	// Get recipient account user ID
	var recipientUserID uint
	err = dbtx.QueryRow(`
		SELECT user_id FROM accounts WHERE account_number = $1
	`, request.RecipientID).Scan(&recipientUserID)
	if err != nil {
		dbtx.Rollback()
		return errors.New("recipient account not found")
	}

	// Insert money request
	var requestID uint
	err = dbtx.QueryRow(`
		INSERT INTO money_requests 
			(requester_id, recipient_id, amount, status, user_id,recipient_user_id, expires_at, requeste_at)
		VALUES 
			($1, $2, $3, $4, $5,$6, NOW(), NOW())
		RETURNING id
	`, request.RequesterID, request.RecipientID, request.Amount, request.Status, request.UserID, recipientUserID).Scan(&requestID)
	if err != nil {
		dbtx.Rollback()
		return err
	}

	// Log audit for the money request creation
	_ = LogAudit(nil, "CREATE", "money_requests", requestID, fmt.Sprintf("Money request of %.2f from %s to %s", request.Amount, request.RequesterID, request.RecipientID))

	// Prepare WebSocket notification message
	message := fmt.Sprintf("User %v requested %.2f from you", request.RequesterID, request.Amount)

	// Insert notification into the database
	_, err = dbtx.Exec(`
		INSERT INTO notifications (user_id, message, created_at)
		VALUES ($1, $2, NOW())
	`, recipientUserID, message)
	if err != nil {
		dbtx.Rollback()
		return errors.New("failed to insert notification")
	}

	// Commit the transaction
	if err := dbtx.Commit(); err != nil {
		return err
	}

	// Send WebSocket notification
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
		SELECT id, requester_id, recipient_id, amount, status ,user_id
		FROM money_requests 
		WHERE id = $1
	`, requestID).Scan(&req.ID, &req.RequesterID, &req.RecipientID, &req.Amount, &req.Status, &req.UserID)
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

	// Use existing transfer logic
	tx := &models.Transaction{
		UserID:      req.UserID,
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

	// Prepare notification message
	message := fmt.Sprintf("Your money request  %v was declined", req.RequesterID)

	// Insert notification into the database
	_, err = dbtx.Exec(`
		INSERT INTO notifications (user_id, message, created_at)
		VALUES ($1, $2, NOW())
	`, userID, message)
	if err != nil {
		dbtx.Rollback()
		return fmt.Errorf("failed to insert notification: %v", err)
	}

	// Commit the transaction
	if err := dbtx.Commit(); err != nil {
		return err
	}

	// Send WebSocket notification
	websocket.NotifyChan <- websocket.NotificationMessage{
		UserID:  userID,
		Message: message,
	}

	return nil
}

func AutoExpireRequests() {
	dbConn := db.DB

	fmt.Println("Checking for expired requests...", time.Now())

	// Step 1: Get all PENDING requests
	rows, err := dbConn.Query(`
		SELECT id, requester_id, recipient_id, amount, status
		FROM money_requests 
		WHERE status = 'PENDING' AND expires_at <= NOW()
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

		// Start a new transaction for each request
		dbtx, err := dbConn.Begin()
		if err != nil {
			log.Printf("Failed to begin transaction for request ID %d: %v\n", req.ID, err)
			continue
		}

		func() {
			defer func() {
				if p := recover(); p != nil {
					dbtx.Rollback()
					panic(p)
				}
			}()

			// Step 2a: Mark as EXPIRED
			_, err = dbtx.Exec(`
				UPDATE money_requests 
				SET status = 'EXPIRED'
				WHERE id = $1
			`, req.ID)
			if err != nil {
				dbtx.Rollback()
				log.Printf("Failed to expire request ID %d: %v\n", req.ID, err)
				return
			}

			// Step 2b: Get requester account
			var requesterAccount models.Account
			err = dbtx.QueryRow(`
				SELECT id, user_id, account_number, balance
				FROM accounts 
				WHERE account_number = $1
			`, req.RequesterID).Scan(&requesterAccount.ID, &requesterAccount.UserID, &requesterAccount.AccountNumber, &requesterAccount.Balance)
			if err != nil {
				dbtx.Rollback()
				log.Printf("Requester account not found for request ID %d: %v\n", req.ID, err)
				return
			}

			// Step 2c: Create notification
			message := fmt.Sprintf("Your money request (Account %v) has expired", requesterAccount.AccountNumber)
			_, err = dbtx.Exec(`
				INSERT INTO notifications (user_id, message, created_at)
				VALUES ($1, $2, NOW())
			`, requesterAccount.UserID, message)
			if err != nil {
				dbtx.Rollback()
				log.Printf("Failed to insert notification for request ID %d: %v\n", req.ID, err)
				return
			}

			// Commit transaction
			if err := dbtx.Commit(); err != nil {
				log.Printf("Failed to commit transaction for request ID %d: %v\n", req.ID, err)
				return
			}

			// Step 2d: Send WebSocket notification
			websocket.NotifyChan <- websocket.NotificationMessage{
				UserID:  requesterAccount.UserID,
				Message: message,
			}
		}()
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
	baseQuery := `SELECT id, user_id, account_id, transaction_type, to_account_id, amount, transaction_date, description 
	              FROM transactions WHERE 1=1`
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
	if filter.DescriptionLike != nil {
		likePattern := "%" + *filter.DescriptionLike + "%"
		params = append(params, likePattern)
		conditions += fmt.Sprintf(" AND description ILIKE $%d", len(params))
	}

	// Final SQL query with ORDER BY
	fullQuery := baseQuery + conditions + " ORDER BY transaction_date DESC"

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

	// Ensure we return empty slice, not nil
	if transactions == nil {
		transactions = []models.Transaction{}
	}

	return transactions, nil
}

func GetMoneyRequestsByUserID(userID uint) ([]models.MoneyRequest, error) {
	query := `
		SELECT id, requester_id, recipient_id, amount, status, user_id, expires_at, requeste_at
		FROM money_requests
		WHERE user_id = $1 or  recipient_user_id= $1 
		ORDER BY requeste_at DESC
	`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.MoneyRequest

	for rows.Next() {
		var req models.MoneyRequest
		err := rows.Scan(
			&req.ID,
			&req.RequesterID,
			&req.RecipientID,
			&req.Amount,
			&req.Status,
			&req.UserID,
			&req.ExpiresAt,
			&req.RequesteAt,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Ensure we return an empty slice, not nil
	if requests == nil {
		requests = []models.MoneyRequest{}
	}

	return requests, nil
}

func GetFilteredNotifications(userID string, filter string) ([]map[string]interface{}, error) {
	baseQuery := "SELECT  message, created_at FROM notifications WHERE user_id = $1"
	var query string

	switch filter {
	case "requests":
		query = baseQuery + " AND message ILIKE '%requested%' ORDER BY created_at DESC"
	case "alert":
		query = baseQuery + " AND (message ILIKE '%declined%' OR message ILIKE '%expired%') ORDER BY created_at DESC"
	default:
		query = baseQuery + " ORDER BY created_at DESC"
	}
	fmt.Println(query, userID)
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {

		var msg string
		var created string
		if err := rows.Scan(&msg, &created); err != nil {
			continue // optionally log this scan error
		}
		results = append(results, map[string]interface{}{

			"message":    msg,
			"created_at": created,
		})
	}

	return results, nil
}

func GetDashboardSummary(userID uint) (*dtos.DashboardSummary, error) {
	var summary dtos.DashboardSummary

	// 1. Wallet balance
	err := db.DB.QueryRow(`
		SELECT COALESCE(SUM(balance), 0)
		FROM accounts
		WHERE user_id = $1
	`, userID).Scan(&summary.WalletBalance)
	if err != nil {
		return nil, err
	}

	// 2. Total transactions
	err = db.DB.QueryRow(`
		SELECT COUNT(*)
		FROM transactions
		WHERE user_id = $1
	`, userID).Scan(&summary.TotalTransactions)
	if err != nil {
		return nil, err
	}

	// 3. Pending money requests
	err = db.DB.QueryRow(`
		SELECT COUNT(*)
		FROM money_requests
		WHERE user_id = $1 AND status = 'pending'
	`, userID).Scan(&summary.PendingRequests)
	if err != nil {
		return nil, err
	}

	// 4. Total transfers
	err = db.DB.QueryRow(`
		SELECT COUNT(*)
		FROM transactions
		WHERE user_id = $1 AND transaction_type = 'DEBIT'
	`, userID).Scan(&summary.TotalTransfers)
	if err != nil {
		return nil, err
	}

	// 5. Total sent amount
	err = db.DB.QueryRow(`
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE user_id = $1 AND transaction_type = 'DEBIT'
	`, userID).Scan(&summary.TotalSentAmount)
	if err != nil {
		return nil, err
	}

	// 6. Total received amount
	err = db.DB.QueryRow(`
		SELECT COALESCE(SUM(t.amount), 0)
		FROM transactions t
		JOIN accounts a ON a.account_number = t.account_id
		WHERE a.user_id = $1 AND t.transaction_type = 'CREDIT'
	`, userID).Scan(&summary.TotalReceivedAmount)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}

func GetMonthlyTransactionVolume(userID uint) ([]dtos.MonthlyTransactionVolume, error) {
	query := `
		WITH months AS (
			SELECT generate_series(1, 12) AS month_number
		),
		monthly_data AS (
			SELECT 
				EXTRACT(MONTH FROM transaction_date)::int AS month_number,
				SUM(amount) AS total_volume
			FROM transactions
			WHERE user_id = $1
			GROUP BY month_number
		)
		SELECT 
			TO_CHAR(TO_DATE(m.month_number::text, 'MM'), 'Mon') AS month,
			COALESCE(md.total_volume, 0) AS total_volume
		FROM months m
		LEFT JOIN monthly_data md ON m.month_number = md.month_number
		ORDER BY m.month_number;
	`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dtos.MonthlyTransactionVolume
	for rows.Next() {
		var entry dtos.MonthlyTransactionVolume
		if err := rows.Scan(&entry.Name, &entry.Total); err != nil {
			return nil, err
		}
		results = append(results, entry)
	}

	return results, nil
}
