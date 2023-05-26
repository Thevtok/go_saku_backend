package usecase

import (
	"fmt"
	"time"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/repository"
)

var now = time.Now().Local()
var date = now.Format("2006-01-02")

type TransactionUseCase interface {
	CreateDepositBank(transaction *model.Deposit) error

	CreateWithdrawal(transaction *model.Withdraw) error
	CreateTransfer(sender *model.User, recipient *model.User, amount int) error
	CreateRedeem(transaction *model.Redeem) error
	FindTxById(userID string) ([]*model.Transaction, error)
	FindByPeId(id int) (*model.PointExchange, error)
	AssignBadge(user *model.User) error
	UpdateDepositStatus(vaNumber, token string) error
}

type transactionUseCase struct {
	transactionRepo repository.TransactionRepository
	userRepo        repository.UserRepository
}

func (uc *transactionUseCase) UpdateDepositStatus(vaNumber, token string) error {
	err := uc.transactionRepo.UpdateDepositStatus(vaNumber, token)
	if err != nil {
		return fmt.Errorf("failed to update deposit status: %v", err)
	}

	return nil
}

func (uc *transactionUseCase) AssignBadge(user *model.User) error {
	err := uc.transactionRepo.AssignBadge(user)
	if err != nil {
		return fmt.Errorf("failed to assign badge: %v", err)
	}
	return nil
}

func (uc *transactionUseCase) FindByPeId(id int) (*model.PointExchange, error) {
	return uc.transactionRepo.GetByPeId(id)
}

func (uc *transactionUseCase) FindTxById(userID string) ([]*model.Transaction, error) {
	return uc.transactionRepo.GetTransactions(userID)
}

func (uc *transactionUseCase) CreateDepositBank(transaction *model.Deposit) error {
	user, err := uc.userRepo.GetByiD(transaction.UserID)
	if err != nil {
		return fmt.Errorf("gagal mendapatkan data pengguna: %v", err)
	}

	// cek apakah pengguna memenuhi syarat untuk bonus poin
	if transaction.Amount >= 50000 {
		newPoint := user.Point + 20
		_ = uc.userRepo.UpdatePoint(user.ID, newPoint)
	} else if transaction.Amount < 50000 {
		newPoint := user.Point
		_ = uc.userRepo.UpdatePoint(user.ID, newPoint)
	}

	err = uc.transactionRepo.CreateDepositBank(transaction)
	if err != nil {
		return fmt.Errorf("gagal membuat transaksi deposit: %v", err)
	}

	return nil
}

func (uc *transactionUseCase) CreateWithdrawal(transaction *model.Withdraw) error {
	user, err := uc.userRepo.GetByiD(transaction.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user data: %v", err)
	}

	// check user balance
	if user.Balance < transaction.Amount {
		return fmt.Errorf("insufficient balance")
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
func (uc *transactionUseCase) CreateTransfer(sender *model.User, recipient *model.User, amount int) error {
	// Update sender balance
	newBalanceS := sender.Balance - amount - 2500
	err := uc.userRepo.UpdateBalance(sender.ID, newBalanceS)
	if err != nil {
		return err
	}

	// Validate sender balance
	if sender.Balance < amount+2500 {
		return fmt.Errorf("insufficient balance")
	}

	// Update recipient balance
	newBalanceR := recipient.Balance + amount
	err = uc.userRepo.UpdateBalance(recipient.ID, newBalanceR)
	if err != nil {
		return err
	}

	// Update sender's point based on transfer amount
	var newPoint int
	if amount >= 50000 {
		newPoint = sender.Point + 20
	} else {
		newPoint = sender.Point
	}
	err = uc.userRepo.UpdatePoint(sender.ID, newPoint)
	if err != nil {
		return err
	}

	// Insert transaction
	newTransfer := model.Transfer{
		SenderID:             sender.ID,
		RecipientID:          recipient.ID,
		SenderPhoneNumber:    sender.Phone_Number,
		RecipientPhoneNumber: recipient.Phone_Number,
		Amount:               amount,
		TransactionType:      "Transfer",
		TransactionDate:      date,
		SenderName:           sender.Username,
		RecipientName:        recipient.Username,
	}
	return uc.transactionRepo.CreateTransfer(&newTransfer)
}

func (uc *transactionUseCase) CreateRedeem(transaction *model.Redeem) error {
	user, err := uc.userRepo.GetByiD(transaction.UserID)
	if err != nil {
		return err
	}
	// Get point exchange by ID
	pointExchange, err := uc.transactionRepo.GetByPeId(transaction.PEID)
	if err != nil {
		return err
	}

	// Check if point exchange reward and price match with transaction data
	if pointExchange.Price != transaction.Amount {
		return fmt.Errorf("reward or price on point exchange data doesn't match with the transaction data")
	}

	// update user balance
	if user.Point < transaction.Amount {
		return fmt.Errorf("your point is not enough to redeem")
	}

	// update user point
	newPoint := user.Point - transaction.Amount
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
