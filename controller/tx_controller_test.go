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

func setupRouterTx() *gin.Engine {
	r := gin.Default()
	return r
}

type MockUtils struct {
	mock.Mock
}

type TxUseCaseMock struct {
	mock.Mock
}

type TxControllerTestSuite struct {
	suite.Suite
	routerMock    *gin.Engine
	bankCaseMock  *BankAccUsecaseMock
	cardCaseMock  *CardUsecaseMock
	txUsecaseMock *TxUseCaseMock

	useCaseMock *UserUseCaseMock
}

func (suite *TxControllerTestSuite) SetupTest() {
	suite.routerMock = gin.Default()
	suite.useCaseMock = new(UserUseCaseMock)
	suite.bankCaseMock = new(BankAccUsecaseMock)
	suite.cardCaseMock = new(CardUsecaseMock)
	suite.txUsecaseMock = new(TxUseCaseMock)

}

func TestTxController(t *testing.T) {
	suite.Run(t, new(TxControllerTestSuite))
}

func (m *TxUseCaseMock) FindByPeId(id uint) ([]*model.PointExchange, error) {
	args := m.Called(id)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.PointExchange), nil
}

func (uc *TxUseCaseMock) FindTxById(txId uint) ([]*model.Transaction, error) {
	args := uc.Called(txId)
	if args[0] == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*model.Transaction), args.Error(1)
}

func (uc *TxUseCaseMock) CreateDepositBank(transaction *model.TransactionBank) error {
	args := uc.Called(transaction)
	if err, ok := args.Get(0).(error); ok {
		return err
	}
	return nil
}

func (uc *TxUseCaseMock) CreateDepositCard(transaction *model.TransactionCard) error {
	args := uc.Called(transaction)
	if err, ok := args.Get(0).(error); ok {
		return err
	}

	return nil
}
func (uc *TxUseCaseMock) CreateWithdrawal(transaction *model.TransactionWithdraw) error {
	args := uc.Called(transaction)
	if err, ok := args.Get(0).(error); ok {
		return err
	}

	return nil
}
func (uc *TxUseCaseMock) CreateTransfer(sender *model.User, recipient *model.User, amount uint) (any, error) {
	args := uc.Called(sender, recipient, amount)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), nil
}

func (uc *TxUseCaseMock) CreateRedeem(transaction *model.TransactionPoint) error {
	args := uc.Called(transaction)
	if err, ok := args.Get(0).(error); ok {
		return err
	}

	return nil
}

func (suite *TxControllerTestSuite) TestCreateDepositBank_Success() {

	bank := &model.TransactionBank{
		SenderID:      uint(1),
		BankAccountID: uint(1),
		Amount:        50000,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.bankCaseMock.On("FindBankAccByAccountID", uint(1)).Return(&model.BankAcc{
		AccountID: uint(1),
		UserID:    uint(1),
	}, nil)

	suite.txUsecaseMock.On("CreateDepositBank", bank).Return(nil)

	router.POST("/user/tx/depo/bank/:user_id/:bank_account_id", controller.CreateDepositBank)
	reqBody, _ := json.Marshal(bank)

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/depo/bank/1/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusCreated, r.Code)
}
func (suite *TxControllerTestSuite) TestCreateDepositCard_Success() {

	card := &model.TransactionCard{
		SenderID: uint(1),
		CardID:   uint(1),
		Amount:   50000,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.cardCaseMock.On("FindCardByCardID", uint(1)).Return(&model.Card{
		CardID: uint(1),
		UserID: uint(1),
	}, nil)

	suite.txUsecaseMock.On("CreateDepositCard", card).Return(nil)

	router.POST("/user/tx/depo/card/:user_id/:card_id", controller.CreateDepositCard)
	reqBody, _ := json.Marshal(card)

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/depo/card/1/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusCreated, r.Code)
}

func (suite *TxControllerTestSuite) TestCreateWithdrawal_Success() {

	wd := &model.TransactionWithdraw{
		SenderID:      uint(1),
		BankAccountID: uint(1),
		Amount:        50000,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.bankCaseMock.On("FindBankAccByAccountID", uint(1)).Return(&model.BankAcc{
		AccountID: uint(1),
		UserID:    uint(1),
	}, nil)

	suite.txUsecaseMock.On("CreateWithdrawal", wd).Return(nil)

	router.POST("/user/tx/wd/:user_id/:bank_account_id", controller.CreateWithdrawal)
	reqBody, _ := json.Marshal(wd)

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/wd/1/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusCreated, r.Code)
}

