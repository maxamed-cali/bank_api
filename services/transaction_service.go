package services

import (
	"bank/db"
	"bank/models"
	"bank/websocket"
	"errors"
	"fmt"

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

