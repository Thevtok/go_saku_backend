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
		UserID:            1,
		BankName:          "Test1",
		AccountNumber:     "123412341111",
		AccountHolderName: "Test1",
	},
	{
		AccountID:         2,
		UserID:            2,
		BankName:          "Test2",
		AccountNumber:     "123412342222",
		AccountHolderName: "Test2",
	},
}

var dummyBankAccResponse = []*model.BankAccResponse{
	{
		UserID:            1,
		BankName:          "Test1",
		AccountNumber:     "123412341111",
		AccountHolderName: "Test1",
	},
	{
		UserID:            2,
		BankName:          "Test2",
		AccountNumber:     "123412342222",
		AccountHolderName: "Test2",
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

func (r *bankaccRepoMock) GetByUserID(id uint) ([]*model.BankAccResponse, error) {
	args := r.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.BankAccResponse), nil
}

func (r *bankaccRepoMock) GetByAccountID(id uint) (*model.BankAcc, error) {
	args := r.Called(id)
	if args.Get(0) == nil {
		return nil, errors.New("bank account not found")
	}
	return args.Get(0).(*model.BankAcc), nil
}

func (r *bankaccRepoMock) Create(id uint, newBankAcc *model.BankAccResponse) (any, error) {
	args := r.Called(id, newBankAcc)
	if args.Get(0) == nil {
		return nil, errors.New("failed to create data")
	}
	return dummyBankAccResponse, nil
}

func (r *bankaccRepoMock) Update(bankAcc *model.BankAcc) string {
	args := r.Called(bankAcc)
	if args.Get(0) == nil {
		return "failed to update Bank Account"
	}
	return "Bank Account updated Successfully"
}

func (r *bankaccRepoMock) DeleteByUserID(id uint) string {
	args := r.Called(id)
	if args.Get(0) == nil {
		return "failed to delete Bank Account"
	}
	return "All Bank Account deleted Successfully"
}

func (r *bankaccRepoMock) DeleteByAccountID(accountID uint) error {
	args := r.Called(accountID)
	if args.Get(0) == nil {
		return errors.New("failed to delete data")
	}
	return nil
}

type BankAccUsecaseTestSuite struct {
	bankaccRepoMock *bankaccRepoMock
	suite.Suite
}

func (suite *BankAccUsecaseTestSuite) TestFindAllBankAcc_Success() {
	bankAccUsecase := NewBankAccUsecase(suite.bankaccRepoMock)
	suite.bankaccRepoMock.On("GetAll").Return(dummyBankAcc)
	res := bankAccUsecase.FindAllBankAcc()
	assert.NotNil(suite.T(), res)
	assert.Equal(suite.T(), dummyBankAcc, res)
}

func (suite *BankAccUsecaseTestSuite) TestFindAllBankAcc_Failed() {
	bankAccUsecase := NewBankAccUsecase(suite.bankaccRepoMock)
	suite.bankaccRepoMock.On("GetAll").Return(nil)
	res := bankAccUsecase.FindAllBankAcc()
	assert.Nil(suite.T(), res)
	assert.Empty(suite.T(), res)
}

func (suite *BankAccUsecaseTestSuite) TestFindByUserID_Success() {
	userID := uint(1)
	bankAcc := dummyBankAccResponse
	bankAccUsecase := NewBankAccUsecase(suite.bankaccRepoMock)
	suite.bankaccRepoMock.On("GetByUserID", userID).Return(bankAcc, nil)
	result, err := bankAccUsecase.FindBankAccByUserID(userID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), bankAcc, result)
}

func (suite *BankAccUsecaseTestSuite) TestFindByUserID_Failed() {
	userID := uint(1)
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

func (suite *BankAccUsecaseTestSuite) SetupTest() {
	suite.bankaccRepoMock = new(bankaccRepoMock)
}

func TestBankAccTestSuite(t *testing.T) {
	suite.Run(t, new(BankAccUsecaseTestSuite))
}
