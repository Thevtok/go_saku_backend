package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/stretchr/testify/assert"
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

type BankAccRepositoryTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSql sqlmock.Sqlmock
}

func (suite *BankAccRepositoryTestSuite) TestGetAll_Success() {
	var users = dummyBankAccResponse[0]
	suite.mockSql.ExpectQuery("SELECT bank_name, account_number, account_holder_name, user_id FROM mst_bank_account").
		WillReturnRows(sqlmock.NewRows([]string{"bank_name", "account_number", "account_holder_name"}).AddRow(users.BankName, users.AccountNumber, users.AccountHolderName))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result := bankAccRepository.GetAll()
	assert.NotNil(suite.T(), result)
}

func (suite *BankAccRepositoryTestSuite) TestGetAll_Failed() {
	suite.mockSql.ExpectQuery("SELECT bank_name, account_number, account_holder_name, user_id FROM mst_bank_account").
		WillReturnError(fmt.Errorf("no data"))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result := bankAccRepository.GetAll()
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "no data", result)
}

func (suite *BankAccRepositoryTestSuite) TestGetByUserID_Success() {
	bankAccs := dummyBankAccResponse[0]
	suite.mockSql.ExpectQuery("SELECT user_id, bank_name, account_number, account_holder_name FROM mst_bank_account WHERE user_id = \\$1").
		WithArgs(bankAccs.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "bank_name", "account_number", "account_holder_name"}).
			AddRow(bankAccs.UserID, bankAccs.BankName, bankAccs.AccountNumber, bankAccs.AccountHolderName))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result, err := bankAccRepository.GetByUserID(bankAccs.UserID)
	assert.NotNil(suite.T(), result)
	assert.Nil(suite.T(), err)
}

func (suite *BankAccRepositoryTestSuite) TestGetByUserID_Failed() {
	bankAccs := dummyBankAccResponse[0]
	suite.mockSql.ExpectQuery("SELECT bank_name, account_number, account_holder_name, user_id FROM mst_bank_account WHERE user_id = \\$1").
		WithArgs(bankAccs.UserID).
		WillReturnError(errors.New("error"))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result, err := bankAccRepository.GetByUserID(1)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *BankAccRepositoryTestSuite) TestGetByAccountID_Success() {
	bankAcc := dummyBankAcc[0]
	suite.mockSql.ExpectQuery("SELECT account_id, bank_name, account_number, account_holder_name, user_id FROM mst_bank_account WHERE account_id = \\$1").
		WithArgs(bankAcc.AccountID).
		WillReturnRows(sqlmock.NewRows([]string{"account_id", "bank_name", "account_number", "account_holder_name", "user_id"}).
			AddRow(bankAcc.AccountID, bankAcc.BankName, bankAcc.AccountNumber, bankAcc.AccountHolderName, bankAcc.UserID))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result, err := bankAccRepository.GetByAccountID(bankAcc.AccountID)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), result)
	assert.Equal(suite.T(), bankAcc, *result)
}

func (suite *BankAccRepositoryTestSuite) TestGetByAccountID_Failed() {
	bankAcc := dummyBankAcc[0]
	suite.mockSql.ExpectQuery("SELECT account_id, bank_name, account_number, account_holder_name, user_id FROM mst_bank_account WHERE account_id = \\$1").
		WithArgs(bankAcc.AccountID).
		WillReturnError(errors.New("bank account not found"))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result, err := bankAccRepository.GetByAccountID(bankAcc.AccountID)
	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "bank account not found", err.Error())
}

func (suite *BankAccRepositoryTestSuite) TestCreateBankA_Success() {
	newBankAcc := dummyBankAccResponse[0]
	userID := uint(1)
	suite.mockSql.ExpectExec("INSERT INTO mst_bank_account \\(user_id, bank_name, account_number, account_holder_name\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)").WithArgs(userID, newBankAcc.BankName, newBankAcc.AccountNumber, newBankAcc.AccountHolderName).WillReturnResult(sqlmock.NewResult(1, 1))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result, err := bankAccRepository.Create(userID, &newBankAcc)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.EqualValues(suite.T(), &newBankAcc, result)
}

func (suite *BankAccRepositoryTestSuite) TestCreate_Failed() {
	newBankAcc := dummyBankAccResponse[0]
	userID := uint(1)
	suite.mockSql.ExpectExec("INSERT INTO mst_bank_account \\(user_id, bank_name, account_number, account_holder_name\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)").
		WithArgs(userID, newBankAcc.BankName, newBankAcc.AccountNumber, newBankAcc.AccountHolderName).WillReturnError(errors.New("failed to create data"))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result, err := bankAccRepository.Create(userID, &newBankAcc)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *BankAccRepositoryTestSuite) TestUpdate_Success() {
	bankAcc := dummyBankAcc[0]
	suite.mockSql.ExpectQuery("SELECT user_id, bank_name, account_number, account_holder_name FROM mst_bank_account WHERE user_id = \\$1").
		WithArgs(bankAcc.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"account_id", "user_id", "bank_name", "account_number", "account_holder_name"}))
	suite.mockSql.ExpectExec("UPDATE mst_bank_account SET bank_name = \\$1, account_number = \\$2, account_holder_name = \\$3 WHERE account_id = \\$4").
		WithArgs(bankAcc.BankName, bankAcc.AccountNumber, bankAcc.AccountHolderName, bankAcc.AccountID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	str := bankAccRepository.Update(&bankAcc)
	assert.NotNil(suite.T(), str)
}

