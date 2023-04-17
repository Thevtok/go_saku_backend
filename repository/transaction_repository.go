package repository

import (
	"github.com/ReygaFitra/inc-final-project.git/model"
)

type TransactionRepository interface {
	FindByID(id uint) (*model.Transaction, error)
	FindByUserID(userID uint) ([]*model.Transaction, error)
	Create(transaction *model.Transaction) error
	Update(transaction *model.Transaction) error
	Delete(transaction *model.Transaction) error
}
