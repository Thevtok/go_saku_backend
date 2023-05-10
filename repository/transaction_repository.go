package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ReygaFitra/inc-final-project.git/model"
)

var now = time.Now().UTC().Truncate(time.Minute)

type TransactionRepository interface {
	CreateDepositBank(tx *model.TransactionBank) error
	CreateDepositCard(tx *model.TransactionCard) error
	CreateWithdrawal(tx *model.TransactionWithdraw) error
	CreateTransfer(tx *model.TransactionTransferResponse) (any, error)
	CreateRedeem(tx *model.TransactionPoint) error
	GetAllPoint() ([]*model.PointExchange, error)
	GetBySenderId(senderId, recipientId uint) ([]*model.Transaction, error)
	GetByPeId(id int) (*model.PointExchange, error)
}

type transactionRepository struct {
	db *sql.DB
}

func (r *transactionRepository) GetBySenderId(senderId, recipientId uint) ([]*model.Transaction, error) {
	var txs []*model.Transaction
	rows, err := r.db.Query(`
        SELECT  transaction_type, sender_id, recipient_id, bank_account_id, card_id, pe_id, amount, point, transaction_date
        FROM tx_transaction
        WHERE sender_id = $1 OR recipient_id = $2
    `, senderId, recipientId)
	if err != nil {
		return nil, fmt.Errorf("error while getting transactions for sender %v: %v", senderId, err)
	}
	defer rows.Close()

	for rows.Next() {

		tx := &model.Transaction{}
		err := rows.Scan(&tx.TransactionType, &tx.SenderID, &tx.RecipientID, &tx.BankAccountID, &tx.CardID, &tx.PointExchangeID, &tx.Amount, &tx.Point, &tx.TransactionDate)
		if err != nil {
			return nil, fmt.Errorf("error while scanning transaction: %v", err)
		}
		if tx.RecipientID == nil {
			tx.RecipientID = new(uint)
		}
		if tx.BankAccountID == nil {
			tx.BankAccountID = new(uint)
		}
		if tx.CardID == nil {
			tx.CardID = new(uint)
		}
		if tx.PointExchangeID == nil {
			tx.PointExchangeID = new(uint)
		}
		if tx.Amount == nil {
			tx.Amount = new(uint)
		}
		if tx.Point == nil {
			tx.Point = new(uint)
		}
		txs = append(txs, tx)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while getting transactions for sender %v: %v", senderId, err)
	}

	return txs, nil
}

func (r *transactionRepository) CreateDepositBank(tx *model.TransactionBank) error {
	query := "INSERT INTO tx_transaction (transaction_type, sender_id, bank_account_id, amount, transaction_date) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(query, "Deposit Bank", tx.SenderID, tx.BankAccountID, tx.Amount, now)
	if err != nil {
		return err
	}

	return nil
}

func (r *transactionRepository) CreateDepositCard(tx *model.TransactionCard) error {
	query := "INSERT INTO tx_transaction (transaction_type, sender_id, card_id, amount, transaction_date) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(query, "Deposit Card", tx.SenderID, tx.CardID, tx.Amount, now)
	if err != nil {
		return err
	}

	return nil
}

func (r *transactionRepository) CreateWithdrawal(tx *model.TransactionWithdraw) error {
	query := "INSERT INTO tx_transaction (transaction_type, bank_account_id, sender_id, amount, transaction_date) VALUES ($1, $2, $3, $4,$5)"
	_, err := r.db.Exec(query, "Withdraw", tx.BankAccountID, tx.SenderID, tx.Amount, now)
	if err != nil {
		return err
	}

	return nil
}

func (r *transactionRepository) CreateTransfer(tx *model.TransactionTransferResponse) (any, error) {
	query := "INSERT INTO tx_transaction (transaction_type, sender_id, recipient_id, amount, transaction_date) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(query, "Transfer", tx.SenderID, tx.RecipientID, tx.Amount, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create data: %v", err)
	}

	return tx, nil
}

func (r *transactionRepository) CreateRedeem(tx *model.TransactionPoint) error {
	query := "INSERT INTO tx_transaction (transaction_type, sender_id, pe_id, point, transaction_date) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(query, "Redeem", tx.SenderID, tx.PointExchangeID, tx.Point, now)
	if err != nil {
		return err
	}

	return nil
}

// Get all point exchanges
func (r *transactionRepository) GetAllPoint() ([]*model.PointExchange, error) {
	query := "SELECT pe_id, reward, price FROM mst_point_exchange"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %v", err)
	}
	defer rows.Close()

	var pointExchanges []*model.PointExchange
	for rows.Next() {
		pe := &model.PointExchange{}
		err := rows.Scan(&pe.PE_ID, &pe.Reward, &pe.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to scan data: %v", err)
		}
		pointExchanges = append(pointExchanges, pe)
	}

	return pointExchanges, nil
}

func (r *transactionRepository) GetByPeId(id int) (*model.PointExchange, error) {
	var peAcc model.PointExchange
	query := "SELECT pe_id, reward, price FROM mst_point_exchange WHERE pe_id = $1"
	err := r.db.QueryRow(query, id).Scan(&peAcc.PE_ID, &peAcc.Reward, &peAcc.Price)
	if err != nil {
		return nil, err
	}
	return &peAcc, nil
}

func NewTxRepository(db *sql.DB) TransactionRepository {
	repo := new(transactionRepository)
	repo.db = db
	return repo
}
