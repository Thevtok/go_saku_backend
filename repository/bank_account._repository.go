package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ReygaFitra/inc-final-project.git/model"
)

type BankAccRepository interface {
	GetByUserID(id string) ([]*model.BankAcc, error)
	GetByAccountID(id uint) (*model.BankAcc, error)
	Create(id string, newBankAcc *model.BankAcc) (any, error)

	DeleteByAccountID(accountID uint) error
}

type bankAccRepository struct {
	db *sql.DB
}

func (r *bankAccRepository) GetByUserID(id string) ([]*model.BankAcc, error) {
	var bankAccs []*model.BankAcc
	query := "SELECT user_id, account_id, bank_name, account_number, account_holder_name FROM mst_bank_account WHERE user_id = $1"
	rows, err := r.db.Query(query, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bankAcc model.BankAcc
		err = rows.Scan(&bankAcc.UserID, &bankAcc.AccountID, &bankAcc.BankName, &bankAcc.AccountNumber, &bankAcc.AccountHolderName)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		bankAccs = append(bankAccs, &bankAcc)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return bankAccs, nil
}

func (r *bankAccRepository) GetByAccountID(id uint) (*model.BankAcc, error) {
	var bankAcc model.BankAcc
	query := "SELECT account_id, bank_name, account_number, account_holder_name, user_id FROM mst_bank_account WHERE account_id = $1"
	row := r.db.QueryRow(query, id)
	err := row.Scan(&bankAcc.AccountID, &bankAcc.BankName, &bankAcc.AccountNumber, &bankAcc.AccountHolderName, &bankAcc.UserID)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return nil, errors.New("bank account not found")
		}
		return nil, err
	}
	return &bankAcc, nil
}

func (r *bankAccRepository) Create(id string, newBankAcc *model.BankAcc) (any, error) {
	query := "INSERT INTO mst_bank_account (user_id, bank_name, account_number, account_holder_name) VALUES ($1, $2, $3, $4) RETURNING account_id"
	_, err := r.db.Exec(query, id, newBankAcc.BankName, newBankAcc.AccountNumber, newBankAcc.AccountHolderName)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to create data")
	}
	return newBankAcc, nil
}

func (r *bankAccRepository) DeleteByAccountID(accountID uint) error {
	_, err := r.GetByAccountID(accountID)
	if err != nil {
		log.Println(err)
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
