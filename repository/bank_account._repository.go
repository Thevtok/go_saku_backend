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
	GetByID(id uint) (*model.BankAccResponse, error)
	Create(newBankAcc *model.BankAccResponse) (any, error)
	Update(bankAcc *model.BankAcc) string
	Delete(bankAcc *model.BankAcc) string
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

func (r *bankAccRepository) GetByID(id uint) (*model.BankAccResponse, error) {
	var bankAcc model.BankAccResponse
	query := "SELECT user_id, bank_name, account_number, account_holder_name FROM mst_bank_account   WHERE user_id = $1"
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

func (r *bankAccRepository) Create(newBankAcc *model.BankAccResponse) (any, error) {
	query := "INSERT INTO mst_bank_account (user_id, bank_name, account_number, account_holder_name) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(query, newBankAcc.UserID, newBankAcc.BankName, newBankAcc.AccountNumber, newBankAcc.AccountHolderName)
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

	query := "UPDATE mst_bank_account SET bank_name = $1, account_number = $2, account_holder_name = $3 WHERE user_id = $4"

	_, err = r.db.Exec(query, bankAcc.BankName, bankAcc.AccountNumber, bankAcc.AccountHolderName, bankAcc.UserID)
	if err != nil {
		log.Println(err)
		return "failed to update Bank Account"
	}

	return "Bank Account updated Successfully"
}

func (r *bankAccRepository) Delete(bankAcc *model.BankAcc) string {
	_, err := r.GetByID(bankAcc.UserID)
	if err != nil {
		return "user not found"
	}

	query := "DELETE FROM mst_bank_account WHERE user_id = $1"
	_, err = r.db.Exec(query, bankAcc.UserID)
	if err != nil {
		log.Println(err)
		return "failed to delete Bank Account"
	}

	return "Bank Account deleted Successfully"
}

func NewBankAccRepository(db *sql.DB) BankAccRepository {
	repo := new(bankAccRepository)
	repo.db = db
	return repo
}
