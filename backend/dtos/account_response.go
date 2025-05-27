package dtos

type AccountResponse struct {
	ID             uint    `json:"id"`
	AccountNumber  string  `json:"account_number"`
	Balance        float64 `json:"balance"`
	UserID         uint    `json:"user_id"`
	AccountTypeID  uint    `json:"account_type_id"`
	TypeName       string  `json:"type_name"`
	Description    string  `json:"description"`
	Currency       string  `json:"currency"`
	Name       string  `json:"name"`
	CreatedAt      string  `json:"created_at"`
}
