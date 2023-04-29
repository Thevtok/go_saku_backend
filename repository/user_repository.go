package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/utils"
	_ "github.com/lib/pq"
)

type UserRepository interface {
	GetByEmailAndPassword(email string, password string) (*model.Credentials, error)
	GetAll() any
	GetByUsername(username string) (*model.UserResponse, error)
	GetByiD(id uint) (*model.User, error)
	Create(user *model.UserCreate) (any, error)
	UpdateProfile(user *model.User) string
	UpdateEmailPassword(user *model.User) string
	Delete(user *model.User) string
	UpdateBalance(userID uint, newBalance uint) error
	UpdatePoint(userID uint, newPoint int) error
}

type userRepository struct {
	db *sql.DB
}

func (r *userRepository) UpdateBalance(userID uint, newBalance uint) error {
	_, err := r.GetByiD(userID)
	if err != nil {
		return err
	}

	query := "UPDATE mst_users SET balance = $1 WHERE user_id = $2"
	_, err = r.db.Exec(query, newBalance, userID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *userRepository) UpdatePoint(userID uint, newPoint int) error {
	fmt.Printf("userID: %d, newPoint: %d\n", userID, newPoint) // print userID and newPoint
	_, err := r.GetByiD(userID)
	if err != nil {
		return err
	}

	query := "UPDATE mst_users SET point = $1 WHERE user_id = $2"
	_, err = r.db.Exec(query, newPoint, userID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *userRepository) GetAll() any {
	var users []model.UserResponse
	query := "SELECT name, username, email, phone_number, address, balance, point from mst_users"
	rows, err := r.db.Query(query)
	if err != nil {
		log.Println(err)
	}
	if rows == nil {
		return "no data"
	}
	defer rows.Close()

	for rows.Next() {
		var user model.UserResponse
		if err := rows.Scan(&user.Name, &user.Username, &user.Email, &user.Phone_Number, &user.Address, &user.Balance, &user.Point); err != nil {
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

func (r *userRepository) GetByUsername(username string) (*model.UserResponse, error) {
	var user model.UserResponse
	err := r.db.QueryRow("SELECT name, username, email, phone_number, address, balance, point FROM mst_users WHERE username = $1", username).Scan(&user.Name, &user.Username, &user.Email, &user.Phone_Number, &user.Address, &user.Balance, &user.Point)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByiD(id uint) (*model.User, error) {
	var user model.User

	err := r.db.QueryRow("SELECT name, user_id, email, phone_number, address, balance, point FROM mst_users WHERE user_id = $1", id).Scan(&user.Name, &user.ID, &user.Email, &user.Phone_Number, &user.Address, &user.Balance, &user.Point)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateProfile(user *model.User) string {
	_, err := r.GetByiD(user.ID)
	if err != nil {
		return "user not found"
	}

	query := "UPDATE mst_users SET name=$1, phone_number=$2, address=$3, username=$4 WHERE user_id=$5"
	_, err = r.db.Exec(query, user.Name, user.Phone_Number, user.Address, user.Username, user.ID)
	if err != nil {
		log.Println(err)
		return "failed to update user"
	}
	return "updated profile successfully"
}

func (r *userRepository) UpdateEmailPassword(user *model.User) string {
	_, err := r.GetByiD(user.ID)
	if err != nil {
		return "user not found"
	}

	query := "UPDATE mst_users SET email=$1, password=$2  WHERE user_id=$3"

	_, err = r.db.Exec(query, user.Email, user.Password, user.ID)
	if err != nil {
		log.Println(err)
		return "failed to update user"
	}
	return "updated email and password successfully"
}

func (r *userRepository) Delete(user *model.User) string {
	// Check if the user exists
	_, err := r.GetByUsername(user.Username)
	if err != nil {
		return "user not found"
	}

	// Execute the delete query
	query := "DELETE FROM mst_users WHERE username = $1"
	_, err = r.db.Exec(query, user.Username)
	if err != nil {
		log.Println(err)
		return "failed to delete user"
	}
	return "deleted user successfully"
}

func (r *userRepository) Create(user *model.UserCreate) (any, error) {
	// Create a copy of the user object
	hashedPassword, err := utils.HasingPassword(user.Password)
	if err != nil {
		log.Println(err)
	}

	user.Password = hashedPassword

	_, err = r.db.Exec("INSERT INTO mst_users (name, username, email, password, phone_number, address, balance, role, point) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", user.Name, user.Username, user.Email, user.Password, user.Phone_Number, user.Address, 0, "user", 0)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByEmailAndPassword(email string, password string) (*model.Credentials, error) {
	var m model.Credentials
	query := "SELECT user_id, username, password, role FROM mst_users WHERE email = $1"
	row := r.db.QueryRow(query, email)
	var hashedPassword string
	err := row.Scan(&m.UserID, &m.Username, &hashedPassword, &m.Role)
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
		Username: m.Username,
		UserID:   m.UserID,
		Role:     m.Role,
	}
	return user, nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}
