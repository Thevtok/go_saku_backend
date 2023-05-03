package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/gin-gonic/gin"
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
		UserID:            1,
		BankName:          "Test2",
		AccountNumber:     "123412341112",
		AccountHolderName: "Test2",
	},
	{
		AccountID:         3,
		UserID:            2,
		BankName:          "Test3",
		AccountNumber:     "123412341113",
		AccountHolderName: "Test3",
	},
	{
		AccountID:         4,
		UserID:            2,
		BankName:          "Test2",
		AccountNumber:     "123412341114",
		AccountHolderName: "Test4",
	},
}

var dummyBankAccResponse = []model.BankAccResponse{
	{
		UserID:            1,
		AccountID:         1,
		BankName:          "Test1",
		AccountNumber:     "123412341111",
		AccountHolderName: "Test1",
	},
	{
		UserID:            1,
		AccountID:         2,
		BankName:          "Test2",
		AccountNumber:     "123412341112",
		AccountHolderName: "Test2",
	},
	{
		UserID:            2,
		AccountID:         3,
		BankName:          "Test3",
		AccountNumber:     "123412341113",
		AccountHolderName: "Test3",
	},
	{
		UserID:            2,
		AccountID:         4,
		BankName:          "Test4",
		AccountNumber:     "123412341114",
		AccountHolderName: "Test4",
	},
}

var dummyBankAccResponse1 = []*model.BankAccResponse{
	{
		UserID:            1,
		AccountID:         1,
		BankName:          "Test1",
		AccountNumber:     "123412341111",
		AccountHolderName: "Test1",
	},
	{
		UserID:            1,
		AccountID:         2,
		BankName:          "Test2",
		AccountNumber:     "123412341112",
		AccountHolderName: "Test2",
	},
	{
		UserID:            2,
		AccountID:         3,
		BankName:          "Test3",
		AccountNumber:     "123412341113",
		AccountHolderName: "Test3",
	},
	{
		UserID:            2,
		AccountID:         4,
		BankName:          "Test4",
		AccountNumber:     "123412341114",
		AccountHolderName: "Test4",
	},
}

type BankAccUsecaseMock struct {
	mock.Mock
}

func setupRouterBankAcc() *gin.Engine {
	r := gin.Default()
	return r
}

func (u *BankAccUsecaseMock) FindAllBankAcc() any {
	args := u.Called()
	if args.Get(0) == nil {
		return nil
	}
	return dummyBankAcc
}

func (u *BankAccUsecaseMock) FindBankAccByUserID(id uint) ([]*model.BankAccResponse, error) {
	args := u.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.BankAccResponse), nil
}

func (u *BankAccUsecaseMock) FindBankAccByAccountID(id uint) (*model.BankAcc, error) {
	args := u.Called(id)
	if args.Get(0) == nil {
		return nil, errors.New("bank account not found")
	}
	return args.Get(0).(*model.BankAcc), nil
}

func (u *BankAccUsecaseMock) Register(id uint, newBankAcc *model.BankAccResponse) (any, error) {
	args := u.Called(id, newBankAcc)
	if args.Get(0) == nil {
		return nil, errors.New("failed to create data")
	}
	return dummyBankAccResponse, nil
}

func (u *BankAccUsecaseMock) Edit(bankAcc *model.BankAcc) string {
	args := u.Called(bankAcc)
	if args.Get(0) == nil {
		return "failed to update Bank Account"
	}
	return "Bank Account updated Successfully"
}

func (u *BankAccUsecaseMock) UnregAll(userID uint) string {
	args := u.Called(userID)
	if args.Get(0) == nil {
		return "failed to delete Bank Account"
	}
	return "All Bank Account deleted Successfully"
}

func (u *BankAccUsecaseMock) UnregByAccountID(accountID uint) error {
	args := u.Called(accountID)
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

func (suite *BankAccControllerTestSuite) TestFindAllBank_Success() {
	bankAccs := dummyBankAcc
	controller := NewBankAccController(suite.usecaseMock)
	router := setupRouterBankAcc()
	router.GET("user/bank", AuthMiddlewareRole(), controller.FindAllBankAcc)

	suite.usecaseMock.On("FindAllBank").Return(bankAccs)
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "user/bank", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}

func (suite *BankAccControllerTestSuite) SetupTest() {
	suite.routerMock = gin.Default()
	suite.usecaseMock = new(BankAccUsecaseMock)
}

func TestBankAccController(t *testing.T) {
	suite.Run(t, new(BankAccControllerTestSuite))
}
