package usecase

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ReygaFitra/inc-final-project.git/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type transactionRepoMock struct {
	mock.Mock
}

type TransactionUseCaseTestSuite struct {
	transactionRepoMock *transactionRepoMock
	userRepoMock        *userRepoMock

	suite.Suite
}

var dummyUsers = []*model.User{
	{
		ID:           "",
		Name:         "name1",
		Username:     "username1",
		Email:        "email1@mail.com",
		Password:     "password1",
		Phone_Number: "08111111",
		Address:      "address1",
		Balance:      1,
		Role:         "user",
		Point:        10,
	},
}

var dummyTX = []*model.Transaction{}

var dummyTxPointExchange = []*model.PointExchange{
	{
		PE_ID:  1,
		Reward: "10K Pulsa",
		Price:  100,
	},
	{
		PE_ID:  2,
		Reward: "20k Pulsa",
		Price:  20,
	},
}

var dummyTxBank = []*model.Deposit{
	{

		TransactionType: "Deposit Bank",

		TransactionDate: "time.Now()",
	},
}

var dummyTxWithdraw = []*model.Withdraw{
	{

		TransactionType: "Deposit Bank",

		TransactionDate: "time.Now()",
	},
}

func (suite *TransactionUseCaseTestSuite) SetupTest() {
	suite.transactionRepoMock = new(transactionRepoMock)
	suite.userRepoMock = new(userRepoMock)
}

func TestTransactionUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionUseCaseTestSuite))
}
func (m *transactionRepoMock) AssignBadge(user *model.User) error {
	args := m.Called(user)

	if args[0] != nil {
		return args.Error(0)
	}
	return nil
}
func (m *transactionRepoMock) UpdateDepositStatus(vaNumber, token string) error {
	args := m.Called(vaNumber, token)

	if args[0] != nil {
		return args.Error(0)
	}
	return nil
}

func (m *transactionRepoMock) GetTransactions(ID string) ([]*model.Transaction, error) {
	args := m.Called(ID)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*model.Transaction), nil
}
func (m *transactionRepoMock) CreateDepositBank(tx *model.Deposit) error {
	args := m.Called(tx)
	if args[0] != nil {
		return args.Error(0)
	}
	return nil

}

func (m *transactionRepoMock) CreateWithdrawal(tx *model.Withdraw) error {
	args := m.Called(tx)

	if args[0] != nil {
		return args.Error(0)
	}
	return nil
}

func (m *transactionRepoMock) CreateTransfer(tx *model.Transfer) error {
	args := m.Called(tx)
	if args[0] != nil {
		return args.Error(0)
	}

	return nil
}

func (m *transactionRepoMock) CreateRedeem(tx *model.Redeem) error {
	args := m.Called(tx)

	return args.Error(0)
}
func (m *transactionRepoMock) GetAllPoint() ([]*model.PointExchange, error) {
	args := m.Called()
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*model.PointExchange), args.Error(1)

}
func (m *transactionRepoMock) GetByPeId(id int) (*model.PointExchange, error) {
	args := m.Called(id)

	// check if the first argument is nil or not
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.PointExchange), args.Error(1)
}

var senderID = "uint(1)"

func (suite *TransactionUseCaseTestSuite) TestFindTxById_Success() {
	// set up expectations

	expectedTxs := dummyTX
	suite.transactionRepoMock.On("GetBySenderId", senderID).Return(expectedTxs, nil)

	// call the method being tested
	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)
	actualTxs, err := uc.FindTxById(senderID)

	// assert the expected results
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedTxs, actualTxs)

	// assert that the expected function calls were made

}

func (suite *TransactionUseCaseTestSuite) TestFindByPeId_Success() {
	// set up expectations
	expectedPEs := dummyTxPointExchange[0]
	suite.transactionRepoMock.On("GetByPeId", 1).Return(expectedPEs, nil)

	// call the method being tested
	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)
	actualPEs, err := uc.FindByPeId(1)

	// assert the expected results
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedPEs, actualPEs)

	// assert that the expected function calls were made

}

