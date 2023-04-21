package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ReygaFitra/inc-final-project.git/model"
)

type BankAccRepository interface {
	GetAll() any
	GetByID(id uint) ([]*model.BankAccResponse, error)
	GetByAccountID(id uint) (*model.BankAcc, error)
	Create(userID uint, newBankAcc *model.BankAccResponse) (any, error)
	Update(bankAcc *model.BankAcc) string
	DeleteByUserID(id uint) string
	DeleteByAccountId(accountID uint) error
}

type bankAccRepository struct {
	db *sql.DB
}

func (r *bankAccRepository) GetAll() any {

	var users []model.BankAccResponse

	query := "SELECT  b.bank_name, b.account_number, b.account_holder_name FROM mst_bank_account b JOIN mst_users u ON b.account_id = u.user_id"

	rows, err := r.db.Query(query)

	if err != nil {
		log.Println(err)
	}
	if rows == nil {
		return "no data"
	}
	defer rows.Close()

	for rows.Next() {

		var user model.BankAccResponse

		err := rows.Scan(&user.BankName, &user.AccountNumber, &user.AccountHolderName)
		if err != nil {
			log.Println(err)
		}

		users = append(users, user)
	}

	return users
}

func (r *bankAccRepository) GetByID(id uint) ([]*model.BankAccResponse, error) {
	var bankAccs []*model.BankAccResponse
	query := "SELECT user_id, bank_name, account_number, account_holder_name FROM mst_bank_account WHERE user_id = $1"
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bankAcc model.BankAccResponse
		err = rows.Scan(&bankAcc.UserID, &bankAcc.BankName, &bankAcc.AccountNumber, &bankAcc.AccountHolderName)
		if err != nil {
			return nil, err
		}
		bankAccs = append(bankAccs, &bankAcc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bankAccs, nil
}
func (r *bankAccRepository) GetByAccountID(id uint) (*model.BankAcc, error) {
	var bankAcc model.BankAcc
	query := "SELECT account_id, bank_name, account_number, account_holder_name FROM mst_bank_account   WHERE account_id = $1"
	row := r.db.QueryRow(query, id)
	err := row.Scan(&bankAcc.UserID, &bankAcc.BankName, &bankAcc.AccountNumber, &bankAcc.AccountHolderName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("bank account not found")
		}
		return nil, err
	}
	return &bankAcc, nil
}

func (r *bankAccRepository) Create(userID uint, newBankAcc *model.BankAccResponse) (any, error) {
	query := "INSERT INTO mst_bank_account (user_id, bank_name, account_number, account_holder_name) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(query, userID, newBankAcc.BankName, newBankAcc.AccountNumber, newBankAcc.AccountHolderName)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to create data")
	}

	return newBankAcc, nil
}

func (r *bankAccRepository) Update(bankAcc *model.BankAcc) string {
	_, err := r.GetByID(bankAcc.UserID)
	if err != nil {
		return "user not found"
	}

	query := "UPDATE mst_bank_account SET bank_name = $1, account_number = $2, account_holder_name = $3 WHERE account_id = $4"

	_, err = r.db.Exec(query, bankAcc.BankName, bankAcc.AccountNumber, bankAcc.AccountHolderName, bankAcc.UserID)
	if err != nil {
		log.Println(err)
		return "failed to update Bank Account"
	}

	return "Bank Account updated Successfully"
}

func (r *bankAccRepository) DeleteByUserID(id uint) string {
	query := "DELETE FROM mst_bank_account WHERE user_id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return "failed to delete bank"
	}

	return "deleted all bank account successfully"
}

func (r *bankAccRepository) DeleteByAccountId(accountID uint) error {
	_, err := r.GetByID(accountID)
	if err != nil {
		return err
	}

	query := "DELETE FROM mst_bank_account WHERE account_id = $1"

	_, err = r.db.Exec(query, accountID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
func NewBankAccRepository(db *sql.DB) BankAccRepository {
	repo := new(bankAccRepository)
	repo.db = db
	return repo
}
