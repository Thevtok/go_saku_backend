package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ReygaFitra/inc-final-project.git/model"
)

var now = time.Now().Local()
var date = now.Format("2006-01-02")

type TransactionRepository interface {
	CreateDepositBank(tx *model.Deposit) error

	CreateWithdrawal(tx *model.Withdraw) error
	CreateTransfer(tx *model.Transfer) error
	CreateRedeem(tx *model.Redeem) error
	GetAllPoint() ([]*model.PointExchange, error)
	GetTransactions(userID string) ([]*model.Transaction, error)
	GetByPeId(id int) (*model.PointExchange, error)
	AssignBadge(user *model.User) error
}

type transactionRepository struct {
	db *sql.DB
}

func (ur *transactionRepository) AssignBadge(user *model.User) error {
	queryCount := "SELECT COUNT(*) FROM tx_transfer WHERE sender_id = $1"
	row := ur.db.QueryRow(queryCount, user.ID)

	var txCount int
	err := row.Scan(&txCount)
	if err != nil {
		return fmt.Errorf("failed to retrieve transaction count: %v", err)
	}

	// Mengambil badge berdasarkan jumlah transaksi
	queryBadge := "SELECT badge_id FROM mst_badges WHERE threshold <= $1 ORDER BY threshold DESC LIMIT 1"
	row = ur.db.QueryRow(queryBadge, txCount)

	var badgeID int
	err = row.Scan(&badgeID)
	if err != nil {
		return fmt.Errorf("failed to retrieve badge ID: %v", err)
	}

	// Update field tx_count dan badge_id pada tabel pengguna
	queryUpdate := "UPDATE mst_users SET tx_count = $1, badge_id = $2 WHERE user_id = $3"
	_, err = ur.db.Exec(queryUpdate, txCount, badgeID, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	user.TxCount = txCount
	user.BadgeID = badgeID

	return nil
}

func (r *transactionRepository) GetTransactions(userID string) ([]*model.Transaction, error) {
	query := `
	SELECT 
    t.tx_id, t.transaction_type, t.transaction_date,
    d.bank_name, d.account_number, d.account_holder_name, d.amount,
    w.bank_name, w.account_number, w.account_holder_name, w.amount,
    tr.sender_name, tr.sender_phone_number, tr.recipient_name, tr.recipient_phone_number, tr.amount,
    CAST(rp.pe_id AS VARCHAR), rp.amount,
    pe.reward
FROM tx_transaction t
LEFT JOIN tx_deposit d ON t.tx_id = d.transaction_id
LEFT JOIN tx_withdraw w ON t.tx_id = w.transaction_id
LEFT JOIN tx_transfer tr ON t.tx_id = tr.transaction_id
LEFT JOIN tx_redeem rp ON t.tx_id = rp.transaction_id
LEFT JOIN mst_point_exchange pe ON rp.pe_id = pe.pe_id
WHERE (t.sender_id = $1 OR t.recipient_id = $1)
   
ORDER BY t.tx_id DESC



	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	transactions := []*model.Transaction{}
	for rows.Next() {
		var (
			txID                       int
			transactionType            string
			transactionDate            string
			depositBankName            sql.NullString
			deposit_bank_number        sql.NullString
			deposit_account_bank_name  sql.NullString
			deposit_amount             sql.NullInt64
			withdrawBankName           sql.NullString
			withdraw_bank_number       sql.NullString
			withdraw_account_bank_name sql.NullString
			withdraw_amount            sql.NullInt64
			transfer_sender_name       sql.NullString
			transfer_sender_phone      sql.NullString
			transfer_recipient_name    sql.NullString
			transfer_recipient_phone   sql.NullString
			transfer_amount            sql.NullInt64
			redeemPEID                 sql.NullString
			redeemAmount               sql.NullInt64
			redeemReward               sql.NullString
		)

		err := rows.Scan(&txID, &transactionType, &transactionDate, &depositBankName, &deposit_bank_number, &deposit_account_bank_name, &deposit_amount, &withdrawBankName, &withdraw_bank_number, &withdraw_account_bank_name, &withdraw_amount, &transfer_sender_name, &transfer_sender_phone, &transfer_recipient_name, &transfer_recipient_phone, &transfer_amount, &redeemPEID, &redeemAmount, &redeemReward)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %v", err)
		}

		transaction := &model.Transaction{
			TxID:            txID,
			TransactionType: transactionType,
			TransactionDate: transactionDate,
		}

		if depositBankName.Valid {
			transaction.DepositBankName = depositBankName.String
		}
		if deposit_bank_number.Valid {
			transaction.DepositBankNumber = deposit_bank_number.String
		}
		if deposit_account_bank_name.Valid {
			transaction.DepositAccountBankName = deposit_account_bank_name.String
		}
		if deposit_amount.Valid {
			transaction.DepositAmount = int(deposit_amount.Int64)
		}

		if withdrawBankName.Valid {
			transaction.WithdrawBankName = withdrawBankName.String
		}
		if withdraw_bank_number.Valid {
			transaction.WithdrawBankNumber = withdraw_bank_number.String
		}
		if withdraw_account_bank_name.Valid {
			transaction.WithdrawAccountBankName = withdraw_account_bank_name.String
		}
		if withdraw_amount.Valid {
			transaction.WithdrawAmount = int(withdraw_amount.Int64)
		}
		if transfer_sender_name.Valid {
			transaction.TransferSenderName = transfer_sender_name.String
		}
		if transfer_sender_phone.Valid {
			transaction.TransferSenderPhone = transfer_sender_phone.String
		}
		if transfer_recipient_name.Valid {
			transaction.TransferRecipientName = transfer_recipient_name.String
		}
		if transfer_recipient_phone.Valid {
			transaction.TransferRecipientPhone = transfer_recipient_phone.String
		}
		if transfer_amount.Valid {
			transaction.TransferAmount = int(transfer_amount.Int64)
		}

		if redeemPEID.Valid {
			transaction.RedeemPEID = redeemPEID.String
		}
		if redeemAmount.Valid {
			transaction.RedeemAmount = int(redeemAmount.Int64)
		}
		if redeemReward.Valid {
			transaction.RedeemReward = redeemReward.String
		}

		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate through result set: %v", err)
	}

	return transactions, nil
}

func (r *transactionRepository) CreateDepositBank(tx *model.Deposit) error {
	query := "INSERT INTO tx_transaction (transaction_type, transaction_date, sender_id) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(query, "Deposit", date, tx.UserID)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %v", err)
	}

	var txID int
	err = r.db.QueryRow("SELECT lastval()").Scan(&txID)
	if err != nil {
		return fmt.Errorf("failed to retrieve transaction ID: %v", err)
	}

	query = "INSERT INTO tx_deposit (transaction_id, amount, bank_name, account_number, account_holder_name) VALUES ($1, $2, $3, $4, $5)"
	_, err = r.db.Exec(query, txID, tx.Amount, tx.BankName, tx.AccountNumber, tx.AccountHolderName)
	if err != nil {
		return fmt.Errorf("failed to insert deposit: %v", err)
	}

	return nil
}
func (r *transactionRepository) CreateWithdrawal(tx *model.Withdraw) error {
	query := "INSERT INTO tx_transaction (transaction_type, transaction_date, sender_id) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(query, "Withdraw", date, tx.UserID)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %v", err)
	}

	var txID int
	err = r.db.QueryRow("SELECT lastval()").Scan(&txID)
	if err != nil {
		return fmt.Errorf("failed to retrieve transaction ID: %v", err)
	}

	query = "INSERT INTO tx_withdraw (transaction_id, amount, bank_name, account_number, account_holder_name) VALUES ($1, $2, $3, $4, $5)"
	_, err = r.db.Exec(query, txID, tx.Amount, tx.BankName, tx.AccountNumber, tx.AccountHolderName)
	if err != nil {
		return fmt.Errorf("failed to insert withdrawal: %v", err)
	}

	return nil
}

func (r *transactionRepository) CreateTransfer(tx *model.Transfer) error {
	query := "INSERT INTO tx_transaction (transaction_type, transaction_date, sender_id,recipient_id) VALUES ($1, $2, $3,$4)"
	_, err := r.db.Exec(query, "Transfer", date, tx.SenderID, tx.RecipientID)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %v", err)
	}

	var txID int
	err = r.db.QueryRow("SELECT lastval()").Scan(&txID)
	if err != nil {
		return fmt.Errorf("failed to retrieve transaction ID: %v", err)
	}

	query = "INSERT INTO tx_transfer (transaction_id, sender_name, recipient_name, amount, sender_phone_number, recipient_phone_number,sender_id,recipient_id) VALUES ($1, $2, $3, $4, $5, $6,$7,$8)"
	_, err = r.db.Exec(query, txID, tx.SenderName, tx.RecipientName, tx.Amount, tx.SenderPhoneNumber, tx.RecipientPhoneNumber, tx.SenderID, tx.RecipientID)
	if err != nil {
		return fmt.Errorf("failed to insert transfer: %v", err)
	}

	return nil
}

func (r *transactionRepository) CreateRedeem(tx *model.Redeem) error {
	query := "INSERT INTO tx_transaction (transaction_type, transaction_date,sender_id) VALUES ($1, $2,$3)"
	_, err := r.db.Exec(query, "Redeem", date, tx.UserID)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %v", err)
	}
	var txID int
	err = r.db.QueryRow("SELECT lastval()").Scan(&txID)
	if err != nil {
		return fmt.Errorf("failed to retrieve transaction ID: %v", err)
	}
	query = "INSERT INTO tx_redeem (transaction_id,amount, pe_id) VALUES ($1, $2, $3)"
	_, err = r.db.Exec(query, txID, tx.Amount, tx.PEID)
	if err != nil {
		return fmt.Errorf("failed to insert redeem: %v", err)
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
			return nil, fmt.Errorf("failed to scan store: %v", err)
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
		return nil, fmt.Errorf("failed to scan pe_id: %v", err)
	}
	return &peAcc, nil
}

func NewTxRepository(db *sql.DB) TransactionRepository {
	repo := new(transactionRepository)
	repo.db = db
	return repo
}
