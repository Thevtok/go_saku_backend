package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/utils"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var newUUID = uuid.New()

type UserRepository interface {
	GetByEmailAndPassword(email string, password string, token string) (*model.User, error)
	GetAll() any
	GetByUsername(username string) (*model.User, error)
	GetByiD(id string) (*model.User, error)
	Create(user *model.User) (any, error)
	UpdateProfile(user *model.User) string
	UpdateEmailPassword(user *model.User) string
	Delete(user *model.User) string
	UpdateBalance(userID string, newBalance int) error
	UpdatePoint(userID string, newPoint int) error
	GetByPhone(phoneNumber string) (*model.User, error)
	SaveDeviceToken(userID string, token string) error
	GetByIDToken(id string) (*model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func (r *userRepository) SaveDeviceToken(userID string, token string) error {
	_, err := r.db.Exec("UPDATE mst_users SET token = $1 WHERE user_id = $2", token, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) UpdateBalance(userID string, newBalance int) error {
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

func (r *userRepository) UpdatePoint(userID string, newPoint int) error {
	fmt.Printf("userID: %s, newPoint: %d\n", userID, newPoint) // print userID and newPoint
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
	var users []model.User
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
		var user model.User
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

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(`
		SELECT u.user_id, u.name, u.username, u.email, u.phone_number, u.address, u.balance, u.point, b.badge_name
		FROM mst_users u
		JOIN mst_badges b ON u.badge_id = b.badge_id
		WHERE u.username = $1
	`, username).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Phone_Number,
		&user.Address,
		&user.Balance,
		&user.Point,
		&user.Badge,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("username not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByPhone(phoneNumber string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow("SELECT user_id,name, username, email, phone_number, address, balance, point,badge_id,tx_count  FROM mst_users WHERE phone_number = $1", phoneNumber).Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Phone_Number, &user.Address, &user.Balance, &user.Point, &user.BadgeID, &user.TxCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("phone not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByiD(id string) (*model.User, error) {
	var user model.User

	err := r.db.QueryRow("SELECT name, user_id, email, phone_number, address, balance, username, point,badge_id,tx_count FROM mst_users WHERE user_id = $1", id).Scan(&user.Name, &user.ID, &user.Email, &user.Phone_Number, &user.Address, &user.Balance, &user.Username, &user.Point, &user.BadgeID, &user.TxCount)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("id not found")
		}
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) GetByIDToken(id string) (*model.User, error) {
	var user model.User

	// Mengambil data pengguna berdasarkan ID
	err := r.db.QueryRow(`
		SELECT name, user_id, email, phone_number, address, balance, username, point,token,badge_id,tx_count
		FROM mst_users
		WHERE user_id = $1
	`, id).Scan(&user.Name, &user.ID, &user.Email, &user.Phone_Number, &user.Address, &user.Balance, &user.Username, &user.Point, &user.Token, &user.BadgeID, &user.TxCount)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("id not found")
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

	query := "UPDATE mst_users SET email=$1, password=$2 WHERE user_id=$3"

	_, err = r.db.Exec(query, user.Email, user.Password, user.ID)
	if err != nil {
		log.Println(err)
		return "failed to update user"
	}
	return "updated email and password successfully"
}

func (r *userRepository) Delete(user *model.User) string {
	// Check if the user exists
	_, err := r.GetByiD(user.ID)
	if err != nil {
		return "user not found"
	}

	// Execute the delete query
	query := "DELETE FROM mst_users WHERE user_id = $1"
	_, err = r.db.Exec(query, user.ID)
	if err != nil {
		log.Println(err)
		return "failed to delete user"
	}
	return "deleted user successfully"
}

func (r *userRepository) Create(user *model.User) (any, error) {
	// Create a copy of the user object
	hashedPassword, err := utils.HasingPassword(user.Password)
	if err != nil {
		log.Println(err)
	}

	user.Password = hashedPassword

	_, err = r.db.Exec("INSERT INTO mst_users (user_id,name, username, email, password, phone_number, address, balance, role, point, token,badge_id,tx_count) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,$10,$11,$12,$13)", newUUID, user.Name, user.Username, user.Email, user.Password, user.Phone_Number, user.Address, 0, "user", 0, "", 1, 0)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return "user created successfully", nil
}
func (r *userRepository) GetByEmailAndPassword(email string, password string, token string) (*model.User, error) {
	var m model.User
	query := "SELECT user_id, username, password, role FROM mst_users WHERE email = $1"
	row := r.db.QueryRow(query, email)
	var hashedPassword string
	err := row.Scan(&m.ID, &m.Username, &hashedPassword, &m.Role)
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
		return nil, fmt.Errorf("invalid credentials \n password = %s\n hashed = %s", password, hashedPassword)
	}

	user := &model.User{
		Password: hashedPassword,
		Username: m.Username,
		ID:       m.ID,
		Role:     m.Role,
	}

	// Save the device token for the user
	err = r.SaveDeviceToken(m.ID, token)
	if err != nil {
		return nil, fmt.Errorf("failed to save device token: %v", err)
	}

	return user, nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}
