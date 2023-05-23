package usecase

import (
	"errors"
	"testing"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyBankAcc = []model.BankAcc{
	{
		AccountID:         1,
		UserID:            "1",
		BankName:          "Test1",
		AccountNumber:     "123412341111",
		AccountHolderName: "Test1",
	},
	{
		AccountID:         2,
		UserID:            "1",
		BankName:          "Test2",
		AccountNumber:     "123412341112",
		AccountHolderName: "Test2",
	},
	{
		AccountID:         3,
		UserID:            "2",
		BankName:          "Test3",
		AccountNumber:     "123412341113",
		AccountHolderName: "Test3",
	},
	{
		AccountID:         4,
		UserID:            "2",
		BankName:          "Test2",
		AccountNumber:     "123412341114",
		AccountHolderName: "Test4",
	},
}

var dummyBankAccResponse1 = []*model.BankAcc{
	{
		UserID:            "1",
		AccountID:         1,
		BankName:          "Test1",
		AccountNumber:     "123412341111",
		AccountHolderName: "Test1",
	},
	{
		UserID:            "1",
		AccountID:         2,
		BankName:          "Test2",
		AccountNumber:     "123412341112",
		AccountHolderName: "Test2",
	},
	{
		UserID:            "2",
		AccountID:         3,
		BankName:          "Test3",
		AccountNumber:     "123412341113",
		AccountHolderName: "Test3",
	},
	{
		UserID:            "2",
		AccountID:         4,
		BankName:          "Test4",
		AccountNumber:     "123412341114",
		AccountHolderName: "Test4",
	},
}

type bankaccRepoMock struct {
	mock.Mock
}

func (r *bankaccRepoMock) GetAll() any {
	args := r.Called()
	if args.Get(0) == nil {
		return nil
	}
	return dummyBankAcc
}

func (r *bankaccRepoMock) GetByUserID(id string) ([]*model.BankAcc, error) {
	args := r.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.BankAcc), nil
}

func (r *bankaccRepoMock) GetByAccountID(id uint) (*model.BankAcc, error) {
	args := r.Called(id)
	if args.Get(0) == nil {
		return nil, errors.New("bank account not found")
	}
	return args.Get(0).(*model.BankAcc), nil
}

func (r *bankaccRepoMock) Create(id string, newBankAcc *model.BankAcc) (any, error) {
	args := r.Called(id, newBankAcc)
	if args.Get(0) == nil {
		return nil, errors.New("failed to create data")
	}
	return dummyBankAcc, nil
}

func (r *bankaccRepoMock) Update(bankAcc *model.BankAcc) string {
	args := r.Called(bankAcc)
	if args.Get(0) == nil {
		return "failed to update Bank Account"
	}
	return "Bank Account updated Successfully"
}

func (r *bankaccRepoMock) DeleteByAccountID(accountID uint) error {
	args := r.Called(accountID)
	if args.Get(0) == nil {
		return nil
	}
	return errors.New("failed to delete data")
}

type BankAccUsecaseTestSuite struct {
	bankaccRepoMock *bankaccRepoMock
	suite.Suite
}

func (suite *BankAccUsecaseTestSuite) TestFindAccByUserID_Success() {
	userID := "uint(1)"
	bankAcc := dummyBankAccResponse1
	bankAccUsecase := NewBankAccUsecase(suite.bankaccRepoMock)
	suite.bankaccRepoMock.On("GetByUserID", userID).Return(bankAcc, nil)
	result, err := bankAccUsecase.FindBankAccByUserID(userID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), bankAcc, result)
}

func (suite *BankAccUsecaseTestSuite) TestFindAccByUserID_Failed() {
	userID := "uint(1)"
	expectedErr := errors.New("failed to get bank account")
	bankAccUsecase := NewBankAccUsecase(suite.bankaccRepoMock)
	suite.bankaccRepoMock.On("GetByUserID", userID).Return(nil, expectedErr)
	result, err := bankAccUsecase.FindBankAccByUserID(userID)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *BankAccUsecaseTestSuite) TestFindAccByAccID_Success() {
	accID := uint(1)
	bankAcc := &dummyBankAcc[0]
	bankAccUsecase := NewBankAccUsecase(suite.bankaccRepoMock)
	suite.bankaccRepoMock.On("GetByAccountID", accID).Return(bankAcc, nil)
	result, err := bankAccUsecase.FindBankAccByAccountID(accID)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), bankAcc, result)
}

func (suite *BankAccUsecaseTestSuite) TestFindAccByAccID_Failed() {
	accID := uint(1)
	expectedErr := errors.New("bank account not found")
	bankAccUsecase := NewBankAccUsecase(suite.bankaccRepoMock)
	suite.bankaccRepoMock.On("GetByAccountID", accID).Return(nil, expectedErr)
	result, err := bankAccUsecase.FindBankAccByAccountID(accID)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *BankAccUsecaseTestSuite) TestRegister_Success() {
	userID := "uint(1)"
	bankAccUsecase := NewBankAccUsecase(suite.bankaccRepoMock)
	suite.bankaccRepoMock.On("Create", userID, &dummyBankAcc[0]).Return(dummyBankAcc, nil)
	result, err := bankAccUsecase.Register(userID, &dummyBankAcc[0])
	assert.NotNil(suite.T(), result)
	assert.Nil(suite.T(), err)
}

func (suite *BankAccUsecaseTestSuite) TestRegister_Failed() {
	userID := " uint(1)"
	expectedErr := errors.New("failed to create data")
	bankAccUsecase := NewBankAccUsecase(suite.bankaccRepoMock)
	suite.bankaccRepoMock.On("Create", userID, &dummyBankAcc[0]).Return(nil, expectedErr)
	result, err := bankAccUsecase.Register(userID, &dummyBankAcc[0])
	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), err)
}

func (suite *BankAccUsecaseTestSuite) TestUnregByAccountID_Success() {
	accountID := uint(1)
	bankAccUsecase := NewBankAccUsecase(suite.bankaccRepoMock)
	suite.bankaccRepoMock.On("DeleteByAccountID", accountID).Return(nil)
	err := bankAccUsecase.UnregByAccountID(accountID)
	assert.NoError(suite.T(), err)
}

func (suite *BankAccUsecaseTestSuite) TestUnregByAccountID_Failed() {
	accountID := uint(1)
	bankAccUsecase := NewBankAccUsecase(suite.bankaccRepoMock)
	suite.bankaccRepoMock.On("DeleteByAccountID", accountID).Return(errors.New("failed to delete data"))
	err := bankAccUsecase.UnregByAccountID(accountID)
	assert.EqualError(suite.T(), err, "failed to delete data")
}

func (suite *BankAccUsecaseTestSuite) SetupTest() {
	suite.bankaccRepoMock = new(bankaccRepoMock)
}

func TestBankAccTestSuite(t *testing.T) {
	suite.Run(t, new(BankAccUsecaseTestSuite))
}
