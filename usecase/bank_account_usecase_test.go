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

var dummyBankAccResponse = []model.BankAccResponse{
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
	return args.Get(0).([]*model.BankAccResponse), args.Error(1)
}

func (r *bankaccRepoMock) GetByAccountID(id uint) (*model.BankAcc, error) {
	args := r.Called(id)
	if args.Get(0) == nil {
		return nil, errors.New("bank account not found")
	}
	return args.Get(0).(*model.BankAcc), args.Error(1)
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
		return nil
	}
	return nil
}

type BankAccUsecaseTestSuite struct {
	bankaccRepoMock *bankaccRepoMock
	suite.Suite
}

func (suite *BankAccUsecaseTestSuite) TestFindAllBankAcc_Success() {
	bankAcc := NewBankAccUsecase(suite.bankaccRepoMock)
	suite.bankaccRepoMock.On("GetAll").Return(dummyBankAcc)
	res := bankAcc.FindAllBankAcc()
	assert.NotNil(suite.T(), res)
	assert.Equal(suite.T(), dummyBankAcc, res)
}

func (suite *BankAccUsecaseTestSuite) TestFindAllBankAcc_Failed() {
	bankAcc := NewBankAccUsecase(suite.bankaccRepoMock)
	suite.bankaccRepoMock.On("GetAll").Return(nil)
	res := bankAcc.FindAllBankAcc()
	assert.Nil(suite.T(), res)
	assert.Empty(suite.T(), dummyBankAcc, res)
}

// func (suite *BankAccUsecaseTestSuite) TestFindByUserID_Success() {
// 	bankAcc := NewBankAccUsecase(suite.bankaccRepoMock)

// 	mockResult := []*model.BankAccResponse{
// 		{
// 			UserID:            1,
// 			BankName:          "Bank A",
// 			AccountNumber:     "1234567890",
// 			AccountHolderName: "John Doe",
// 		},
// 		{
// 			UserID:            1,
// 			BankName:          "Bank B",
// 			AccountNumber:     "0987654321",
// 			AccountHolderName: "Jane Doe",
// 		},
// 	}

// 	suite.bankaccRepoMock.On("GetByUserID", uint(1)).Return(mockResult, nil)
// 	result, err := bankAcc.FindByUserID(1)

// 	assert.NoError(suite.T(), err)
// 	assert.NotNil(suite.T(), result)
// 	assert.Equal(suite.T(), mockResult, result)
// }

func (suite *BankAccUsecaseTestSuite) SetupTest() {
	suite.bankaccRepoMock = new(bankaccRepoMock)
}

func TestBankAccTestSuite(t *testing.T) {
	suite.Run(t, new(BankAccUsecaseTestSuite))
}
