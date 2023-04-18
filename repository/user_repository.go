package repository

import (
	"database/sql"
	"log"

	"github.com/ReygaFitra/inc-final-project.git/model"
)


type UserRepository interface {
	FindByID(id uint) (*model.Users, error)
	FindByEmail(email string) (*model.Users, error)
	FindByPhone(phone string) (*model.Users, error)
	Create(user *model.Users) error
	Update(user *model.Users) error
	Delete(user *model.Users) error
}

type UserRepo interface {
	Register(newUser *model.Users) string
}

type userRepo struct {
	db *sql.DB
}

func (r *userRepo) Create(newUser *model.Users) string {
	query := "INSERT INTO (name, email, password, phone_number, address, balance) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := r.db.Exec(query, newUser.Name, newUser.Email, newUser.Password, newUser.Phone_Number, newUser.Address, newUser.Balance) 
	if err != nil {
		log.Println(err.Error())
		return "Failed Register"
	}
	return "Register Successfully"
}

