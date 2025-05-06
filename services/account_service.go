package services

import (
	"bank/db"
	"bank/models"
	"fmt"
)

// Create a wallet
func CreateWallet(account *models.Account) error {
	fmt.Print(account,"account")

	return db.DB.Create(account).Error
}

// Delete a wallet
func DeleteWallet(accountID uint) error {
	return db.DB.Delete(&models.Account{}, accountID).Error
}

// Rename wallet (change account number)
func RenameWallet(accountID uint, newNumber string) error {
	return db.DB.Model(&models.Account{}).Where("id = ?", accountID).Update("account_number", newNumber).Error
}



func GetBalanceByAccountNumber(accountNumber string) (float64, error) {
	var account models.Account
	err := db.DB.Where("account_number = ? AND is_active = ?", accountNumber, true).First(&account).Error
	if err != nil {
		return 0, err
	}
	return account.Balance, nil
}


// View balances grouped by currency
type WalletInfo struct {
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
}

type CurrencyWallets struct {
	Currency     string       `json:"currency"`
	TotalBalance float64      `json:"total_balance"`
	Wallets      []WalletInfo `json:"wallets"`
}

func GetBalancesGroupedByCurrency() ([]CurrencyWallets, error) {
	var accounts []models.Account
	err := db.DB.Preload("AccountType").Where("is_active = ?", true).Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	grouped := make(map[string]*CurrencyWallets)
	for _, acc := range accounts {
		currency := acc.AccountType.Currency
		if _, exists := grouped[currency]; !exists {
			grouped[currency] = &CurrencyWallets{
				Currency:     currency,
				Wallets:      []WalletInfo{},
				TotalBalance: 0,
			}
		}
		grouped[currency].Wallets = append(grouped[currency].Wallets, WalletInfo{
			AccountNumber: acc.AccountNumber,
			Balance:       acc.Balance,
		})
		grouped[currency].TotalBalance += acc.Balance
	}

	// Convert map to slice
	var result []CurrencyWallets
	for _, g := range grouped {
		result = append(result, *g)
	}
	return result, nil
}
