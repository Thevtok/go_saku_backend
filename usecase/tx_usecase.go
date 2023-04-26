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
	CreateWithdrawal(transaction *model.TransactionWithdraw) error
	CreateTransfer(sender *model.User, recipient *model.User, amount uint) (any, error)
	CreateRedeem(transaction *model.TransactionPoint) error
	FindTxById(senderId uint) ([]*model.Transaction, error)
	FindByPeId(id uint) ([]*model.PointExchange, error)
}

type transactionUseCase struct {
	transactionRepo repository.TransactionRepository
	userRepo        repository.UserRepository
}

func (uc *transactionUseCase) FindByPeId(id uint) ([]*model.PointExchange, error) {
	return uc.transactionRepo.GetByPeId(id)
}

func (uc *transactionUseCase) FindTxById(senderId uint) ([]*model.Transaction, error) {
	return uc.transactionRepo.GetBySenderId(senderId)
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

	// check if user is eligible for bonus points
	newPoint := user.Point + 20 // change from 20 to 100
	err = uc.userRepo.UpdatePoint(user.ID, newPoint)
	if err != nil {

		return err
	}

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

	newPoint := user.Point + 20 // change from 20 to 100
	err = uc.userRepo.UpdatePoint(user.ID, newPoint)
	if err != nil {
		return err
	}

	// check if user is eligible for bonus points

	// insert transaction
	err = uc.transactionRepo.CreateDepositCard(transaction)
	if err != nil {
		return fmt.Errorf("failed to create deposit transaction: %v", err)
	}

	return nil
}

func (uc *transactionUseCase) CreateWithdrawal(transaction *model.TransactionWithdraw) error {
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
	// validate sender balance
	if sender.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	// update recipient balance
	newBalanceR := recipient.Balance + amount
	err = uc.userRepo.UpdateBalance(recipient.ID, newBalanceR)
	if err != nil {
		return newBalanceR, err
	}

	// check if sender is eligible for bonus points
	newPoint := sender.Point + 20 // change from 20 to 100
	err = uc.userRepo.UpdatePoint(sender.ID, newPoint)
	if err != nil {
		return nil, err
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

	// Get all point exchanges
	pointExchanges, err := uc.transactionRepo.GetAllPoint()
	if err != nil {
		return err
	}

	// Find the point exchange with matching pe_id
	var pointExchange *model.PointExchange
	for _, pe := range pointExchanges {
		if pe.PE_ID == transaction.PointExchangeID {
			pointExchange = pe
			break
		}
	}

	// Check if point exchange was found
	if pointExchange == nil {
		return fmt.Errorf("point exchange with pe_id %d not found", transaction.PointExchangeID)
	}

	// Check if point exchange reward and price match with transaction data
	if pointExchange.Price != transaction.Point {
		return fmt.Errorf("reward or price on point exchange data doesn't match with the transaction data")
	}

	// update user balance
	if user.Point < transaction.Point {
		return fmt.Errorf("your point is not enough to redeem")
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

func NewTransactionUseCase(transactionRepo repository.TransactionRepository, userRepo repository.UserRepository) TransactionUseCase {
	return &transactionUseCase{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
	}
}
