package usecase

import (
	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/repository"
)

type BankAccUsecase interface {
	FindAllBankAcc() any
	FindBankAccByID(id uint) any
	Register(newBankAcc *model.BankAcc) string
	Edit(bankAcc *model.BankAcc) string
	Unreg(id uint) string
}

type bankAccUsecase struct {
	bankAccRepo repository.BankAccRepository
}

func (u *bankAccUsecase) FindAllBankAcc() any {
	return u.bankAccRepo.GetAll()
}

func (u *bankAccUsecase) FindBankAccByID(id uint) any {
	return u.bankAccRepo.GetByID(id)
}

func (u *bankAccUsecase) Register(newBankAcc *model.BankAcc) string {
	return u.bankAccRepo.Create(newBankAcc)
}

func (u *bankAccUsecase) Edit(bankAcc *model.BankAcc) string {
	return u.bankAccRepo.Update(bankAcc)
}

func (u *bankAccUsecase) Unreg(id uint) string {
	return u.bankAccRepo.Delete(id)
}

func NewBankAccUsecase(bankAccRepo repository.BankAccRepository) BankAccUsecase {
	return &bankAccUsecase{
		bankAccRepo: bankAccRepo,
	}
}
