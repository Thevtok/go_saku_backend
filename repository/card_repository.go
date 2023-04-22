package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ReygaFitra/inc-final-project.git/model"
)

type CardRepository interface {
	GetAll() any
	GetByID(id uint) ([]*model.CardResponse, error)
	GetByCardID(id uint) (*model.Card, error)
	Create(userID uint, newCard *model.CardResponse) (any, error)
	Update(card *model.Card) string
	DeleteByUserID(id uint) string
	DeleteByCardId(cardID uint) error
}

type cardRepository struct {
	db *sql.DB
}

func (r *cardRepository) GetAll() any {

	var users []model.BankAccResponse

	query := "SELECT  bank_name, account_number, account_holder_name,username FROM mst_bank_account"

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

		err := rows.Scan(&user.BankName, &user.AccountNumber, &user.AccountHolderName, &user.Username)
		if err != nil {
			log.Println(err)
		}

		users = append(users, user)
	}

	return users
}

func (r *cardRepository) GetByID(id uint) ([]*model.CardResponse, error) {
	var cards []*model.CardResponse
	query := "SELECT user_id, card_type, card_number, expiration_date, cvv FROM mst_card WHERE user_ID = $1"
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var card model.CardResponse
		err = rows.Scan(&card.UserID, &card.CardType, &card.CardNumber, &card.ExpirationDate, &card.CVV)
		if err != nil {
			return nil, err
		}
		cards = append(cards, &card)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

func (r *cardRepository) GetByCardID(id uint) (*model.Card, error) {
	var card model.Card
	query := "SELECT card_id, card_type, card_number, expiration_date, cvv WHERE card_id = $1"
	row := r.db.QueryRow(query, id)
	err := row.Scan(&card.UserID, &card.CardType, &card.CardNumber, &card.ExpirationDate, &card.CVV)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("card id not founf")
		}
		return nil, err
	}
	return &card, nil
}

func (r *cardRepository) Create(userID uint, newCard *model.CardResponse) (any, error) {
	query := "INSERT INTO mst_card (user_ID, card_type, card_number, expiration_date, cvv) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(query, userID, newCard.CardType, newCard.CardNumber, newCard.ExpirationDate, newCard.CVV)
	if err != nil {
		return nil, fmt.Errorf("failed to create data")
	}
	return newCard, nil
}

func (r *cardRepository) Update(card *model.Card) string {
	_, err := r.GetByID(card.UserID)
	if err != nil {
		return "user not found"
	}

	query := "UPDATE mst_card SET card_type = $1, card_number = $2, expiration_date = $3, cvv = $4 WHERE card_id = $5"

	_, err = r.db.Exec(query, card.CardType, card.CardNumber, card.ExpirationDate, card.CVV, card.UserID)
	if err != nil {
		log.Println("failed to update Card ID")
	}

	return "Card ID updated Successfully"
}

func (r *cardRepository) DeleteByUserID(id uint) string {
	query := "DELETE FROM mst_card WHERE user_id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return "failed to delete card"
	}

	return "Deleted All Card ID Successfully"
}

func (r *cardRepository) DeleteByCardId(cardID uint) error {
	_, err := r.GetByID(cardID)
	if err != nil {
		return err
	}

	query := "DELETE FROM mst_card WHERE card_id = $1"

	_, err = r.db.Exec(query, cardID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func NewCardRepository(db *sql.DB) CardRepository {
	return &cardRepository{db}
}
