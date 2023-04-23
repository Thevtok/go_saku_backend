package usecase

import (
	"errors"
	"fmt"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/repository"
)

type TransactionUseCase interface {
	CreateDepositBank(transaction *model.TransactionBank) error
	CreateDepositCard(transaction *model.TransactionCard) error
	CreateWithdrawal(transaction *model.TransactionBank) error
	CreateTransfer(sender *model.User, recipient *model.User, amount uint) (any, error)
	CreateRedeem(transaction *model.TransactionPoint) error
}

type transactionUseCase struct {
	transactionRepo repository.TransactionRepository
	userRepo        repository.UserRepository
}

func NewTransactionUseCase(transactionRepo repository.TransactionRepository, userRepo repository.UserRepository) TransactionUseCase {
	return &transactionUseCase{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
	}
}

func (uc *transactionUseCase) CreateDepositBank(transaction *model.TransactionBank) error {
	user, err := uc.userRepo.GetByiD(transaction.SenderID)
	if err != nil {
		return fmt.Errorf("failed to get user data: %v", err)
	}

	// update user balance
	newBalance := user.Balance + transaction.Amount
	err = uc.userRepo.UpdateBalance(user.ID, newBalance)
	if err != nil {
		return fmt.Errorf("failed to update user balance: %v", err)
	}

	// insert transaction
	err = uc.transactionRepo.CreateDepositBank(transaction)
	if err != nil {
		return fmt.Errorf("failed to create deposit transaction: %v", err)
	}

	return nil
}

func (uc *transactionUseCase) CreateDepositCard(transaction *model.TransactionCard) error {
	user, err := uc.userRepo.GetByiD(transaction.SenderID)
	if err != nil {
		return fmt.Errorf("failed to get user data: %v", err)
	}

	// update user balance
	newBalance := user.Balance + transaction.Amount
	err = uc.userRepo.UpdateBalance(user.ID, newBalance)
	if err != nil {
		return fmt.Errorf("failed to update user balance: %v", err)
	}

	// insert transaction
	err = uc.transactionRepo.CreateDepositCard(transaction)
	if err != nil {
		return fmt.Errorf("failed to create deposit transaction: %v", err)
	}

	return nil
}

func (uc *transactionUseCase) CreateWithdrawal(transaction *model.TransactionBank) error {
	user, err := uc.userRepo.GetByiD(transaction.SenderID)
	if err != nil {
		return fmt.Errorf("failed to get user data: %v", err)
	}

	// check user balance
	if user.Balance < transaction.Amount {
		return errors.New("insufficient balance")
	}

	// update user balance
	newBalance := user.Balance - transaction.Amount
	err = uc.userRepo.UpdateBalance(user.ID, newBalance)
	if err != nil {
		return fmt.Errorf("failed to update user balance: %v", err)
	}

	// insert transaction
	err = uc.transactionRepo.CreateWithdrawal(transaction)
	if err != nil {
		return fmt.Errorf("failed to create withdrawal transaction: %v", err)
	}

	return nil
}
func (uc *transactionUseCase) CreateTransfer(sender *model.User, recipient *model.User, amount uint) (any, error) {
	// update sender balance
	newBalanceS := sender.Balance - amount
	err := uc.userRepo.UpdateBalance(sender.ID, newBalanceS)
	if err != nil {
		return newBalanceS, err
	}

	// update recipient balance
	newBalanceR := recipient.Balance + amount
	err = uc.userRepo.UpdateBalance(recipient.ID, newBalanceR)
	if err != nil {
		return newBalanceR, err
	}

	// insert transaction
	newTransfer := model.TransactionTransfer{
		SenderID:    sender.ID,
		RecipientID: recipient.ID,
		Amount:      amount,
	}
	return uc.transactionRepo.CreateTransfer(&newTransfer)
}

func (uc *transactionUseCase) CreateRedeem(transaction *model.TransactionPoint) error {
	user, err := uc.userRepo.GetByiD(transaction.SenderID)
	if err != nil {
		return err
	}

	// update user balance
	newBalance := user.Balance + transaction.Point
	err = uc.userRepo.UpdateBalance(user.ID, newBalance)
	if err != nil {
		return err
	}

	// update user point
	newPoint := user.Point - transaction.Point
	err = uc.userRepo.UpdatePoint(user.ID, newPoint)
	if err != nil {
		return err
	}

	// insert transaction
	err = uc.transactionRepo.CreateRedeem(transaction)
	if err != nil {
		return err
	}

	return nil
}
