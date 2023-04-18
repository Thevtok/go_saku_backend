package repository

import (
	"database/sql"
	"log"

	"github.com/ReygaFitra/inc-final-project.git/model"
)

type BankAccRepository interface {
	GetAll() any
	GetByID(id int) any
	Create(newBankAcc *model.BankAcc) string
	Update(bankacc *model.BankAcc) string
	Delete(id int) string
}

type bankAccRepository struct {
	db *sql.DB
}

func NewBankAccRepository(db *sql.DB) BankAccRepository {
	repo := new(bankAccRepository)
	repo.db = db
	return repo
}

func (r *bankAccRepository) GetAll() any {
	var users []model.BankAcc

	query := "SELECT account_id, user_id, bank_account, account_number, account_holder_name FROM mst_bank_account"
	rows, err := r.db.Query(query)

	if err != nil {
		log.Println(err)
	}
	if rows == nil {
		return "no data"
	}
	defer rows.Close()

	for rows.Next() {
		var user model.BankAcc

		err := rows.Scan(&user.AccountId, &user.UserId, &user.BankName, &user.AccountNumber, &user.AccountHolderName)
		if err != nil {
			log.Println(err)
		}

		users = append(users, user)
	}

	return users
}

func (r *bankAccRepository) GetByID(id int) any {
	var user model.BankAcc

	query := "SELECT account_id, bank_account, account_number, account_holder_name WHERE user_id = $1"
	row := r.db.QueryRow(query, id)
	err := row.Scan(&user.AccountId, &user.BankName, &user.AccountNumber, &user.AccountHolderName)
	if err != nil {
		log.Println(err)
	}

	return user
}
func (r *bankAccRepository) Create(newBankAcc *model.BankAcc) string {
	query := "INSERT INTO mst_bank_account (account_id, user_id, bank_account, account_number, account_holder_name) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(query, newBankAcc.AccountId, newBankAcc.UserId, newBankAcc.BankName, newBankAcc.AccountNumber, newBankAcc.AccountHolderName)
	if err != nil {
		log.Println(err)
		return "failed to create data"
	}

	return "Bank Account created Successfully"
}

func (r *bankAccRepository) Update(bankAcc *model.BankAcc) string {
	result := r.GetByID(bankAcc.UserId)

	if result == "Bank Account not found" {
		return result.(string)
	}

	query := "UPDATE mst_bank_account SET account_id = $1, bank_account = $2, account_number = $3, account_holder_name = $4 WHERE user_id = $5"
	_, err := r.db.Exec(query, bankAcc.AccountId, bankAcc.BankName, bankAcc.AccountNumber, bankAcc.AccountHolderName)
	if err != nil {
		log.Println(err)
		return "failed to update Bank Account"
	}

	return "Bank Account updated Successfully"
}

func (r *bankAccRepository) Delete(id int) string {
	query := "DELETE FROM mst_bank_account WHERE user_id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Println(err)
		return "failed to delete Bank Account"
	}

	return "Bank Account deleted Successfully"
}
