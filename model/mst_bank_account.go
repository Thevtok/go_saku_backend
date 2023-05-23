package model

type BankAcc struct {
	AccountID         uint   `json:"account_id"`
	UserID            string `json:"user_id"`
	BankName          string `json:"bank_name"`
	AccountNumber     string `json:"account_number"`
	AccountHolderName string `json:"account_holder_name"`
}