func (suite *TransactionUseCaseTestSuite) TestCreateDepositBank() {
	user := dummyUsers[0]
	bank := dummyTxBank[0]

	suite.userRepoMock.On("GetByiD", "").Return(user, nil)

	newBalance := user.Balance + bank.Amount

	suite.userRepoMock.On("UpdateBalance", user.ID, newBalance).Return(nil)

	newPoint := user.Point + 20

	suite.userRepoMock.On("UpdatePoint", user.ID, newPoint).Return(nil)

	suite.transactionRepoMock.On("CreateDepositBank", bank).Return(nil)

	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)
	err := uc.CreateDepositBank(bank)

	// assert the expected results
	assert.NoError(suite.T(), err)
}
func (suite *TransactionUseCaseTestSuite) TestCreateDepositBank_UserNotFound() {
	transaction := &model.Deposit{

		Amount: 10000,
	}
	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)
	suite.userRepoMock.On("GetByiD", "transaction.SenderID").Return(nil, errors.New("user not found"))

	err := uc.CreateDepositBank(transaction)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "failed to get user data: user not found", err.Error())
}

func (suite *TransactionUseCaseTestSuite) TestCreateDepositBank_UpdateBalanceError() {
	transaction := &model.Deposit{

		Amount: 10000,
	}
	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)
	suite.userRepoMock.On("GetByiD", "transaction.SenderID").Return(&model.User{}, nil)
	suite.userRepoMock.On("UpdateBalance", mock.Anything, mock.Anything).Return(errors.New("balance update error"))

	err := uc.CreateDepositBank(transaction)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "failed to update user balance: balance update error", err.Error())
}

func (suite *TransactionUseCaseTestSuite) TestCreateDepositBank_Failed() {
	// Create a dummy transaction
	transaction := &model.Deposit{

		Amount: 100,
	}

	// Set up the mock repository to return an error
	expectedErr := errors.New("failed to create deposit transaction")
	suite.userRepoMock.On("GetByiD", "transaction.SenderID").Return(&model.User{}, nil)
	suite.userRepoMock.On("UpdateBalance", mock.Anything, mock.Anything).Return(nil)
	suite.userRepoMock.On("UpdatePoint", mock.Anything, mock.Anything).Return(nil)
	suite.transactionRepoMock.On("CreateDepositBank", transaction).Return(expectedErr)

	// Create the use case and call the function being tested
	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)
	err := uc.CreateDepositBank(transaction)

	// Verify that the function returns an error
	assert.EqualError(suite.T(), err, fmt.Sprintf("failed to create deposit transaction: %v", expectedErr))
}

func (suite *TransactionUseCaseTestSuite) TestCreateWithdrawal() {
	user := dummyUsers[0]
	withdraw := dummyTxWithdraw[0]

	suite.userRepoMock.On("GetByiD", "withdraw.SenderID").Return(user, nil)

	newBalance := user.Balance - withdraw.Amount

	suite.userRepoMock.On("UpdateBalance", user.ID, newBalance).Return(nil)

	suite.transactionRepoMock.On("CreateWithdrawal", withdraw).Return(nil)

	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)
	err := uc.CreateWithdrawal(withdraw)

	// assert the expected results
	assert.NoError(suite.T(), err)
}
func (suite *TransactionUseCaseTestSuite) TestCreateWithdrawal_UserNotFound() {
	transaction := &model.Withdraw{

		Amount: 10000,
	}
	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)
	suite.userRepoMock.On("GetByiD", "transaction.SenderID").Return(nil, errors.New("user not found"))

	err := uc.CreateWithdrawal(transaction)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "failed to get user data: user not found", err.Error())
}
func (suite *TransactionUseCaseTestSuite) TestCreateWithdrawal_InsufficientBalance() {
	transaction := &model.Withdraw{

		Amount: 10000,
	}
	user := &model.User{
		ID:      "",
		Name:    "Test User",
		Email:   "test@example.com",
		Balance: 5000,
	}
	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)
	suite.userRepoMock.On("GetByiD", "transaction.SenderID").Return(user, nil)
	suite.userRepoMock.On("UpdateBalance", user.ID, user.Balance-transaction.Amount).Return(nil)

	err := uc.CreateWithdrawal(transaction)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "insufficient balance", err.Error())
}
func (suite *TransactionUseCaseTestSuite) TestCreateWithdrawal_UpdateBalance() {
	transaction := &model.Withdraw{

		Amount: 10000,
	}
	user := &model.User{
		ID:      "1",
		Name:    "Test User",
		Email:   "test@example.com",
		Balance: 15000,
	}
	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)
	suite.userRepoMock.On("GetByiD", "transaction.SenderID").Return(user, nil)
	suite.userRepoMock.On("UpdateBalance", user.ID, user.Balance-transaction.Amount).Return(errors.New("failed to update user balance"))

	err := uc.CreateWithdrawal(transaction)

	assert.NotNil(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to update user balance")
}
func (suite *TransactionUseCaseTestSuite) TestCreateWithdrawal_CreateTransactionError() {
	transaction := &model.Withdraw{

		Amount: 10000,
	}
	user := &model.User{
		ID:      "1",
		Name:    "Test User",
		Email:   "test@example.com",
		Balance: 15000,
	}
	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)
	suite.userRepoMock.On("GetByiD", "transaction.SenderID").Return(user, nil)
	suite.userRepoMock.On("UpdateBalance", user.ID, user.Balance-transaction.Amount).Return(nil)
	suite.transactionRepoMock.On("CreateWithdrawal", transaction).Return(errors.New("failed to create withdrawal transaction"))

	err := uc.CreateWithdrawal(transaction)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "failed to create withdrawal transaction: failed to create withdrawal transaction", err.Error())
}

