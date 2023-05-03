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
	GetByUserID(id uint) ([]*model.CardResponse, error)
	GetByCardID(id uint) (*model.Card, error)
	Create(id uint, newCard *model.CardResponse) (any, error)
	Update(card *model.Card) string
	DeleteByUserID(userID uint) string
	DeleteByCardID(cardID uint) error
}

type cardRepository struct {
	db *sql.DB
}

func (r *cardRepository) GetAll() any {
	var users []model.CardResponse
	query := "SELECT user_id, card_id, card_type, card_number, expiration_date, cvv FROM mst_card"
	rows, err := r.db.Query(query)
	if err != nil {
		log.Println(err)
	}
	if rows == nil {
		return "no data"
	}
	defer rows.Close()

	for rows.Next() {
		var user model.CardResponse
		err := rows.Scan(&user.UserID, &user.CardID, &user.CardType, &user.CardNumber, &user.ExpirationDate, &user.CVV)
		if err != nil {
			log.Println(err)
		}
		users = append(users, user)
	}
	return users
}

func (r *cardRepository) GetByUserID(id uint) ([]*model.CardResponse, error) {
	var cards []*model.CardResponse
	query := "SELECT user_id, card_id, card_type, card_number, expiration_date, cvv FROM mst_card WHERE user_id = $1"
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var card model.CardResponse
		err = rows.Scan(&card.UserID, &card.CardID, &card.CardType, &card.CardNumber, &card.ExpirationDate, &card.CVV)
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
	query := "SELECT card_id, card_type, card_number, expiration_date, cvv, user_id FROM mst_card WHERE card_id = $1"
	row := r.db.QueryRow(query, id)
	err := row.Scan(&card.CardID, &card.CardType, &card.CardNumber, &card.ExpirationDate, &card.CVV, &card.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("card is not found")
		}
		return nil, err
	}
	return &card, nil
}

func (r *cardRepository) Create(id uint, newCard *model.CardResponse) (any, error) {
	query := "INSERT INTO mst_card (user_id, card_type, card_number, expiration_date, cvv) VALUES ($1, $2, $3, $4, $5) 	RETURNING card_id"
	_, err := r.db.Exec(query, id, newCard.CardType, newCard.CardNumber, newCard.ExpirationDate, newCard.CVV)
	if err != nil {
		return nil, fmt.Errorf("failed to create data")
	}
	return newCard, nil
}

func (r *cardRepository) Update(card *model.Card) string {
	_, err := r.GetByUserID(card.UserID)
	if err != nil {
		return "user not found"
	}

	query := "UPDATE mst_card SET card_type = $1, card_number = $2, expiration_date = $3, cvv = $4 WHERE card_id = $5"
	_, err = r.db.Exec(query, card.CardType, card.CardNumber, card.ExpirationDate, card.CVV, card.CardID)
	if err != nil {
		log.Println("failed to update Card ID")
	}

	return "Card ID updated Successfully"
}

func (r *cardRepository) DeleteByUserID(userID uint) string {
	query := "DELETE FROM mst_card WHERE user_id = $1"
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return "failed to delete card"
	}
	return "Deleted All Card ID Successfully"
}

func (r *cardRepository) DeleteByCardID(cardID uint) error {
	_, err := r.GetByCardID(cardID)
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
	repo := new(cardRepository)
	repo.db = db
	return repo
}
