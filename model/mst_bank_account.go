package model

type BankAcc struct {
	AccountId         uint    `json:"account_id"`
	UserID            uint    `json:"user_id"`
	BankName          string  `json:"bank_name"`
	AccountNumber     string  `json:"account_number"`
	AccountHolderName string  `json:"account_holder_name"`
	Email             string  `json:"email"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         *string `json:"updated_at,omitempty"`
	Name              string  `json:"name"`
}

type BankAccResponse struct {
	UserID uint `json:"user_id"`

	BankName          string `json:"bank_name"`
	AccountNumber     string `json:"account_number"`
	AccountHolderName string `json:"account_holder_name"`
}