func (suite *BankAccRepositoryTestSuite) TestUpdateScanID_Failed() {
	bankAcc := dummyBankAcc[0]
	expectedError := fmt.Errorf("failed to update Bank Account")
	userID := bankAcc.UserID
	suite.mockSql.ExpectQuery("SELECT user_id, bank_name, account_number, account_holder_name FROM mst_bank_account WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnError(errors.New("user not found"))
	suite.mockSql.ExpectExec("UPDATE mst_bank_account SET bank_name = \\$1, account_number = \\$2, account_holder_name = \\$3 WHERE account_id = \\$4").
		WithArgs(bankAcc.BankName, bankAcc.AccountNumber, bankAcc.AccountHolderName, bankAcc.AccountID).WillReturnError(expectedError)
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	str := bankAccRepository.Update(&bankAcc)
	assert.NotNil(suite.T(), str)
}

func (suite *BankAccRepositoryTestSuite) TestUpdate_Failed() {
	bankAcc := dummyBankAcc[0]
	expectedError := fmt.Errorf("failed to update Bank Account")
	suite.mockSql.ExpectQuery("SELECT user_id, bank_name, account_number, account_holder_name FROM mst_bank_account WHERE user_id = \\$1").
		WithArgs(bankAcc.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"account_id", "user_id", "bank_name", "account_number", "account_holder_name"}))
	suite.mockSql.ExpectExec("UPDATE mst_bank_account SET bank_name = \\$1, account_number = \\$2, account_holder_name = \\$3 WHERE account_id = \\$4").
		WithArgs(bankAcc.BankName, bankAcc.AccountNumber, bankAcc.AccountHolderName, bankAcc.AccountID).WillReturnError(expectedError)
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	str := bankAccRepository.Update(&bankAcc)
	assert.NotNil(suite.T(), str)
}

func (suite *BankAccRepositoryTestSuite) TestDeleteByUserID_Success() {
	userID := uint(1)
	suite.mockSql.ExpectExec("DELETE FROM mst_bank_account WHERE user_id = \\$1").WithArgs(userID).WillReturnResult(sqlmock.NewResult(1, 1))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	str := bankAccRepository.DeleteByUserID(userID)
	assert.EqualValues(suite.T(), "All Bank Account Deleted Successfully", str)
}

func (suite *BankAccRepositoryTestSuite) TestDeleteByUserID_Failed() {
	userID := uint(1)
	expectedError := fmt.Errorf("failed to delete Bank Account")
	suite.mockSql.ExpectExec("DELETE FROM mst_bank_account WHERE user_id = \\$1").WithArgs(userID).WillReturnError(expectedError)
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	str := bankAccRepository.DeleteByUserID(userID)
	assert.NotNil(suite.T(), str)
}

func (suite *BankAccRepositoryTestSuite) TestDeleteByAccountID_Success() {
	bankAcc := dummyBankAcc[0]
	suite.mockSql.ExpectQuery("SELECT account_id, bank_name, account_number, account_holder_name, user_id FROM mst_bank_account WHERE account_id = \\$1").
		WithArgs(bankAcc.AccountID).
		WillReturnRows(sqlmock.NewRows([]string{"account_id", "user_id", "bank_name", "account_number", "account_holder_name"}))
	accountID := uint(1)
	suite.mockSql.ExpectExec("DELETE FROM mst_bank_account WHERE account_id = \\$1").
		WithArgs(accountID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	err := bankAccRepository.DeleteByAccountID(bankAcc.AccountID)
	assert.NotNil(suite.T(), err)
}

func (suite *BankAccRepositoryTestSuite) TestDeleteAccountID_Failed() {
	accountID := dummyBankAcc[0].AccountID
	suite.mockSql.ExpectQuery("SELECT account_id").WithArgs(accountID).WillReturnRows(sqlmock.NewRows([]string{"account_id"}).AddRow(accountID))
	suite.mockSql.ExpectExec("DELETE FROM mst_bank_account WHERE account_id").WithArgs(accountID).WillReturnError(errors.New("Failed to delete AccountID"))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	err := bankAccRepository.DeleteByAccountID(accountID)
	assert.NotNil(suite.T(), err)
}

// Setup test
func (suite *BankAccRepositoryTestSuite) SetupTest() {
	mockDb, mockSql, err := sqlmock.New()
	if err != nil {
		log.Println("Error database", err)
	}
	suite.mockDB = mockDb
	suite.mockSql = mockSql
}

func (suite *BankAccRepositoryTestSuite) TearDownTest() {
	suite.mockDB.Close()
}

func TestBankAccRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(BankAccRepositoryTestSuite))
}
