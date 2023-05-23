package usecase

import (
	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/repository"
)

type BankAccUsecase interface {
	FindBankAccByUserID(id string) ([]*model.BankAcc, error)
	FindBankAccByAccountID(id uint) (*model.BankAcc, error)
	Register(id string, newBankAcc *model.BankAcc) (any, error)

	UnregByAccountID(accountID uint) error
}

type bankAccUsecase struct {
	bankAccRepo repository.BankAccRepository
}

func (u *bankAccUsecase) FindBankAccByUserID(id string) ([]*model.BankAcc, error) {
	return u.bankAccRepo.GetByUserID(id)
}
func (u *bankAccUsecase) FindBankAccByAccountID(id uint) (*model.BankAcc, error) {
	return u.bankAccRepo.GetByAccountID(id)
}

func (u *bankAccUsecase) Register(id string, newBankAcc *model.BankAcc) (any, error) {
	newBankAcc.UserID = id
	return u.bankAccRepo.Create(id, newBankAcc)
}

func (u *bankAccUsecase) UnregByAccountID(accountID uint) error {
	err := u.bankAccRepo.DeleteByAccountID(accountID)
	if err != nil {
		return err
	}
	return nil
}

func NewBankAccUsecase(bankAccRepo repository.BankAccRepository) BankAccUsecase {
	return &bankAccUsecase{
		bankAccRepo: bankAccRepo,
	}
}
