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

type BankAccRepositoryTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSql sqlmock.Sqlmock
}

func (suite *BankAccRepositoryTestSuite) TestGetByUserID_Success() {
	bankAccs := dummyBankAcc[:2]
	suite.mockSql.ExpectQuery("SELECT user_id, account_id, bank_name, account_number, account_holder_name FROM mst_bank_account WHERE user_id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "account_id", "bank_name", "account_number", "account_holder_name"}).
			AddRow(bankAccs[0].UserID, bankAccs[0].AccountID, bankAccs[0].BankName, bankAccs[0].AccountNumber, bankAccs[0].AccountHolderName).
			AddRow(bankAccs[1].UserID, bankAccs[1].AccountID, bankAccs[1].BankName, bankAccs[1].AccountNumber, bankAccs[1].AccountHolderName))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result, err := bankAccRepository.GetByUserID("1")
	assert.NotNil(suite.T(), result)
	assert.Nil(suite.T(), err)
}

func (suite *BankAccRepositoryTestSuite) TestGetByUserIDScanRows_Failed() {
	bankAccs := dummyBankAcc[:2]
	suite.mockSql.ExpectQuery("SELECT user_id, account_id, bank_name, account_number, account_holder_name FROM mst_bank_account WHERE user_id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "account_id", "bank_name", "account_number", "account_holder_name"}).
			AddRow(nil, bankAccs[0].AccountID, bankAccs[0].BankName, bankAccs[0].AccountNumber, bankAccs[0].AccountHolderName).
			AddRow(bankAccs[1].UserID, bankAccs[1].AccountID, bankAccs[1].BankName, bankAccs[1].AccountNumber, bankAccs[1].AccountHolderName))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result, err := bankAccRepository.GetByUserID("1")
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *BankAccRepositoryTestSuite) TestGetByUserIDButRowsError_Failed() {
	bankAccs := dummyBankAcc[:2]
	expectedErr := fmt.Errorf("Rows error")
	suite.mockSql.ExpectQuery("SELECT user_id, account_id, bank_name, account_number, account_holder_name FROM mst_bank_account WHERE user_id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "account_id", "bank_name", "account_number", "account_holder_name"}).
			AddRow(bankAccs[0].UserID, bankAccs[0].AccountID, bankAccs[0].BankName, bankAccs[0].AccountNumber, bankAccs[0].AccountHolderName).
			RowError(1, expectedErr))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result, err := bankAccRepository.GetByUserID("1")
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

func (suite *BankAccRepositoryTestSuite) TestGetByAccountID_Success() {
	bankAcc := &dummyBankAcc[0]
	suite.mockSql.ExpectQuery("SELECT account_id, bank_name, account_number, account_holder_name, user_id FROM mst_bank_account WHERE account_id = \\$1").
		WithArgs(bankAcc.AccountID).
		WillReturnRows(sqlmock.NewRows([]string{"account_id", "bank_name", "account_number", "account_holder_name", "user_id"}).
			AddRow(bankAcc.AccountID, bankAcc.BankName, bankAcc.AccountNumber, bankAcc.AccountHolderName, bankAcc.UserID))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result, err := bankAccRepository.GetByAccountID(bankAcc.AccountID)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), bankAcc, result)
}

func (suite *BankAccRepositoryTestSuite) TestGetByAccountID_Failed() {
	accountID := uint(4)
	suite.mockSql.ExpectQuery("SELECT account_id, bank_name, account_number, account_holder_name, user_id FROM mst_bank_account WHERE account_id = \\$1").
		WithArgs(accountID).
		WillReturnError(sql.ErrNoRows)
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result, err := bankAccRepository.GetByAccountID(accountID)
	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "bank account not found", err.Error())
}

