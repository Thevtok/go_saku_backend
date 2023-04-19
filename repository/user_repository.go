package repository

import (
	"database/sql"

	"fmt"

	"log"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/utils"

	_ "github.com/lib/pq"
)

type UserRepository interface {
	GetByUsernameAndPassword(username string, password string) (*model.Credentials, error)
	GetAll() any
	GetByID(id uint) any

	Create(user *model.User) (any, error)
	Update(user *model.User) string
	Delete(id uint) string
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}

}

func (r *userRepository) GetAll() any {
	var users []model.UserGetAll

	query := `SELECT name,email,phone_number,address,balance from mst_users`
	rows, err := r.db.Query(query)

	if err != nil {
		log.Println(err)
	}
	if rows == nil {
		return "no data"
	}
	defer rows.Close()

	for rows.Next() {
		var user model.UserGetAll

		if err := rows.Scan(&user.Name, &user.Email, &user.Phone_Number, &user.Address, &user.Balance); err != nil {
			log.Println(err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
	}

	if len(users) == 0 {
		return "no data"
	}

	return users
}

func (r *userRepository) GetByID(id uint) any {
	var user model.UserGetAll
	row := r.db.QueryRow("SELECT name,email,phone_number,address,balance from mst_users WHERE user_id = $1", id)
	err := row.Scan(&user.Name, &user.Email, &user.Phone_Number, &user.Address, &user.Balance)
	if err != nil {

		log.Println(err)
	}
	return user
}

func (r *userRepository) Update(user *model.User) string {
	res := r.GetByID(user.ID)

	if res == "user not found" {
		return res.(string)
	}

	_, err := r.db.Exec("UPDATE mst_users SET name=$1, email=$2, password=$3, phone_number=$4, address=$5, balance=$6 WHERE user_id=$7",
		user.Name, user.Email, user.Password, user.Phone_Number, user.Address, user.Balance, user.ID)
	if err != nil {
		log.Println(err)
		return "failed to update user"
	}
	return "updated user successfully"
}

func (r *userRepository) Delete(id uint) string {
	query := "DELETE FROM mst_users WHERE user_id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		fmt.Println("failed to delete student")
	}
	return "deleted user successfully"
}

func (r *userRepository) Create(user *model.User) (interface{}, error) {
	// Create a copy of the user object
	hashedPassword, err := utils.HasingPassword(user.Password)
	if err != nil {
		log.Println(err)
	}

	user.Password = hashedPassword

	_, err = r.db.Exec("INSERT INTO mst_users (name, email, password, phone_number, address, balance) VALUES ($1, $2, $3, $4, $5, $6)", user.Name, user.Email, user.Password, user.Phone_Number, user.Address, user.Balance)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByUsernameAndPassword(email string, password string) (*model.Credentials, error) {
	// Query the database to retrieve the user's hashed password by email
	query := "SELECT password FROM mst_users WHERE email = $1"
	row := r.db.QueryRow(query, email)
	var hashedPassword string
	err := row.Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		log.Println(err)
		return nil, fmt.Errorf("failed to get user")
	}

	// Verify that the retrieved password is a valid hash
	err = utils.CheckPasswordHash(password, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials \n password = %s\n hased = %s", password, hashedPassword)
	}

	user := &model.Credentials{
		Password: hashedPassword,
	}

	return user, nil
}
