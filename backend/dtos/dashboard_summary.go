package dtos

type DashboardSummary struct {
	WalletBalance       float64 `json:"wallet_balance"`
	TotalTransactions   int     `json:"total_transactions"`
	PendingRequests     int     `json:"pending_requests"`
	TotalTransfers      int     `json:"total_transfers"`
	TotalSentAmount     float64 `json:"total_sent_amount"`
	TotalReceivedAmount float64 `json:"total_received_amount"`
}


type MonthlyTransactionVolume struct {
	Name        string  `json:"name"`
	Total  float64 `json:"total"`
}

