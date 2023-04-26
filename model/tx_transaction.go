package model

import "time"

type Transaction struct {
	TransactionType string    `json:"transaction_type"`
	SenderID        uint      `json:"sender_id"`
	RecipientID     *uint     `json:"recipient_id"`
	BankAccountID   *uint     `json:"bank_account_id"`
	CardID          *uint     `json:"card_id"`
	PointExchangeID *uint     `json:"pe_id"`
	Amount          *uint     `json:"amount"`
	Point           *uint     `json:"point"`
	TransactionDate time.Time `json:"transaction_date"`
}

type TransactionBank struct {
	ID              uint      `json:"tx_id"`
	TransactionType string    `json:"transaction_type"`
	SenderID        uint      `json:"sender_id"`
	BankAccountID   uint      `json:"bank_account_id"`
	Amount          uint      `json:"amount"`
	TransactionDate time.Time `json:"transaction_date"`
}

type TransactionCard struct {
	ID              uint      `json:"tx_id"`
	TransactionType string    `json:"transaction_type"`
	SenderID        uint      `json:"sender_id"`
	CardID          uint      `json:"card_id"`
	Amount          uint      `json:"amount"`
	TransactionDate time.Time `json:"transaction_date"`
}

type TransactionWithdraw struct {
	ID              uint      `json:"tx_id"`
	TransactionType string    `json:"transaction_type"`
	BankAccountID   uint      `json:"bank_account_id"`
	SenderID        uint      `json:"sender_id"`
	Amount          uint      `json:"amount"`
	TransactionDate time.Time `json:"transaction_date"`
}

type TransactionTransfer struct {
	ID              uint      `json:"tx_id"`
	TransactionType string    `json:"transaction_type"`
	SenderID        uint      `json:"sender_id"`
	RecipientID     uint      `json:"recipient_id"`
	Amount          uint      `json:"amount"`
	TransactionDate time.Time `json:"transaction_date"`
}

type TransactionPoint struct {
	ID              uint      `json:"tx_id"`
	Reward          string    `json:"reward"`
	TransactionType string    `json:"transaction_type"`
	SenderID        uint      `json:"sender_id"`
	Point           int       `json:"point"`
	PointExchangeID int       `json:"pe_id"`
	TransactionDate time.Time `json:"transaction_date"`
}
