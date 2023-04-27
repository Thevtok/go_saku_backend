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
	suite.mockSql.ExpectQuery("SELECT bank_name, account_number, account_holder_name, user_id FROM mst_bank_account").WillReturnRows(sqlmock.NewRows([]string{"bank_name", "account_number", "account_holder_name"}).AddRow(users.BankName, users.AccountNumber, users.AccountHolderName))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result := bankAccRepository.GetAll()
	assert.NotNil(suite.T(), result)
}

func (suite *BankAccRepositoryTestSuite) TestGetAll_Failed() {
	suite.mockSql.ExpectQuery("SELECT bank_name, account_number, account_holder_name, user_id FROM mst_bank_account").WillReturnError(fmt.Errorf("no data"))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result := bankAccRepository.GetAll()
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "no data", result)
}

func (suite *BankAccRepositoryTestSuite) TestGetByUserID_Success() {
	bankAccs := dummyBankAccResponse[0]
	suite.mockSql.ExpectQuery("SELECT user_id, bank_name, account_number, account_holder_name FROM mst_bank_account WHERE user_id = \\$1").WithArgs(bankAccs.UserID).WillReturnRows(sqlmock.NewRows([]string{"user_id", "bank_name", "account_number", "account_holder_name"}).AddRow(bankAccs.UserID, bankAccs.BankName, bankAccs.AccountNumber, bankAccs.AccountHolderName))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result, err := bankAccRepository.GetByUserID(bankAccs.UserID)
	assert.NotNil(suite.T(), result)
	assert.Nil(suite.T(), err)
}

func (suite *BankAccRepositoryTestSuite) TestGetByUserID_Failed() {
	suite.mockSql.ExpectQuery("SELECT bank_name, account_number, account_holder_name, user_id FROM mst_bank_account WHERE user_id = \\$1").WithArgs(1).WillReturnError(errors.New("some error"))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result, err := bankAccRepository.GetByUserID(1)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *BankAccRepositoryTestSuite) TestGetByAccountID_Success() {
	bankAcc := dummyBankAcc[0]
	suite.mockSql.ExpectQuery("SELECT account_id, bank_name, account_number, account_holder_name, user_id FROM mst_bank_account WHERE account_id = \\$1").WithArgs(bankAcc.AccountID).WillReturnError(sql.ErrNoRows)
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result, err := bankAccRepository.GetByAccountID(bankAcc.AccountID)
	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "Bank Account not found", err.Error())
}

func

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
