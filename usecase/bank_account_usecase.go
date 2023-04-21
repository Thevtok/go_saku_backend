package usecase

import (
	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/repository"
)

type BankAccUsecase interface {
	FindAllBankAcc() any
	FindBankAccByID(id uint) (*model.BankAccResponse, error)
	Register(newBankAcc *model.BankAccResponse) (any, error)
	Edit(bankAcc *model.BankAcc) string
	Unreg(bankAcc *model.BankAcc) string
}

type bankAccUsecase struct {
	bankAccRepo repository.BankAccRepository
}

func (u *bankAccUsecase) FindAllBankAcc() any {
	return u.bankAccRepo.GetAll()
}

func (u *bankAccUsecase) FindBankAccByID(id uint) (*model.BankAccResponse, error) {
	return u.bankAccRepo.GetByID(id)
}

func (u *bankAccUsecase) Register(newBankAcc *model.BankAccResponse) (any, error) {
	return u.bankAccRepo.Create(newBankAcc)
}

func (u *bankAccUsecase) Edit(bankAcc *model.BankAcc) string {
	return u.bankAccRepo.Update(bankAcc)
}

func (u *bankAccUsecase) Unreg(bankAcc *model.BankAcc) string {
	return u.bankAccRepo.Delete(bankAcc)
}

func NewBankAccUsecase(bankAccRepo repository.BankAccRepository) BankAccUsecase {
	return &bankAccUsecase{
		bankAccRepo: bankAccRepo,
	}
}