func (suite *BankAccRepositoryTestSuite) TestGetByAccountID_Error() {
	accountID := uint(4)
	expectedErr := errors.New("unexpected error")
	suite.mockSql.ExpectQuery("SELECT account_id, bank_name, account_number, account_holder_name, user_id FROM mst_bank_account WHERE account_id = \\$1").
		WithArgs(accountID).
		WillReturnError(expectedErr)
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result, err := bankAccRepository.GetByAccountID(accountID)
	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *BankAccRepositoryTestSuite) TestCreate_Success() {
	newBankAcc := dummyBankAcc[0]
	userID := ""
	suite.mockSql.ExpectExec("INSERT INTO mst_bank_account \\(user_id, bank_name, account_number, account_holder_name\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING account_id").
		WithArgs(userID, newBankAcc.BankName, newBankAcc.AccountNumber, newBankAcc.AccountHolderName).
		WillReturnResult(sqlmock.NewResult(1, 1))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result, err := bankAccRepository.Create(userID, &newBankAcc)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.EqualValues(suite.T(), &newBankAcc, result)
}

func (suite *BankAccRepositoryTestSuite) TestCreate_Failed() {
	newBankAcc := dummyBankAcc[0]
	userID := ""
	suite.mockSql.ExpectExec("INSERT INTO mst_bank_account \\(user_id, bank_name, account_number, account_holder_name\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING account_id").
		WithArgs(userID, newBankAcc.BankName, newBankAcc.AccountNumber, newBankAcc.AccountHolderName).
		WillReturnError(errors.New("failed to create data"))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	result, err := bankAccRepository.Create(userID, &newBankAcc)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *BankAccRepositoryTestSuite) TestDeleteByAccountID_Success() {
	bankAcc := dummyBankAcc[0]
	suite.mockSql.ExpectQuery("SELECT account_id, bank_name, account_number, account_holder_name, user_id FROM mst_bank_account WHERE account_id = \\$1").
		WithArgs(bankAcc.AccountID).
		WillReturnRows(sqlmock.NewRows([]string{"account_id", "user_id", "bank_name", "account_number", "account_holder_name"}).
			AddRow(bankAcc.AccountID, bankAcc.BankName, bankAcc.AccountNumber, bankAcc.AccountHolderName, bankAcc.UserID))
	accountID := uint(1)
	suite.mockSql.ExpectExec("DELETE FROM mst_bank_account WHERE account_id = \\$1").
		WithArgs(accountID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	err := bankAccRepository.DeleteByAccountID(bankAcc.AccountID)
	assert.Nil(suite.T(), err)
}

func (suite *BankAccRepositoryTestSuite) TestDeleteByAccountIDScan_Failed() {
	accountID := uint(4)
	suite.mockSql.ExpectQuery("SELECT account_id, bank_name, account_number, account_holder_name, user_id FROM mst_bank_account WHERE account_id = \\$1").
		WithArgs(accountID).
		WillReturnError(sql.ErrNoRows)
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	err := bankAccRepository.DeleteByAccountID(accountID)
	assert.NotNil(suite.T(), err)
}

func (suite *BankAccRepositoryTestSuite) TestDeleteByAccountID_Failed() {
	bankAcc := dummyBankAcc[0]
	suite.mockSql.ExpectQuery("SELECT account_id, bank_name, account_number, account_holder_name, user_id FROM mst_bank_account WHERE account_id = \\$1").
		WithArgs(bankAcc.AccountID).
		WillReturnRows(sqlmock.NewRows([]string{"account_id", "user_id", "bank_name", "account_number", "account_holder_name"}).
			AddRow(bankAcc.AccountID, bankAcc.BankName, bankAcc.AccountNumber, bankAcc.AccountHolderName, bankAcc.UserID))
	accountID := uint(1)
	suite.mockSql.ExpectExec("DELETE FROM mst_bank_account WHERE account_id = \\$1").
		WithArgs(accountID).
		WillReturnError(fmt.Errorf("failed to delete bank account with ID %d", accountID))
	bankAccRepository := NewBankAccRepository(suite.mockDB)
	err := bankAccRepository.DeleteByAccountID(bankAcc.AccountID)
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
