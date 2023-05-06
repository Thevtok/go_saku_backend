package controller

import (
	"bytes"
	"encoding/json"
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
	return dummyBankAccResponse1, nil
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
	router.GET("/user/bank", controller.FindAllBankAcc)

	suite.usecaseMock.On("FindAllBankAcc").Return(bankAccs)
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/bank", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}

func (suite *BankAccControllerTestSuite) TestFindAllBankAcc_Failed() {
	controller := NewBankAccController(suite.usecaseMock)
	router := setupRouterBankAcc()
	router.GET("/user/bank", controller.FindAllBankAcc)

	suite.usecaseMock.On("FindAllBankAcc").Return(nil)
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/bank", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}

func (suite *BankAccControllerTestSuite) TestFindBankAccByUserID_Success() {
	bankAcc := dummyBankAccResponse1[:4]
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

	suite.usecaseMock.On("FindBankAccByAccountID", bankAcc.AccountID).Return(bankAcc, nil)
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

func (suite *BankAccControllerTestSuite) TestEdit_Success() {
	bankAcc := dummyBankAcc[0]
	controller := NewBankAccController(suite.usecaseMock)
	router := setupRouterBankAcc()
	router.PUT("/user/bank/update/:user_id/:account_id", controller.Edit)

	suite.usecaseMock.On("FindBankAccByAccountID", uint(1)).Return(&bankAcc, nil)
	suite.usecaseMock.On("Edit", mock.Anything).Return("Bank Account updated Successfully")
	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(bankAcc)
	request, _ := http.NewRequest(http.MethodPut, "/user/bank/update/1/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)
	response := r.Body.String()
	var result model.BankAcc
	json.Unmarshal([]byte(response), &result)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}

func (suite *BankAccControllerTestSuite) TestEdit_InvalidAccountID() {
	bankAcc := dummyBankAcc[0]
	controller := NewBankAccController(suite.usecaseMock)
	router := setupRouterBankAcc()
	router.PUT("/user/bank/update/:user_id/:account_id", controller.Edit)

	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(bankAcc)
	request, _ := http.NewRequest(http.MethodPut, "/user/bank/update/1/invalid_id", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)
	response := r.Body.String()
	var result model.BankAcc
	json.Unmarshal([]byte(response), &result)
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}

func (suite *BankAccControllerTestSuite) TestEdit_AccountNotFound() {
	bankAcc := dummyBankAcc[0]
	controller := NewBankAccController(suite.usecaseMock)
	router := setupRouterBankAcc()
	router.PUT("/user/bank/update/:user_id/:account_id", controller.Edit)

	suite.usecaseMock.On("FindBankAccByAccountID", uint(5)).Return(nil, errors.New("Bank Account not found"))
	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(bankAcc)
	request, _ := http.NewRequest(http.MethodPut, "/user/bank/update/1/5", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)
	response := r.Body.String()
	var result model.Card
	json.Unmarshal([]byte(response), &result)
	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}

// func (suite *BankAccControllerTestSuite) TestEdit_Failed() {
// 	bankAcc := dummyBankAcc[0]
// 	controller := NewBankAccController(suite.usecaseMock)
// 	router := setupRouterBankAcc()
// 	router.PUT("/user/bank/update/:user_id/:account_id", controller.Edit)

// 	suite.usecaseMock.On("FindBankAccByAccountID", uint(1)).Return(&bankAcc, nil)
// 	expectedErr := errors.New("Failed to edit Bank")
// 	suite.usecaseMock.On("Edit", mock.Anything).Return(expectedErr)
// 	reqBody, _ := json.Marshal(bankAcc)
// 	request, _ := http.NewRequest(http.MethodPut, "/user/bank/update/1/1", bytes.NewBuffer(reqBody))
// 	r := httptest.NewRecorder()
// 	router.ServeHTTP(r, request)

// 	assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)
// 	suite.usecaseMock.AssertExpectations(suite.T())
// }

func (suite *BankAccControllerTestSuite) TestUnregAll_Success() {
	bankAcc := dummyBankAcc[0]
	controller := NewBankAccController(suite.usecaseMock)
	router := setupRouterBankAcc()
	router.DELETE("/user/bank/:user_id", controller.UnregAll)
	suite.usecaseMock.On("UnregAll", bankAcc.UserID).Return("Bank Account deleted Successfully")

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, "/user/bank/1", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}

func (suite *BankAccControllerTestSuite) TestUnregAll_InvalidUserID() {
	controller := NewBankAccController(suite.usecaseMock)
	router := setupRouterBankAcc()
	router.DELETE("/user/bank/:user_id", controller.UnregAll)

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, "/user/bank/invalid_userID", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
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