func (suite *TxControllerTestSuite) TestCreateTransfer_Success() {

	tf := &model.TransactionTransfer{
		SenderID:    uint(1),
		RecipientID: uint(2),
		Amount:      50000,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.useCaseMock.On("FindById", tf.SenderID).Return(&model.User{
		ID: uint(1),
	}, nil)
	suite.useCaseMock.On("FindById", tf.RecipientID).Return(&model.User{
		ID: uint(2),
	}, nil)

	suite.txUsecaseMock.On("CreateTransfer", &model.User{
		ID: uint(1),
	}, &model.User{
		ID: uint(2),
	}, tf.Amount).Return(tf, nil)

	router.POST("/user/tx/tf/:user_id", controller.CreateTransferTransaction)
	reqBody, _ := json.Marshal(tf)

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/tf/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusCreated, r.Code)
}

func (suite *TxControllerTestSuite) TestCreateRedeem_Success() {

	redeem := &model.TransactionPoint{
		SenderID: uint(1),

		PointExchangeID: 1,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.txUsecaseMock.On("FindByPeId", uint(1)).Return([]*model.PointExchange{}, nil)

	suite.txUsecaseMock.On("CreateRedeem", redeem).Return(nil)

	router.POST("/user/tx/redeem/:user_id/:pe_id", controller.CreateRedeemTransaction)
	reqBody, _ := json.Marshal(redeem)

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/redeem/1/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusCreated, r.Code)
}

func (suite *TxControllerTestSuite) TestGetTxBySenderId_Success() {

	tx := &model.Transaction{
		SenderID: uint(1),
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.useCaseMock.On("FindById", tx.SenderID).Return(&model.User{
		ID: uint(1),
	}, nil)

	suite.txUsecaseMock.On("FindTxById", tx.SenderID).Return([]*model.Transaction{}, nil)

	router.GET("/user/tx/:user_id", controller.GetTxBySenderId)
	reqBody, _ := json.Marshal(tx)

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/tx/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusOK, r.Code)

}

