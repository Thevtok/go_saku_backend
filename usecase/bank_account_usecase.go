package usecase

import (
	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/repository"
)

type BankAccUsecase interface {
	FindAllBankAcc() any
	FindBankAccByUserID(id uint) (any, error)
	FindBankAccByAccountID(id uint) (*model.BankAcc, error)
	Register(id uint, newBankAcc *model.BankAccResponse) (any, error)
	Edit(bankAcc *model.BankAcc) string
	UnregAll(bankAcc *model.BankAcc) string
	UnregByAccountID(accountID uint) error
}

type bankAccUsecase struct {
	bankAccRepo repository.BankAccRepository
}

func (u *bankAccUsecase) FindAllBankAcc() any {
	return u.bankAccRepo.GetAll()
}

func (u *bankAccUsecase) FindBankAccByUserID(id uint) (any, error) {
	return u.bankAccRepo.GetByUserID(id)
}
func (u *bankAccUsecase) FindBankAccByAccountID(id uint) (*model.BankAcc, error) {
	return u.bankAccRepo.GetByAccountID(id)
}

func (u *bankAccUsecase) Register(id uint, newBankAcc *model.BankAccResponse) (any, error) {
	newBankAcc.UserID = id
	return u.bankAccRepo.Create(id, newBankAcc)
}

func (u *bankAccUsecase) Edit(bankAcc *model.BankAcc) string {
	return u.bankAccRepo.Update(bankAcc)
}

func (u *bankAccUsecase) UnregAll(bankAcc *model.BankAcc) string {
	return u.bankAccRepo.DeleteByUserID(bankAcc.UserID)
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