// func (suite *TransactionUseCaseTestSuite) TestCreateRedeem_Success() {
// 	// set up test data
// 	sender := &model.User{
// 		ID:      1,
// 		Name:    "Sender",
// 		Balance: 0,
// 		Point:   100,
// 	}
// 	pointExchange := &model.PointExchange{
// 		PE_ID: 1,

// 		Reward: "10K Pulsa",
// 		Price:  100,
// 	}
// 	transaction := &model.TransactionPoint{
// 		SenderID:        sender.ID,
// 		PointExchangeID: pointExchange.PE_ID,
// 		Point:           pointExchange.Price,
// 	}

// 	// set up mock repository behavior
// 	suite.userRepoMock.On("GetByiD", transaction.SenderID).
// 		Return(sender, nil)
// 	suite.transactionRepoMock.On("GetByPeId", transaction.PointExchangeID).
// 		Return(pointExchange, nil)
// 	suite.userRepoMock.On("UpdatePoint", sender.ID, sender.Point-transaction.Point).
// 		Return(nil)
// 	suite.transactionRepoMock.On("CreateRedeem", transaction).
// 		Return(nil)

// 	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)

// 	// call the use case
// 	err := uc.CreateRedeem(transaction)

// 	// check the result and error
// 	assert.Nil(suite.T(), err)
// 	assert.Equal(suite.T(), sender.Point-transaction.Point, 0)
// }

// func (suite *TransactionUseCaseTestSuite) TestCreateRedeem_UserRepoGetByIDError() {
// 	// set up test data
// 	transaction := &model.TransactionPoint{
// 		SenderID:        1,
// 		PointExchangeID: 1,
// 		Point:           10,
// 	}
// 	suite.userRepoMock.On("GetByiD", transaction.SenderID).
// 		Return(nil, errors.New("failed to get user by ID"))
// 	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)

// 	// call the use case
// 	err := uc.CreateRedeem(transaction)

// 	// check the result and error
// 	assert.NotNil(suite.T(), err)
// 	assert.Equal(suite.T(), "failed to get user by ID", err.Error())
// }
// func (suite *TransactionUseCaseTestSuite) TestCreateRedeem_PointExchangeNotFound() {
// 	// set up test data
// 	user := &model.User{
// 		ID:      1,
// 		Name:    "User",
// 		Balance: 0,
// 		Point:   100,
// 	}
// 	transaction := &model.TransactionPoint{
// 		SenderID:        user.ID,
// 		PointExchangeID: 999, // ID yang tidak ada
// 		Point:           10,
// 	}
// 	suite.userRepoMock.On("GetByiD", user.ID).Return(user, nil)
// 	suite.transactionRepoMock.On("GetByPeId", transaction.PointExchangeID).Return(nil, fmt.Errorf("point exchange with pe_id %d not found", transaction.PointExchangeID))
// 	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)
// 	// call the use case
// 	err := uc.CreateRedeem(transaction)

// 	// check the error
// 	assert.NotNil(suite.T(), err)
// 	assert.Equal(suite.T(), "point exchange with pe_id 999 not found", err.Error())
// }

// func (suite *TransactionUseCaseTestSuite) TestCreateRedeem_PointExchangePriceNotMatch() {
// 	// set up test data
// 	user := &model.User{
// 		ID:    1,
// 		Name:  "John Doe",
// 		Point: 100,
// 	}
// 	transaction := &model.TransactionPoint{
// 		SenderID:        user.ID,
// 		PointExchangeID: 1,
// 		Point:           50,
// 	}
// 	pointExchange := &model.PointExchange{
// 		PE_ID:  1,
// 		Reward: "Free Coffee",
// 		Price:  30,
// 	}