func (suite *TxControllerTestSuite) TestCreateDepositBank_InvalidUserID() {
	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	invalidUserID := "invalidUserID"
	router.POST("/user/tx/depo/bank/:user_id/:bank_account_id", controller.CreateDepositBank)
	reqBody, _ := json.Marshal(&model.TransactionBank{
		SenderID:      uint(1),
		BankAccountID: uint(1),
		Amount:        50000,
	})

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/user/tx/depo/bank/%s/1", invalidUserID), bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *TxControllerTestSuite) TestCreateDepositCard_InvalidUserID() {
	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	invalidUserID := "invalidUserID"
	router.POST("/user/tx/depo/card/:user_id/:card_id", controller.CreateDepositCard)
	reqBody, _ := json.Marshal(&model.TransactionCard{
		SenderID: uint(1),
		CardID:   uint(1),
		Amount:   50000,
	})

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/user/tx/depo/card/%s/1", invalidUserID), bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}
func (suite *TxControllerTestSuite) TestCreateWithdrawal_InvalidUserID() {
	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	invalidUserID := "invalidUserID"
	router.POST("/user/tx/wd/:user_id/:bank_account_id", controller.CreateWithdrawal)
	reqBody, _ := json.Marshal(&model.TransactionWithdraw{
		SenderID:      uint(1),
		BankAccountID: uint(1),
		Amount:        50000,
	})

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/user/tx/wd/%s/1", invalidUserID), bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *TxControllerTestSuite) TestCreateRedeem_InvalidUserID() {
	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	invalidUserID := "invalidUserID"
	router.POST("/user/tx/redeem/:user_id/:pe_id", controller.CreateRedeemTransaction)
	reqBody, _ := json.Marshal(&model.TransactionPoint{
		SenderID:        uint(1),
		PointExchangeID: 1,
	})

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/user/tx/redeem/%s/1", invalidUserID), bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *TxControllerTestSuite) TestCreateTransfer_InvalidUserID() {
	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	invalidUserID := "invalidUserID"
	router.POST("/user/tx/tf/:user_id", controller.CreateTransferTransaction)
	reqBody, _ := json.Marshal(&model.TransactionTransfer{
		SenderID: uint(1),

		Amount: 50000,
	})

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/user/tx/tf/%s", invalidUserID), bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *TxControllerTestSuite) TestGetTx_InvalidUserID() {
	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	invalidUserID := "invalidUserID"
	router.GET("/user/tx/:user_id", controller.GetTxBySenderId)
	reqBody, _ := json.Marshal(&model.Transaction{
		SenderID: uint(1),
	})

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user/tx/%s", invalidUserID), bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *TxControllerTestSuite) TestCreateDepositBank_InvalidBank_Account_ID() {
	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	invalidBankAccountID := "InvalidBankAccountID"
	router.POST("/user/tx/depo/bank/:user_id/:bank_account_id", controller.CreateDepositBank)
	reqBody, _ := json.Marshal(&model.TransactionBank{
		SenderID:      uint(1),
		BankAccountID: uint(1),
		Amount:        50000,
	})

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/user/tx/depo/bank/1/%s", invalidBankAccountID), bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}
func (suite *TxControllerTestSuite) TestCreateDepositCard_InvalidCard_ID() {
	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	invalidCardID := "InvalidCardID"
	router.POST("/user/tx/depo/card/:user_id/:card_id", controller.CreateDepositCard)
	reqBody, _ := json.Marshal(&model.TransactionCard{
		SenderID: uint(1),
		CardID:   uint(1),
		Amount:   50000,
	})

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/user/tx/depo/card/1/%s", invalidCardID), bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}
func (suite *TxControllerTestSuite) TestCreateWithdrawal_InvalidBank_Account_ID() {
	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	invalidBankAccountID := "InvalidBankAccountID"
	router.POST("/user/tx/wd/:user_id/:bank_account_id", controller.CreateWithdrawal)
	reqBody, _ := json.Marshal(&model.TransactionWithdraw{
		SenderID:      uint(1),
		BankAccountID: uint(1),
		Amount:        50000,
	})

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/user/tx/wd/1/%s", invalidBankAccountID), bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *TxControllerTestSuite) TestCreateRedeem_InvalidPeID() {
	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	invalidPeiD := "invalidPeID"
	router.POST("/user/tx/redeem/:user_id/:pe_id", controller.CreateRedeemTransaction)
	reqBody, _ := json.Marshal(&model.TransactionPoint{
		SenderID:        uint(1),
		PointExchangeID: 1,
	})

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/user/tx/redeem/1/%s", invalidPeiD), bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *TxControllerTestSuite) TestCreateDepositBank_BankAccIDNotFound() {
	// Prepare
	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.bankCaseMock.On("FindBankAccByAccountID", uint(1)).Return(nil, errors.New("Bank_account_id not found"))

	router.POST("/user/tx/depo/bank/:user_id/:bank_account_id", controller.CreateDepositBank)
	reqBody, _ := json.Marshal(&model.TransactionBank{
		SenderID:      uint(1),
		BankAccountID: uint(1),
		Amount:        50000,
	})

	// Execute
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/depo/bank/1/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
}
func (suite *TxControllerTestSuite) TestCreateDepositCard_CardIdNotFound() {
	// Prepare
	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.cardCaseMock.On("FindCardByCardID", uint(1)).Return(nil, errors.New("Card_id not found"))

	router.POST("/user/tx/depo/card/:user_id/:card_id", controller.CreateDepositCard)
	reqBody, _ := json.Marshal(&model.TransactionCard{
		SenderID: uint(1),
		CardID:   uint(1),
		Amount:   50000,
	})

	// Execute
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/depo/card/1/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
}

func (suite *TxControllerTestSuite) TestCreateWithdrawal_BankAccIDNotFound() {
	// Prepare
	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.bankCaseMock.On("FindBankAccByAccountID", uint(1)).Return(nil, errors.New("Bank_account_id not found"))

	router.POST("/user/tx/wd/:user_id/:bank_account_id", controller.CreateWithdrawal)
	reqBody, _ := json.Marshal(&model.TransactionWithdraw{
		SenderID:      uint(1),
		BankAccountID: uint(1),
		Amount:        50000,
	})

	// Execute
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/wd/1/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
}

func (suite *TxControllerTestSuite) TestCreaRedeem_PeIDNotFound() {
	// Prepare
	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.txUsecaseMock.On("FindByPeId", uint(1)).Return(nil, errors.New("Pe_iD not found"))

	router.POST("/user/tx/redeem/:user_id/:pe_id", controller.CreateRedeemTransaction)
	reqBody, _ := json.Marshal(&model.TransactionPoint{
		SenderID:        uint(1),
		PointExchangeID: 1,
	})

	// Execute
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/redeem/1/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
}
func (suite *TxControllerTestSuite) TestCreateTransfer_FailedToGetSenderUser() {
	// Prepare
	tf := &model.TransactionTransfer{
		SenderID:    uint(1),
		RecipientID: uint(2),
		Amount:      50000,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()
	suite.useCaseMock.On("FindById", uint(1)).Return(nil, errors.New("failed to get sender user"))

	router.POST("/user/tx/tf/:user_id", controller.CreateTransferTransaction)
	reqBody, _ := json.Marshal(tf)

	// Execute
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/tf/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
}

func (suite *TxControllerTestSuite) TestGetTx_FailedToGetSenderUser() {
	// Prepare
	tx := &model.Transaction{
		SenderID: uint(1),
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()
	suite.useCaseMock.On("FindById", uint(1)).Return(nil, errors.New("failed to get sender user"))

	router.GET("/user/tx/:user_id", controller.GetTxBySenderId)
	reqBody, _ := json.Marshal(tx)

	// Execute
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/tx/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
}

func (suite *TxControllerTestSuite) TestCreateTransfer_FailedToGetRecipientUser() {
	// Prepare
	tf := &model.TransactionTransfer{
		SenderID:    uint(1),
		RecipientID: uint(2),
		Amount:      50000,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()
	suite.useCaseMock.On("FindById", uint(1)).Return(&model.User{
		ID: 1,
	}, nil)
	suite.useCaseMock.On("FindById", uint(2)).Return(nil, errors.New("failed to get recipient user"))

	router.POST("/user/tx/tf/:user_id", controller.CreateTransferTransaction)
	reqBody, _ := json.Marshal(tf)

	// Execute
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/tf/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
}
func (suite *TxControllerTestSuite) TestCreateDepositBank_BankAccountNotBelongToGivenUser() {
	// Prepare
	bank := &model.TransactionBank{
		SenderID:      uint(1),
		BankAccountID: uint(1),
		Amount:        50000,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.bankCaseMock.On("FindBankAccByAccountID", uint(1)).Return(&model.BankAcc{
		AccountID: uint(1),
		UserID:    uint(2),
	}, nil)

	router.POST("/user/tx/depo/bank/:user_id/:bank_account_id", controller.CreateDepositBank)
	reqBody, _ := json.Marshal(bank)

	// Execute
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/depo/bank/1/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *TxControllerTestSuite) TestCreateDepositCard_CardAccountNotBelongToGivenUser() {
	// Prepare
	card := &model.TransactionCard{
		SenderID: uint(1),
		CardID:   uint(1),
		Amount:   50000,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.cardCaseMock.On("FindCardByCardID", uint(1)).Return(&model.Card{

		UserID: uint(5),
		CardID: uint(1),
	}, nil)

	router.POST("/user/tx/depo/card/:user_id/:card_id", controller.CreateDepositCard)
	reqBody, _ := json.Marshal(card)

	// Execute
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/depo/card/5/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *TxControllerTestSuite) TestCreateWithdrawal_BankAccountNotBelongToGivenUser() {
	// Prepare
	wd := &model.TransactionWithdraw{
		SenderID:      uint(1),
		BankAccountID: uint(1),
		Amount:        50000,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.bankCaseMock.On("FindBankAccByAccountID", uint(1)).Return(&model.BankAcc{
		AccountID: uint(1),
		UserID:    uint(2),
	}, nil)

	router.POST("/user/tx/wd/:user_id/:bank_account_id", controller.CreateWithdrawal)
	reqBody, _ := json.Marshal(wd)

	// Execute
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/wd/1/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *TxControllerTestSuite) TestCreateDepositBank_InvalidRequestBody() {
	// Prepare
	bank := &model.TransactionBank{
		SenderID:      uint(1),
		BankAccountID: uint(1),
		Amount:        50000,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.bankCaseMock.On("FindBankAccByAccountID", uint(1)).Return(&model.BankAcc{
		AccountID: uint(1),
		UserID:    uint(1),
	}, nil)

	suite.txUsecaseMock.On("CreateDepositBank", bank).Return(nil)

	router.POST("/user/tx/depo/bank/:user_id/:bank_account_id", controller.CreateDepositBank)
	reqBody := []byte("invalid JSON")

	// Execute
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/depo/bank/1/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *TxControllerTestSuite) TestCreateDepositCard_InvalidRequestBody() {
	// Prepare
	card := &model.TransactionCard{
		SenderID: uint(1),
		CardID:   uint(1),
		Amount:   50000,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.cardCaseMock.On("FindCardByCardID", uint(1)).Return(&model.Card{
		CardID: uint(1),
		UserID: uint(1),
	}, nil)

	suite.txUsecaseMock.On("CreateDepositCard", card).Return(nil)

	router.POST("/user/tx/depo/card/:user_id/:card_id", controller.CreateDepositCard)
	reqBody := []byte("invalid JSON")

	// Execute
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/depo/card/1/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *TxControllerTestSuite) TestCreateWithdrawal_InvalidRequestBody() {
	// Prepare
	wd := &model.TransactionWithdraw{
		SenderID:      uint(1),
		BankAccountID: uint(1),
		Amount:        50000,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.bankCaseMock.On("FindBankAccByAccountID", uint(1)).Return(&model.BankAcc{
		AccountID: uint(1),
		UserID:    uint(1),
	}, nil)

	suite.txUsecaseMock.On("CreateWithdrawal", wd).Return(nil)

	router.POST("/user/tx/wd/:user_id/:bank_account_id", controller.CreateWithdrawal)
	reqBody := []byte("invalid JSON")

	// Execute
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/wd/1/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *TxControllerTestSuite) TestCreateTransfer_InvalidRequestBody() {
	// Prepare
	tf := &model.TransactionTransfer{
		SenderID:    uint(1),
		RecipientID: uint(2),
		Amount:      50000,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.txUsecaseMock.On("CreateTransfer", &model.User{
		ID: 1,
	}, &model.User{
		ID: 2,
	}, tf.Amount).Return(tf, nil)

	router.POST("/user/tx/tf/:user_id", controller.CreateTransferTransaction)
	reqBody := []byte("invalid JSON")

	// Execute
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/tf/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *TxControllerTestSuite) TestCreateRedeem_InvalidRequestBody() {
	// Prepare
	redeem := &model.TransactionPoint{
		SenderID:        uint(1),
		PointExchangeID: 1,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()
	suite.txUsecaseMock.On("FindByPeId", uint(1)).Return([]*model.PointExchange{}, nil)
	suite.txUsecaseMock.On("CreateRedeem", redeem).Return(nil)

	router.POST("/user/tx/redeem/:user_id/:pe_id", controller.CreateRedeemTransaction)
	reqBody := []byte("invalid JSON")

	// Execute
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/redeem/1/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}
func (suite *TxControllerTestSuite) TestCreateDepositBank_Failed() {
	bank := &model.TransactionBank{
		SenderID:      uint(1),
		BankAccountID: uint(1),
		Amount:        50000,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.bankCaseMock.On("FindBankAccByAccountID", uint(1)).Return(&model.BankAcc{
		AccountID: uint(1),
		UserID:    uint(1),
	}, nil)

	err := errors.New("failed to create deposit transaction")
	suite.txUsecaseMock.On("CreateDepositBank", bank).Return(err)

	router.POST("/user/tx/depo/bank/:user_id/:bank_account_id", controller.CreateDepositBank)
	reqBody, _ := json.Marshal(bank)

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/depo/bank/1/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)
	assert.Contains(suite.T(), r.Body.String(), "Failed to create Deposit Transaction")
}

func (suite *TxControllerTestSuite) TestCreateDepositCard_Failed() {
	card := &model.TransactionCard{
		SenderID: uint(1),
		CardID:   uint(1),
		Amount:   50000,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.cardCaseMock.On("FindCardByCardID", uint(1)).Return(&model.Card{
		CardID: uint(1),
		UserID: uint(1),
	}, nil)

	err := errors.New("failed to create deposit transaction")
	suite.txUsecaseMock.On("CreateDepositCard", card).Return(err)

	router.POST("/user/tx/depo/card/:user_id/:card_id", controller.CreateDepositCard)
	reqBody, _ := json.Marshal(card)

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/depo/card/1/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)
	assert.Contains(suite.T(), r.Body.String(), "Failed to create Deposit Transaction")
}
func (suite *TxControllerTestSuite) TestCreateWithdrawal_InsufficientBalance() {
	wd := &model.TransactionWithdraw{
		SenderID:      uint(1),
		BankAccountID: uint(1),
		Amount:        50000,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.bankCaseMock.On("FindBankAccByAccountID", uint(1)).Return(&model.BankAcc{
		AccountID: uint(1),
		UserID:    uint(1),
	}, nil)

	suite.txUsecaseMock.On("CreateWithdrawal", wd).Return(errors.New("insufficient balance"))

	router.POST("/user/tx/wd/:user_id/:bank_account_id", controller.CreateWithdrawal)
	reqBody, _ := json.Marshal(wd)

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/wd/1/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusUnprocessableEntity, r.Code)
}
func (suite *TxControllerTestSuite) TestCreateWithdrawal_InternalServerError() {
	wd := &model.TransactionWithdraw{
		SenderID:      uint(1),
		BankAccountID: uint(1),
		Amount:        50000,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.bankCaseMock.On("FindBankAccByAccountID", uint(1)).Return(&model.BankAcc{
		AccountID: uint(1),
		UserID:    uint(1),
	}, nil)

	suite.txUsecaseMock.On("CreateWithdrawal", wd).Return(errors.New("error"))

	router.POST("/user/tx/wd/:user_id/:bank_account_id", controller.CreateWithdrawal)
	reqBody, _ := json.Marshal(wd)

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/wd/1/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)
}

func (suite *TxControllerTestSuite) TestCreateTransfer_Failed() {
	// Prepare
	tf := &model.TransactionTransfer{
		SenderID:    uint(1),
		RecipientID: uint(2),
		Amount:      50000,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()
	suite.useCaseMock.On("FindById", uint(1)).Return(&model.User{
		ID: 1,
	}, nil)
	suite.useCaseMock.On("FindById", uint(2)).Return(&model.User{
		ID: 2,
	}, nil)
	suite.txUsecaseMock.On("CreateTransfer", &model.User{
		ID: uint(1),
	}, &model.User{
		ID: uint(2),
	}, tf.Amount).Return(nil, errors.New("err"))

	router.POST("/user/tx/tf/:user_id", controller.CreateTransferTransaction)
	reqBody, _ := json.Marshal(tf)

	// Execute
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/tf/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	// Assert
	assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)
}

func (suite *TxControllerTestSuite) TestCreateRedeem_Failed() {

	redeem := &model.TransactionPoint{
		SenderID: uint(1),

		PointExchangeID: 1,
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.txUsecaseMock.On("FindByPeId", uint(1)).Return([]*model.PointExchange{}, nil)

	suite.txUsecaseMock.On("CreateRedeem", redeem).Return(errors.New("err"))

	router.POST("/user/tx/redeem/:user_id/:pe_id", controller.CreateRedeemTransaction)
	reqBody, _ := json.Marshal(redeem)

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/user/tx/redeem/1/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)
}

func (suite *TxControllerTestSuite) TestGetTxBySenderId_Failed() {

	tx := &model.Transaction{
		SenderID: uint(1),
	}

	controller := NewTransactionController(suite.txUsecaseMock, suite.useCaseMock, suite.bankCaseMock, suite.cardCaseMock)
	router := setupRouterTx()

	suite.useCaseMock.On("FindById", tx.SenderID).Return(&model.User{
		ID: uint(1),
	}, nil)

	suite.txUsecaseMock.On("FindTxById", tx.SenderID).Return(nil, errors.New("err"))

	router.GET("/user/tx/:user_id", controller.GetTxBySenderId)
	reqBody, _ := json.Marshal(tx)

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/tx/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)

}
