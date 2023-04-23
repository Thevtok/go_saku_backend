package model

import "time"

type Transaction struct {
	ID              uint      `json:"id"`
	TransactionType string    `json:"transaction_type"`
	SenderID        uint      `json:"sender_id"`
	RecipientID     uint      `json:"recipient_id"`
	BankAccountID   uint      `json:"bank_account_id"`
	CardID          uint      `json:"card_id"`
	PointExchangeID uint      `json:"point_exchange_id"`
	Amount          uint      `json:"amount"`
	Point           uint      `json:"point"`
	Timestamp       time.Time `json:"timestamp"`
}
type TransactionBank struct {
	ID              uint   `json:"tx_id"`
	TransactionType string `json:"transaction_type"`
	SenderID        uint   `json:"sender_id"`

	BankAccountID uint `json:"bank_account_id"`

	Amount uint `json:"amount"`

	Timestamp time.Time `json:"timestamp"`
}
type TransactionCard struct {
	ID              uint   `json:"tx_id"`
	TransactionType string `json:"transaction_type"`
	SenderID        uint   `json:"sender_id"`

	CardID uint `json:"card_id"`

	Amount uint `json:"amount"`

	Timestamp time.Time `json:"timestamp"`
}

type TransactionTransfer struct {
	ID              uint   `json:"tx_id"`
	TransactionType string `json:"transaction_type"`
	SenderID        uint   `json:"sender_id"`
	RecipientID     uint   `json:"recipient_id"`

	Amount uint `json:"amount"`

	Timestamp time.Time `json:"timestamp"`
}

type TransactionPoint struct {
	ID              uint   `json:"tx_id"`
	TransactionType string `json:"transaction_type"`
	SenderID        uint   `json:"sender_id"`
	Point           uint   `json:"point"`
	PointExchangeID uint   `json:"point_exchange_id"`

	Timestamp time.Time `json:"timestamp"`
}
