package model

type BankAcc struct {
	AccountID         uint    `json:"account_id"`
	Username          string  `json:"username"`
	BankName          string  `json:"bank_name"`
	AccountNumber     string  `json:"account_number"`
	AccountHolderName string  `json:"account_holder_name"`
	Email             string  `json:"email"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         *string `json:"updated_at,omitempty"`
	Name              string  `json:"name"`
}

type BankAccResponse struct {
	Username          string `json:"username"`
	BankName          string `json:"bank_name"`
	AccountNumber     string `json:"account_number"`
	AccountHolderName string `json:"account_holder_name"`
}
