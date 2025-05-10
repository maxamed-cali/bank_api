package services

import (
	"bank/db"
	"bank/models"
	"bank/websocket"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)



func MoneyTransfer(tx *models.Transaction) error {
	database := db.DB

	if tx.ToAccountID == nil {
		return errors.New("missing destination account for transfer")
	}

	var sender models.Account
	if err := database.First(&sender, "account_number = ?", tx.AccountID).Error; err != nil {
	
    return fmt.Errorf("sender account not found: %v", err)
}


	var receiver models.Account
	if err := database.First(&receiver, "account_number = ?", *tx.ToAccountID).Error; err != nil {
			
		return fmt.Errorf("receiver account not found: %v", err)
	}

	if sender.AccountNumber == receiver.AccountNumber {
		return errors.New("cannot transfer to self")
	}

	if !sender.IsActive || !receiver.IsActive {
		return errors.New("both accounts must be active")
	}

	if sender.Balance < tx.Amount {
		return errors.New("insufficient balance")
	}

	return database.Transaction(func(txDB *gorm.DB) error {
		// Update balances
		sender.Balance -= tx.Amount
		receiver.Balance += tx.Amount

		if err := txDB.Save(&sender).Error; err != nil {
			return err
		}
		if err := txDB.Save(&receiver).Error; err != nil {
			return err
		}

		// Record sender transaction (DEBIT)
		senderTx := models.Transaction{
			AccountID:       sender.AccountNumber,
			ToAccountID:     &receiver.AccountNumber,
			TransactionType: "DEBIT",
			Amount:          tx.Amount,
			Description:     fmt.Sprintf("Transferred to Account ID %s", receiver.AccountNumber),
		}
		if err := txDB.Create(&senderTx).Error; err != nil {
			return err
		}

		// Record receiver transaction (CREDIT)
		receiverTx := models.Transaction{
			AccountID:       receiver.AccountNumber,
			ToAccountID:     &sender.AccountNumber,
			TransactionType: "CREDIT",
			Amount:          tx.Amount,
			Description:     fmt.Sprintf("Received from Account ID %s", sender.AccountNumber),
		}
		if err := txDB.Create(&receiverTx).Error; err != nil {
			return err
		}

		// Save notification to DB
		message := fmt.Sprintf("You received %.2f from %s", tx.Amount, sender.AccountNumber)
		notification := models.Notification{
			UserID:  receiver.UserID,
			Message: message,
		}
		if err := txDB.Create(&notification).Error; err != nil {
			return err
		}

		// Send notification to WebSocket client (real-time)
		websocket.NotifyChan <- websocket.NotificationMessage{
			UserID:  receiver.UserID,
			Message: message,
		}

		return nil
	})
}

func MoneyRequest(request *models.MoneyRequest) error {
	database := db.DB

	// Optional validation checks
	if request.Amount <= 0 {
		return errors.New("invalid amount")
	}
    
 	if request.RequesterID == request.RecipientID {
		fmt.Println("Requester and recipient cannot be the same", request.RequesterID, request.RecipientID)
		return errors.New("cannot request from self")
	}

	request.Status = "PENDING"
	request.ExpiresAt = time.Now().Add(24 * time.Hour) // auto-expiry in 24h

	if err := database.Create(&request).Error; err != nil {
		return err
	}

	//  Get UserID of RecipientID (assuming RecipientID is an account number)
	var recipientAccount models.Account
	if err := database.First(&recipientAccount, "account_number = ?", request.RecipientID).Error; err != nil {
		return fmt.Errorf("recipient account not found: %v", err)
	}

	// Notify recipient via WebSocket
	message := fmt.Sprintf("User %v requested %.2f from you", request.RequesterID, request.Amount)
	websocket.NotifyChan <- websocket.NotificationMessage{
		UserID:  recipientAccount.UserID,
		Message: message,
	}

	return nil
}

func AcceptMoneyRequest(requestID uint) error {
	database := db.DB

	var req models.MoneyRequest
	if err := database.First(&req, requestID).Error; err != nil {
		return err
	}

	if req.Status != "PENDING" {
		return errors.New("request is no longer active")
	}

	// Mark accepted before transfer (to avoid double execution)
	req.Status = "ACCEPTED"
	if err := database.Save(&req).Error; err != nil {
		return err
	}

	// Perform fund transfer
	tx := &models.Transaction{
		AccountID:    req.RecipientID,
		ToAccountID:  &req.RequesterID,
		Amount:       req.Amount,
		Description:  fmt.Sprintf("Accepted request ID %d", req.ID),
	}
	return MoneyTransfer(tx)
}
