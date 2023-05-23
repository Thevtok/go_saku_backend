package controller

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyBankAcc = []*model.BankAcc{
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

type BankAccUsecaseMock struct {
	mock.Mock
}

func setupRouterBankAcc() *gin.Engine {
	r := gin.Default()
	return r
}

func (u *BankAccUsecaseMock) FindBankAccByUserID(id string) ([]*model.BankAcc, error) {
	args := u.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return dummyBankAcc, nil
}

func (u *BankAccUsecaseMock) FindBankAccByAccountID(id uint) (*model.BankAcc, error) {
	args := u.Called(id)
	if args.Get(0) == nil {
		return nil, errors.New("bank account not found")
	}
	return args.Get(0).(*model.BankAcc), nil
}

func (u *BankAccUsecaseMock) Register(id string, newBankAcc *model.BankAcc) (any, error) {
	args := u.Called(id, newBankAcc)
	if args.Get(0) == nil {
		return nil, errors.New("failed to create data")
	}
	return dummyBankAcc, nil
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

func (suite *BankAccControllerTestSuite) TestFindBankAccByUserID_Success() {
	bankAcc := dummyBankAcc[:4]
	controller := NewBankAccController(suite.usecaseMock)
	router := setupRouterBankAcc()
	router.GET("/user/bank/:user_id", controller.FindBankAccByUserID)

	suite.usecaseMock.On("FindBankAccByUserID", uint(1)).Return(bankAcc, nil)
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/bank/1", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}

func (suite *BankAccControllerTestSuite) TestFindBankAccByUserID_InvalidUserID() {
	controller := NewBankAccController(suite.usecaseMock)
	router := setupRouterBankAcc()
	router.GET("/user/bank/:user_id", controller.FindBankAccByUserID)

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/bank/invalid_id", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}

func (suite *BankAccControllerTestSuite) TestFinBankAccByUserID_InvalidAccountID() {
	controller := NewBankAccController(suite.usecaseMock)
	router := setupRouterBankAcc()
	router.GET("/user/bank/:user_id", controller.FindBankAccByUserID)
	suite.usecaseMock.On("FindBankAccByUserID", uint(3)).Return(nil, errors.New("Invalid user ID"))

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/bank/3", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}

func (suite *BankAccControllerTestSuite) TestFindBankAccByAccountID_Success() {
	bankAcc := &dummyBankAcc[0]
	controller := NewBankAccController(suite.usecaseMock)
	router := setupRouterBankAcc()
	router.GET("/user/bank/:user_id/:account_id", controller.FindBankAccByAccountID)

	suite.usecaseMock.On("FindBankAccByAccountID", 1).Return(bankAcc, nil)
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/bank/1/1", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}

func (suite *BankAccControllerTestSuite) TestFindBankAccByAccountID_InvalidAccountID() {
	controller := NewBankAccController(suite.usecaseMock)
	router := setupRouterBankAcc()
	router.GET("/user/bank/:user_id/:account_id", controller.FindBankAccByAccountID)

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/bank/1/invalid_ID", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}

func (suite *BankAccControllerTestSuite) TestFindBankAccByAccountID_AccountNotFound() {
	controller := NewBankAccController(suite.usecaseMock)
	router := setupRouterBankAcc()
	router.GET("/user/bank/:user_id/:account_id", controller.FindBankAccByAccountID)

	suite.usecaseMock.On("FindBankAccByAccountID", uint(5)).Return(nil, errors.New("Bank Account not found"))
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/bank/1/5", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}

func (suite *BankAccControllerTestSuite) TestUnregByAccountID_Success() {
	bankAcc := dummyBankAcc[0]
	controller := NewBankAccController(suite.usecaseMock)
	router := setupRouterBankAcc()
	router.DELETE("/user/bank/:user_id/:account_id", controller.UnregByAccountID)
	suite.usecaseMock.On("UnregByAccountID", bankAcc.AccountID).Return(nil)

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, "/user/bank/1/1", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}

func (suite *BankAccControllerTestSuite) TestUnregByAccountId_InvalidAccountID() {
	controller := NewBankAccController(suite.usecaseMock)
	router := setupRouterBankAcc()
	router.DELETE("/user/bank/:user_id/:account_id", controller.UnregByAccountID)

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, "/user/bank/1/invalid_accountID", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}

func (suite *BankAccControllerTestSuite) TestUnregByAccountID_Failed() {
	controller := NewBankAccController(suite.usecaseMock)
	router := setupRouterBankAcc()
	router.DELETE("/user/bank/:user_id/:account_id", controller.UnregByAccountID)

	accountID := uint(1)
	expectedErr := errors.New("failed to delete bank account")
	suite.usecaseMock.On("UnregByAccountID", accountID).Return(expectedErr)

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/user/bank/1/%d", accountID), nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}

func (suite *BankAccControllerTestSuite) SetupTest() {
	suite.routerMock = gin.Default()
	suite.usecaseMock = new(BankAccUsecaseMock)
}

func TestBankAccController(t *testing.T) {
	suite.Run(t, new(BankAccControllerTestSuite))
}
