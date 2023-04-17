package repository

import (
	"github.com/ReygaFitra/inc-final-project.git/model"
)

type BankAccountRepository interface {
	FindByID(id uint) (*model.BankAccount, error)
	FindByUserID(userID uint) ([]*model.BankAccount, error)
	Create(bankAccount *model.BankAccount) error
	Update(bankAccount *model.BankAccount) error
	Delete(bankAccount *model.BankAccount) error
}
