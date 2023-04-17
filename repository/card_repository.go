package repository

import (
	"github.com/ReygaFitra/inc-final-project.git/model"
)

type CardRepository interface {
	FindByID(id uint) (*model.Card, error)
	FindByUserID(userID uint) ([]*model.Card, error)
	Create(card *model.Card) error
	Update(card *model.Card) error
	Delete(card *model.Card) error
}
