package controller

import (
	"errors"
	"testing"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/gin-gonic/gin"
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

var dummyBankAccResponse1 = []*model.BankAccResponse{
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

type BankAccUsecaseMock struct {
	mock.Mock
}

func (r *BankAccUsecaseMock) FindAllBankAcc() any {
	args := r.Called()
	if args.Get(0) == nil {
		return nil
	}
	return dummyBankAcc
}

func (r *BankAccUsecaseMock) FindBankAccByUserID(id uint) ([]*model.BankAccResponse, error) {
	args := r.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.BankAccResponse), nil
}

func (r *BankAccUsecaseMock) FindBankAccByAccountID(id uint) (*model.BankAcc, error) {
	args := r.Called(id)
	if args.Get(0) == nil {
		return nil, errors.New("bank account not found")
	}
	return args.Get(0).(*model.BankAcc), nil
}

func (r *BankAccUsecaseMock) Register(id uint, newBankAcc *model.BankAccResponse) (any, error) {
	args := r.Called(id, newBankAcc)
	if args.Get(0) == nil {
		return nil, errors.New("failed to create data")
	}
	return dummyBankAccResponse, nil
}

func (r *BankAccUsecaseMock) Edit(bankAcc *model.BankAcc) string {
	args := r.Called(bankAcc)
	if args.Get(0) == nil {
		return "failed to update Bank Account"
	}
	return "Bank Account updated Successfully"
}

func (r *BankAccUsecaseMock) UnregAll(id uint) string {
	args := r.Called(id)
	if args.Get(0) == nil {
		return "failed to delete Bank Account"
	}
	return "All Bank Account deleted Successfully"
}

func (r *BankAccUsecaseMock) UnregByAccountID(accountID uint) error {
	args := r.Called(accountID)
	if args.Get(0) != nil {
		return errors.New("failed to delete data")
	}
	return nil
}

type BankAccControllerTestSuite struct {
	suite.Suite
	routerMock  *gin.Engine
	usecaseMock *BankAccUsecaseMock
}

func TestBankAccController(t *testing.T) {
	suite.Run(t, new(BankAccControllerTestSuite))
}
