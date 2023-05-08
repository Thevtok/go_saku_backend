package model

type BankAcc struct {
	AccountID         uint   `json:"account_id"`
	UserID            uint   `json:"user_id"`
	BankName          string `json:"bank_name"`
	AccountNumber     string `json:"account_number"`
	AccountHolderName string `json:"account_holder_name"`
}

type BankAccResponse struct {
	UserID            uint   `json:"user_id"`
	AccountID         uint   `json:"account_id"`
	BankName          string `json:"bank_name"`
	AccountNumber     string `json:"account_number"`
	AccountHolderName string `json:"account_holder_name"`
}

type CreateBankAcc struct {
	UserID            uint   `json:"user_id"`
	BankName          string `json:"bank_name"`
	AccountNumber     string `json:"account_number"`
	AccountHolderName string `json:"account_holder_name"`
}
