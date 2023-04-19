package repository

import (
	"database/sql"
	"log"

	"github.com/ReygaFitra/inc-final-project.git/model"
)

type BankAccRepository interface {
	GetAll() any
	GetByID(id uint) any
	Create(newBankAcc *model.BankAcc) string
	Update(bankAcc *model.BankAcc) string
	Delete(id uint) string
}

type bankAccRepository struct {
	db *sql.DB
}

func (r *bankAccRepository) GetAll() any {
	var users []model.BankResponse

	query := "SELECT u.name , b.bank_name, b.account_number, b.account_holder_name FROM mst_users u join mst_bank_account b ON b.user_id = u.user_id "
	rows, err := r.db.Query(query)

	if err != nil {
		log.Println(err)
	}
	if rows == nil {
		return "no data"
	}
	defer rows.Close()

	for rows.Next() {
		var user model.BankResponse

		err := rows.Scan(&user.Name, &user.BankName, &user.AccountNumber, &user.AccountHolderName)
		if err != nil {
			log.Println(err)
		}

		users = append(users, user)
	}

	return users
}

func (r *bankAccRepository) GetByID(id uint) any {
	var user model.BankResponse

	query := "SELECT u.name , b.bank_name, b.account_number, b.account_holder_name FROM mst_users u join mst_bank_account b ON b.user_id = u.user_id  WHERE account_id = $1"
	row := r.db.QueryRow(query, id)
	err := row.Scan(&user.Name, &user.BankName, &user.AccountNumber, &user.AccountHolderName)
	if err != nil {
		log.Println(err)
	}

	return user
}
func (r *bankAccRepository) Create(newBankAcc *model.BankAcc) string {
	query := "INSERT INTO mst_bank_account (user_id, bank_name, account_number, account_holder_name) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(query, newBankAcc.UserId, newBankAcc.BankName, newBankAcc.AccountNumber, newBankAcc.AccountHolderName)
	if err != nil {
		log.Println(err)
		return "failed to create data"
	}

	return "Bank Account created Successfully"
}

func (r *bankAccRepository) Update(bankAcc *model.BankAcc) string {
	result := r.GetByID(bankAcc.AccountId)

	if result == "Bank Account not found" {
		return result.(string)
	}

	query := "UPDATE mst_bank_account SET user_id = $1, bank_name = $2, account_number = $3, account_holder_name = $4 WHERE user_id = $5"
	_, err := r.db.Exec(query, bankAcc.UserId, bankAcc.BankName, bankAcc.AccountNumber, bankAcc.AccountHolderName)
	if err != nil {
		log.Println(err)
		return "failed to update Bank Account"
	}

	return "Bank Account updated Successfully"
}

func (r *bankAccRepository) Delete(id uint) string {
	result := r.GetByID(id)

	if result == "Bank Account not found" {
		return result.(string)
	}

	query := "DELETE FROM mst_bank_account WHERE user_id = $1"
	_, err := r.db.Exec(query, id)
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
