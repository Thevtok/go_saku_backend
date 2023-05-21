package model

type Transaction struct {
	TransactionType   string  `json:"transaction_type"`
	SenderID          *uint   `json:"sender_id"`
	RecipientID       *uint   `json:"recipient_id"`
	BankAccountID     *uint   `json:"bank_account_id"`
	CardID            *uint   `json:"card_id"`
	PointExchangeID   *uint   `json:"pe_id"`
	Amount            *uint   `json:"amount"`
	Point             *uint   `json:"point"`
	TransactionDate   string  `json:"transaction_date"`
	SenderNumber      *string `json:"sender_phone_number"`
	RecipientNumber   *string `json:"recipient_phone_number"`
	SenderName        *string `json:"sender_name"`
	RecipientName     *string `json:"recipient_name"`
	BankName          *string `json:"bank_name"`
	BankAccountNumber *string `json:"bank_account_number"`
}

type TransactionBank struct {
	ID                uint   `json:"tx_id"`
	TransactionType   string `json:"transaction_type"`
	SenderID          uint   `json:"sender_id"`
	BankAccountID     uint   `json:"bank_account_id"`
	Amount            uint   `json:"amount"`
	TransactionDate   string `json:"transaction_date"`
	SenderName        string `json:"sender_name"`
	BankName          string `json:"bank_name"`
	BankAccountNumber string `json:"bank_account_number"`
}

type TransactionCard struct {
	ID              uint   `json:"tx_id"`
	TransactionType string `json:"transaction_type"`
	SenderID        uint   `json:"sender_id"`
	CardID          uint   `json:"card_id"`
	Amount          uint   `json:"amount"`
	TransactionDate string `json:"transaction_date"`
}

type TransactionWithdraw struct {
	ID                uint   `json:"tx_id"`
	TransactionType   string `json:"transaction_type"`
	BankAccountID     uint   `json:"bank_account_id"`
	SenderID          uint   `json:"sender_id"`
	Amount            uint   `json:"amount"`
	TransactionDate   string `json:"transaction_date"`
	SenderName        string `json:"sender_name"`
	BankName          string `json:"bank_name"`
	BankAccountNumber string `json:"bank_account_number"`
}

type TransactionTransfer struct {
	ID              uint   `json:"tx_id"`
	TransactionType string `json:"transaction_type"`
	SenderID        uint   `json:"sender_id"`
	RecipientID     uint   `json:"recipient_id"`
	Amount          uint   `json:"amount"`
	TransactionDate string `json:"transaction_date"`
}
type TransactionTransferResponse struct {
	TransactionType string `json:"transaction_type"`
	SenderID        uint   `json:"sender_id"`
	RecipientID     uint   `json:"recipient_id"`
	Amount          uint   `json:"amount"`
	TransactionDate string `json:"transaction_date"`
	SenderNumber    string `json:"sender_phone_number"`
	RecipientNumber string `json:"recipient_phone_number"`
	SenderName      string `json:"sender_name"`
	RecipientName   string `json:"recipient_name"`
}

type TransactionPoint struct {
	ID uint `json:"tx_id"`

	TransactionType string `json:"transaction_type"`
	SenderID        uint   `json:"sender_id"`
	Point           int    `json:"point"`
	PointExchangeID int    `json:"pe_id"`
	TransactionDate string `json:"transaction_date"`
}
