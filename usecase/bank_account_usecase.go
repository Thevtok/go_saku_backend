package usecase

import (
	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/repository"
)

type BankAccUsecase interface {
	FindAllBankAcc() any
	FindBankAccByID(id int) any
	Register(newbankAcc *model.BankAcc) string
	Edit(bankAcc *model.BankAcc) string
	Unreg(id int) string
}

type bankAccUsecase struct {
	bankAccRepo repository.BankAccRepository
}

func NewBankAccUsecase(bankAccRepo repository.BankAccRepository) BankAccUsecase {
	return &bankAccUsecase{
		bankAccRepo: bankAccRepo,
	}
}

func (u *bankAccUsecase) FindAllBankAcc() any {
	return u.bankAccRepo.GetAll()
}

func (u *bankAccUsecase) FindBankAccByID(userID int) any {
	return u.bankAccRepo.GetByID(userID)
}

func (u *bankAccUsecase) Register(newbankAcc *model.BankAcc) string {
	return u.bankAccRepo.Create(newbankAcc)
}

func (u *bankAccUsecase) Edit(bankAcc *model.BankAcc) string {
	return u.bankAccRepo.Update(bankAcc)
}

func (u *bankAccUsecase) Unreg(userId int) string {
	return u.bankAccRepo.Delete(userId)
}
