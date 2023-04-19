package repository

import (
	"database/sql"
	"fmt"

	"github.com/ReygaFitra/inc-final-project.git/model"
)

type CardRepository interface {
	GetByID(id uint) (*model.Card, error)
	GetByUserID(userID uint) ([]*model.Card, error)
	Create(card *model.Card) error
	Update(card *model.Card) error
	Delete(card *model.Card) error
}

type cardRepository struct {
	db *sql.DB
}

func NewCardRepository(db *sql.DB) CardRepository {
	return &cardRepository{db}
}

func (r *cardRepository) GetByID(id uint) (*model.Card, error) {
	card := &model.Card{}
	err := r.db.QueryRow("SELECT * FROM MST_Card WHERE Card_ID = $1", id).
		Scan(&card.Card_ID, &card.User_ID, &card.Card_Type, &card.Card_Number, &card.Expiration_Date, &card.CVV)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("card with ID %d not found", id)
		}
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	return card, nil
}

func (r *cardRepository) GetByUserID(userID uint) ([]*model.Card, error) {
	rows, err := r.db.Query("SELECT * FROM MST_Card WHERE User_ID = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	cards := []*model.Card{}
	for rows.Next() {
		card := &model.Card{}
		err := rows.Scan(&card.Card_ID, &card.User_ID, &card.Card_Type, &card.Card_Number, &card.Expiration_Date, &card.CVV)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		cards = append(cards, card)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return cards, nil
}

func (r *cardRepository) Create(card *model.Card) error {
	err := r.db.QueryRow("INSERT INTO MST_Card (User_ID, Card_Type, Card_Number, Expiration_Date, CVV) VALUES ($1, $2, $3, $4, $5) RETURNING Card_ID",
		card.User_ID, card.Card_Type, card.Card_Number, card.Expiration_Date, card.CVV).Scan(&card.Card_ID)
	if err != nil {
		return fmt.Errorf("error inserting row: %v", err)
	}
	return nil
}

func (r *cardRepository) Update(card *model.Card) error {
	_, err := r.db.Exec("UPDATE MST_Card SET User_ID=$1, Card_Type=$2, Card_Number=$3, Expiration_Date=$4, CVV=$5 WHERE Card_ID=$6",
		card.User_ID, card.Card_Type, card.Card_Number, card.Expiration_Date, card.CVV, card.Card_ID)
	if err != nil {
		return fmt.Errorf("error updating row: %v", err)
	}
	return nil
}

func (r *cardRepository) Delete(card *model.Card) error {
	_, err := r.db.Exec("DELETE FROM MST_Card WHERE Card_ID=$1", card.Card_ID)
	if err != nil {
		return fmt.Errorf("error deleting row: %v", err)
	}
	return nil
}
