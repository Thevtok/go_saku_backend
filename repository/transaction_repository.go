package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ReygaFitra/inc-final-project.git/model"
)

var now = time.Now().UTC().Truncate(time.Minute)
var waktu = now.Format("2006-01-02 15:04")

type TransactionRepository interface {
	CreateDepositBank(tx *model.TransactionBank) error
	CreateDepositCard(tx *model.TransactionCard) error
	CreateWithdrawal(tx *model.TransactionBank) error
	CreateTransfer(tx *model.TransactionTransfer) (any, error)
	CreateRedeem(tx *model.TransactionPoint) error
}

type transactionRepository struct {
	db *sql.DB
}

func NewTxRepository(db *sql.DB) TransactionRepository {
	repo := new(transactionRepository)
	repo.db = db
	return repo
}

func (r *transactionRepository) CreateDepositBank(tx *model.TransactionBank) error {
	query := `INSERT INTO tx_transaction (transaction_type, sender_id, bank_account_id, amount, timestamp)
              VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(query, "Deposit Bank", tx.SenderID, tx.BankAccountID, tx.Amount, waktu)
	if err != nil {
		return err
	}

	return nil
}
func (r *transactionRepository) CreateDepositCard(tx *model.TransactionCard) error {
	query := `INSERT INTO tx_transaction (transaction_type, sender_id, card_id, amount, timestamp)
              VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(query, "Deposit Card", tx.SenderID, tx.CardID, tx.Amount, waktu)
	if err != nil {
		return err
	}

	return nil
}

func (r *transactionRepository) CreateWithdrawal(tx *model.TransactionBank) error {
	query := `INSERT INTO tx_transaction (transaction_type, sender_id, bank_account_id, amount, timestamp)
              VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(query, "Withdraw", tx.SenderID, tx.BankAccountID, tx.Amount, waktu)
	if err != nil {
		return err
	}

	return nil
}
func (r *transactionRepository) CreateTransfer(tx *model.TransactionTransfer) (any, error) {

	query := `INSERT INTO tx_transaction (transaction_type, sender_id, recipient_id, amount, timestamp)
              VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(query, "Transfer", tx.SenderID, tx.RecipientID, tx.Amount, waktu)
	if err != nil {
		return nil, fmt.Errorf("failed to create data: %v", err)
	}

	return tx, nil
}

func (r *transactionRepository) CreateRedeem(tx *model.TransactionPoint) error {
	query := `INSERT INTO tx_transaction (transaction_type, sender_id, point_exchange_id, point, timestamp)
              VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(query, tx.TransactionType, tx.SenderID, tx.PointExchangeID, tx.Point, waktu)
	if err != nil {
		return err
	}

	return nil
}
