package usecase

import (
	"fmt"
	"testing"
	"time"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type transactionRepoMock struct {
	mock.Mock
}
type userRepoMock struct {
	mock.Mock
	users []model.Credentials
}

var dummyUsers = []*model.User{
	{
		ID:           uint(1),
		Name:         "fikri",
		Username:     "fikri",
		Email:        "fikri@",
		Password:     "fikri",
		Phone_Number: "0838",
		Address:      "smi",
		Balance:      1000,
		Point:        20,
		Role:         "user",
	},
}

type TransactionUseCaseTestSuite struct {
	transactionRepoMock *transactionRepoMock
	userRepoMock        *userRepoMock

	suite.Suite
}

func (suite *TransactionUseCaseTestSuite) SetupTest() {
	suite.transactionRepoMock = new(transactionRepoMock)
}

func TestTransactionUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionUseCaseTestSuite))
}

func (m *transactionRepoMock) GetBySenderId(senderId uint) ([]*model.Transaction, error) {
	args := m.Called(senderId)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*model.Transaction), nil
}
func (m *transactionRepoMock) CreateDepositBank(tx *model.TransactionBank) error {
	args := m.Called(tx)

	return args.Error(0)
}
func (m *transactionRepoMock) CreateDepositCard(tx *model.TransactionCard) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *transactionRepoMock) CreateWithdrawal(tx *model.TransactionWithdraw) error {
	args := m.Called(tx)

	return args.Error(0)
}
func (m *transactionRepoMock) CreateTransfer(tx *model.TransactionTransfer) (any, error) {
	args := m.Called(tx)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0), nil
}
func (m *transactionRepoMock) CreateRedeem(tx *model.TransactionPoint) error {
	args := m.Called(tx)

	return args.Error(0)
}
func (m *transactionRepoMock) GetAllPoint() ([]*model.PointExchange, error) {
	args := m.Called()
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*model.PointExchange), nil
}
func (m *transactionRepoMock) GetByPeId(id uint) ([]*model.PointExchange, error) {
	args := m.Called(id)

	// check if the first argument is nil or not
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*model.PointExchange), args.Error(1)
}

func (m *userRepoMock) UpdateBalance(userID uint, newBalance uint) error {
	args := m.Called(userID, newBalance)
	return args.Error(0)
}
func (m *userRepoMock) Create(user *model.UserCreate) (any, error) {
	args := m.Called(user)
	result := args.Get(0)
	err := args.Error(1)
	if result == nil {
		return nil, err
	}
	return result.(*model.UserCreate), err
}
func (m *userRepoMock) UpdatePoint(userID uint, newPoint int) error {
	args := m.Called(userID, newPoint)
	return args.Error(0)
}
func (m *userRepoMock) GetAll() any {
	args := m.Called()
	return args.Get(0)
}
func (m *userRepoMock) GetByUsername(username string) (*model.UserResponse, error) {
	args := m.Called(username)
	result := args.Get(0)
	err := args.Error(1)
	if result == nil {
		return nil, err
	}
	return result.(*model.UserResponse), err
}

func (m *userRepoMock) GetByiD(id uint) (*model.User, error) {
	args := m.Called(id)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*model.User), args.Error(1)
}

func (r *userRepoMock) UpdateProfile(user *model.User) string {
	if user.ID == 0 {
		return "user not found"
	}

	return "updated profile successfully"
}
func (m *userRepoMock) UpdateEmailPassword(user *model.User) string {
	args := m.Called(user)
	return args.String(0)
}
func (m *userRepoMock) Delete(user *model.User) string {
	if user.Username == "not-found" {
		return "user not found"
	}

	return "deleted user successfully"
}
func (r *userRepoMock) GetByEmailAndPassword(email string, password string) (*model.Credentials, error) {
	// Check if email and password match with any user
	for _, user := range r.users {
		if user.Email == email && utils.CheckPasswordHash(password, user.Password) == nil {
			return &model.Credentials{
				Username: user.Username,
				UserID:   user.UserID,
				Role:     user.Role,
			}, nil
		}
	}

	return nil, fmt.Errorf("invalid credentials")
}

var senderID = uint(1)

func (suite *TransactionUseCaseTestSuite) TestFindTxById() {
	// set up expectations

	expectedTxs := []*model.Transaction{
		&model.Transaction{SenderID: senderID},
		&model.Transaction{SenderID: senderID},
	}
	suite.transactionRepoMock.On("GetBySenderId", senderID).Return(expectedTxs, nil)

	// call the method being tested
	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)
	actualTxs, err := uc.FindTxById(senderID)

	// assert the expected results
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedTxs, actualTxs)

	// assert that the expected function calls were made
	suite.transactionRepoMock.AssertExpectations(suite.T())
}
func (suite *TransactionUseCaseTestSuite) TestFindByPeId() {
	// set up expectations
	expectedPEs := []*model.PointExchange{
		&model.PointExchange{PE_ID: 1, Reward: "10K Pulsa", Price: 100},
		&model.PointExchange{PE_ID: 2, Reward: "20k pulsa", Price: 200},
	}
	suite.transactionRepoMock.On("GetByPeId", uint(1)).Return(expectedPEs, nil)

	// call the method being tested
	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)
	actualPEs, err := uc.FindByPeId(uint(1))

	// assert the expected results
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedPEs, actualPEs)

	// assert that the expected function calls were made
	suite.transactionRepoMock.AssertExpectations(suite.T())
}
func (suite *TransactionUseCaseTestSuite) TestCreateDepositBank() {
	// set up expectations
	tx_type := "Deposit Bank"
	tx_date := time.Now()

	amount := uint(50000)
	bank_account_id := uint(1)
	expectedUser := dummyUsers[0]
	suite.userRepoMock.On("GetByiD", dummyUsers[0].ID).Return(expectedUser, nil)

	suite.userRepoMock.On("UpdateBalance", dummyUsers[0].ID, expectedUser.Balance+(amount)).Return(nil)
	suite.userRepoMock.On("UpdatePoint", dummyUsers[0].ID, expectedUser.Point+20).Return(nil)
	expectedTransactionBank := &model.TransactionBank{TransactionType: tx_type, SenderID: senderID, Amount: amount, BankAccountID: bank_account_id, TransactionDate: tx_date}
	suite.transactionRepoMock.On("CreateDepositBank", expectedTransactionBank).Return(nil)
	// call the method being tested
	uc := NewTransactionUseCase(suite.transactionRepoMock, suite.userRepoMock)
	err := uc.CreateDepositBank(expectedTransactionBank)

	// add log for debugging
	fmt.Printf("expectedUser before update balance: %+v\n", expectedUser)
	fmt.Printf("expectedTransactionBank: %+v\n", expectedTransactionBank)
	fmt.Printf("err: %+v\n", err)

	// assert the expected results
	assert.NoError(suite.T(), err)

	// assert that the expected function calls were made
	suite.userRepoMock.AssertExpectations(suite.T())
	suite.transactionRepoMock.AssertExpectations(suite.T())
}