// 	// set up mock repository behavior
// 	suite.userRepoMock.On("GetByiD", transaction.SenderID).Return(user, nil)
// 	suite.transactionRepoMock.On("GetByPeId", transaction.PointExchangeID).Return(pointExchange, nil)

// 	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)

// 	// call the use case
// 	err := uc.CreateRedeem(transaction)

// 	// check the result and error
// 	assert.NotNil(suite.T(), err)
// 	assert.Equal(suite.T(), "reward or price on point exchange data doesn't match with the transaction data", err.Error())
// }

// func (suite *TransactionUseCaseTestSuite) TestCreateRedeem_InsufficientPoint() {
// 	// set up test data
// 	sender := &model.User{
// 		ID:      1,
// 		Name:    "Sender",
// 		Balance: 100,
// 		Point:   10,
// 	}
// 	pe := &model.PointExchange{
// 		PE_ID:  1,
// 		Reward: "baso",
// 		Price:  50,
// 	}
// 	transaction := &model.TransactionPoint{
// 		SenderID:        sender.ID,
// 		PointExchangeID: pe.PE_ID,
// 		Point:           50,
// 	}

// 	// set up mock repository behavior
// 	suite.userRepoMock.On("GetByiD", transaction.SenderID).Return(sender, nil)
// 	suite.transactionRepoMock.On("GetByPeId", transaction.PointExchangeID).Return(pe, nil)

// 	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)

// 	// call the use case
// 	err := uc.CreateRedeem(transaction)

// 	// check the result and error
// 	assert.NotNil(suite.T(), err)
// 	assert.Equal(suite.T(), "your point is not enough to redeem", err.Error())
// }

// func (suite *TransactionUseCaseTestSuite) TestCreateRedeem_UpdatePointError() {
// 	// set up test data
// 	transaction := &model.TransactionPoint{
// 		SenderID:        1,
// 		PointExchangeID: 2,
// 		Point:           30,

// 		TransactionType: "REDEEM",
// 	}
// 	user := &model.User{
// 		ID:    1,
// 		Name:  "User",
// 		Point: 30,
// 	}

// 	// set up mock repository behavior
// 	suite.userRepoMock.On("GetByiD", transaction.SenderID).Return(user, nil)
// 	suite.transactionRepoMock.On("GetByPeId", transaction.PointExchangeID).Return(&model.PointExchange{
// 		PE_ID:  2,
// 		Reward: "bakso",
// 		Price:  30,
// 	}, nil)
// 	suite.userRepoMock.On("UpdatePoint", user.ID, user.Point-transaction.Point).Return(errors.New("failed to update point"))

// 	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)

// 	// call the use case
// 	err := uc.CreateRedeem(transaction)

// 	// check the result and error
// 	assert.NotNil(suite.T(), err)
// 	assert.Equal(suite.T(), "failed to update point", err.Error())
// }

// func (suite *TransactionUseCaseTestSuite) TestCreateRedeem_Error() {
// 	// set up test data
// 	sender := &model.User{
// 		ID:      1,
// 		Name:    "Sender",
// 		Balance: 0,
// 		Point:   100,
// 	}
// 	pointExchange := &model.PointExchange{
// 		PE_ID: 1,

// 		Reward: "10K Pulsa",
// 		Price:  100,
// 	}
// 	transaction := &model.TransactionPoint{
// 		SenderID:        sender.ID,
// 		PointExchangeID: pointExchange.PE_ID,
// 		Point:           pointExchange.Price,
// 	}

// 	// set up mock repository behavior
// 	suite.userRepoMock.On("GetByiD", transaction.SenderID).
// 		Return(sender, nil)
// 	suite.transactionRepoMock.On("GetByPeId", transaction.PointExchangeID).
// 		Return(pointExchange, nil)
// 	suite.userRepoMock.On("UpdatePoint", sender.ID, sender.Point-transaction.Point).
// 		Return(nil)
// 	suite.transactionRepoMock.On("CreateRedeem", transaction).
// 		Return(errors.New("err"))

// 	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)

// 	// call the use case
// 	err := uc.CreateRedeem(transaction)

// 	// check the result and error
// 	assert.NotNil(suite.T(), err)
// 	assert.Equal(suite.T(), "err", err.Error())
// }
