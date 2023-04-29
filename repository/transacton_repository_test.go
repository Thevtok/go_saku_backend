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

var dummyTxBank = []*model.TransactionBank{
	{},
}
var dummyTxCard = []*model.TransactionCard{
	{},
}
var dummyTxWd = []*model.TransactionWithdraw{
	{},
}
var dummyTxRd = []*model.TransactionPoint{
	{},
}
var dummyTxTf = []*model.TransactionTransfer{
	{},
}

var value = uint(2)
var value1 = uint(3)
var txs = []*model.Transaction{
	{
		SenderID:        1,
		TransactionType: "Transfer",
		RecipientID:     &value,
		BankAccountID:   &value,
		CardID:          &value,
		PointExchangeID: &value,
		Point:           &value,
		Amount:          &value,
		TransactionDate: now,
	},
	{
		SenderID:        1,
		TransactionType: "Transfer",
		RecipientID:     &value1,
		BankAccountID:   &value1,
		CardID:          &value1,
		PointExchangeID: &value1,
		Point:           &value1,
		Amount:          &value1,
		TransactionDate: now,
	},
}
var txP = []*model.TransactionPoint{
	{
		SenderID:        1,
		TransactionType: "Redeem",

		PointExchangeID: 1,
		Point:           1,

		TransactionDate: now,
	},
}

type TransactionRepositoryTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSql sqlmock.Sqlmock
}

func (suite *TransactionRepositoryTestSuite) TestGetBySenderId_Success() {
	// Create some test data
	senderID := txs[0].SenderID

	// Create a mock database connection and repository

	repo := NewTxRepository(suite.mockDB)

	// Expect the query to be executed with the correct arguments
	suite.mockSql.ExpectQuery("SELECT").WithArgs(senderID).WillReturnRows(
		sqlmock.NewRows([]string{"transaction_type", "sender_id", "recipient_id", "bank_account_id", "card_id", "pe_id", "amount", "point", "transaction_date"}).
			AddRow(txs[0].TransactionType, txs[0].SenderID, txs[0].RecipientID, txs[0].BankAccountID, txs[0].CardID, txs[0].PointExchangeID, txs[0].Amount, txs[0].Point, txs[0].TransactionDate).
			AddRow(txs[1].TransactionType, txs[1].SenderID, txs[1].RecipientID, txs[1].BankAccountID, txs[1].CardID, txs[1].PointExchangeID, txs[1].Amount, txs[1].Point, txs[1].TransactionDate))

	// Call the GetBySenderId method
	result, err := repo.GetBySenderId(senderID)

	// Assert that no errors occurred and all expectations were met
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), len(txs), len(result))
	assert.NotNil(suite.T(), result)

	assert.NoError(suite.T(), suite.mockSql.ExpectationsWereMet())
}

func (suite *TransactionRepositoryTestSuite) TestGetBySenderId_Error() {
	senderID := uint(1)
	expectedError := fmt.Errorf("dummy error")

	suite.mockSql.ExpectQuery("SELECT").WithArgs(senderID).WillReturnError(expectedError)
	repository := NewTxRepository(suite.mockDB)

	txs, err := repository.GetBySenderId(senderID)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "error while getting transactions for sender 1: dummy error", err.Error())

	assert.Nil(suite.T(), txs)
}
func (suite *TransactionRepositoryTestSuite) TestGetBySenderId_ScanNilValues() {
	// Prepare mock rows
	senderID := uint(1)
	rows := sqlmock.NewRows([]string{"transaction_type", "sender_id", "recipient_id", "bank_account_id", "card_id", "point_exchange_id", "amount", "point", "transaction_date"}).
		AddRow("debit", senderID, nil, nil, nil, nil, nil, nil, now)

	// Set up expectations
	suite.mockSql.ExpectQuery("SELECT").WithArgs(senderID).WillReturnRows(rows)
	repository := NewTxRepository(suite.mockDB)

	// Call the function
	txs, err := repository.GetBySenderId(senderID)

	// Check the results
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), txs)

	assert.Len(suite.T(), txs, 1)
	assert.Equal(suite.T(), "debit", txs[0].TransactionType)
	assert.Equal(suite.T(), senderID, txs[0].SenderID)
	assert.Nil(suite.T(), nil)
	assert.Nil(suite.T(), nil)
	assert.Nil(suite.T(), nil)
	assert.Nil(suite.T(), nil)
	assert.Nil(suite.T(), nil)
	assert.Nil(suite.T(), nil)
	assert.Equal(suite.T(), now, txs[0].TransactionDate)

}

func (suite *TransactionRepositoryTestSuite) TestCreateDepositBank_Success() {
	// Create a new transaction object with test data
	tx := &model.TransactionBank{
		SenderID:      1,
		BankAccountID: 1,
		Amount:        50000,
	}

	// Create a mock database connection and repository

	repo := NewTxRepository(suite.mockDB)

	// Expect the query to be executed with the correct arguments
	suite.mockSql.ExpectExec("INSERT INTO tx_transaction").WithArgs("Deposit Bank", tx.SenderID, tx.BankAccountID, tx.Amount, now).WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the CreateDepositBank method
	err := repo.CreateDepositBank(tx)

	// Assert that no errors occurred and all expectations were met
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mockSql.ExpectationsWereMet())
}
func (suite *TransactionRepositoryTestSuite) TestCreateDepositBank_Error() {
	tx := dummyTxBank[0]

	expectedErr := errors.New("database error")

	suite.mockSql.ExpectExec("INSERT INTO tx_transaction").WillReturnError(expectedErr)

	repository := NewTxRepository(suite.mockDB)

	err := repository.CreateDepositBank(tx)

	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *TransactionRepositoryTestSuite) TestCreateDepositCard_Success() {
	// Create a new transaction object with test data
	tx := &model.TransactionCard{
		SenderID: 1,
		CardID:   1,
		Amount:   50000,
	}

	// Create a mock database connection and repository

	repo := NewTxRepository(suite.mockDB)

	// Expect the query to be executed with the correct arguments
	suite.mockSql.ExpectExec("INSERT INTO tx_transaction").WithArgs("Deposit Card", tx.SenderID, tx.CardID, tx.Amount, now).WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the CreateDepositBank method
	err := repo.CreateDepositCard(tx)

	// Assert that no errors occurred and all expectations were met
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mockSql.ExpectationsWereMet())
}
func (suite *TransactionRepositoryTestSuite) TestCreateDepositCard_Error() {
	tx := dummyTxCard[0]

	expectedErr := errors.New("database error")

	suite.mockSql.ExpectExec("INSERT INTO tx_transaction").WillReturnError(expectedErr)

	repository := NewTxRepository(suite.mockDB)

	err := repository.CreateDepositCard(tx)

	assert.Equal(suite.T(), expectedErr, err)
}
func (suite *TransactionRepositoryTestSuite) TestCreateWitdrawal_Success() {
	// Create a new transaction object with test data
	tx := &model.TransactionWithdraw{
		SenderID:      1,
		BankAccountID: 1,
		Amount:        50000,
	}

	repo := NewTxRepository(suite.mockDB)

	// Expect the query to be executed with the correct arguments
	suite.mockSql.ExpectExec("INSERT INTO tx_transaction").WithArgs("Withdraw", tx.SenderID, tx.BankAccountID, tx.Amount, now).WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the CreateDepositBank method
	err := repo.CreateWithdrawal(tx)

	// Assert that no errors occurred and all expectations were met
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mockSql.ExpectationsWereMet())
}
func (suite *TransactionRepositoryTestSuite) TestCreateWithdrawal_Error() {
	tx := dummyTxWd[0]

	expectedErr := errors.New("database error")

	suite.mockSql.ExpectExec("INSERT INTO tx_transaction").WillReturnError(expectedErr)

	repository := NewTxRepository(suite.mockDB)

	err := repository.CreateWithdrawal(tx)

	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *TransactionRepositoryTestSuite) TestCreateRedeem_Success() {
	// Create some test data
	senderID := txP[0].SenderID
	redem := txP[0]

	repo := NewTxRepository(suite.mockDB)

	// Expect the query to be executed with the correct arguments
	suite.mockSql.ExpectExec("INSERT INTO tx_transaction").WithArgs(txP[0].TransactionType, senderID, txP[0].PointExchangeID, txP[0].Point, txP[0].TransactionDate).WillReturnResult(
		sqlmock.NewResult(1, 1))

	// Call the GetBySenderId method
	err := repo.CreateRedeem(redem)

	// Assert that no errors occurred and all expectations were met
	assert.NoError(suite.T(), err)

	assert.NoError(suite.T(), suite.mockSql.ExpectationsWereMet())
}
func (suite *TransactionRepositoryTestSuite) TestCreateRedeem_Error() {
	tx := dummyTxRd[0]

	expectedErr := errors.New("database error")

	suite.mockSql.ExpectExec("INSERT INTO tx_transaction").WillReturnError(expectedErr)

	repository := NewTxRepository(suite.mockDB)

	err := repository.CreateRedeem(tx)

	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *TransactionRepositoryTestSuite) TestCreateTransfer_Success() {
	// Create some test data
	tx := &model.TransactionTransfer{
		TransactionType: "Transfer",
		SenderID:        1,
		RecipientID:     2,
		Amount:          50000,
		TransactionDate: now,
	}

	repo := NewTxRepository(suite.mockDB)

	// Expect the query to be executed with the correct arguments
	suite.mockSql.ExpectExec("INSERT INTO tx_transaction").WithArgs(tx.TransactionType, tx.SenderID, tx.RecipientID, tx.Amount, tx.TransactionDate).WillReturnResult(
		sqlmock.NewResult(1, 1))

	// Call the GetBySenderId method
	result, err := repo.CreateTransfer(tx)

	// Assert that no errors occurred and all expectations were met
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)

	assert.NoError(suite.T(), suite.mockSql.ExpectationsWereMet())
}
func (suite *TransactionRepositoryTestSuite) TestCreateTransfer_Error() {
	tx := dummyTxTf[0]

	expectedErr := fmt.Errorf("failed to create data: database error")

	suite.mockSql.ExpectExec("INSERT INTO tx_transaction").WillReturnError(errors.New("database error"))

	repository := NewTxRepository(suite.mockDB)

	_, err := repository.CreateTransfer(tx)

	assert.Equal(suite.T(), expectedErr, err)
}
func (suite *TransactionRepositoryTestSuite) TestGetAllPoint_Success() {
	// Create some test data
	pointExchanges := []*model.PointExchange{
		{
			PE_ID:  1,
			Reward: "10k Pulsa",
			Price:  100,
		},
		{
			PE_ID:  2,
			Reward: "20k Pulsa",
			Price:  200,
		},
	}

	// Create a mock database connection and repository

	repo := NewTxRepository(suite.mockDB)

	// Expect the query to be executed and return some test data
	rows := sqlmock.NewRows([]string{"pe_id", "reward", "price"}).
		AddRow(pointExchanges[0].PE_ID, pointExchanges[0].Reward, pointExchanges[0].Price).
		AddRow(pointExchanges[1].PE_ID, pointExchanges[1].Reward, pointExchanges[1].Price)

	suite.mockSql.ExpectQuery("SELECT pe_id, reward, price FROM mst_point_exchange").WillReturnRows(rows)

	// Call the GetAllPoint method
	results, err := repo.GetAllPoint()

	// Assert that no errors occurred and all expectations were met
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pointExchanges, results)
	assert.NotNil(suite.T(), results)
	assert.NotNil(suite.T(), rows)

	assert.NoError(suite.T(), suite.mockSql.ExpectationsWereMet())
}
func (suite *TransactionRepositoryTestSuite) TestGetAllPoint_Error() {
	expectedErr := fmt.Errorf("failed to get data: database error")

	suite.mockSql.ExpectQuery("SELECT pe_id, reward, price FROM mst_point_exchange").WillReturnError(errors.New("database error"))

	repository := NewTxRepository(suite.mockDB)

	_, err := repository.GetAllPoint()

	assert.Equal(suite.T(), expectedErr, err)
}

func (suite *TransactionRepositoryTestSuite) TestGetByPeId_Success() {
	// Create a mock database connection and repository

	repo := NewTxRepository(suite.mockDB)

	// Create some test data
	peID := uint(1)
	reward := "10k"
	price := 100

	// Expect the query to be executed with the correct arguments
	suite.mockSql.ExpectQuery("SELECT pe_id, reward, price FROM mst_point_exchange WHERE pe_id = ?").WithArgs(peID).WillReturnRows(
		sqlmock.NewRows([]string{"pe_id", "reward", "price"}).AddRow(peID, reward, price))

	// Call the GetByPeId method
	pointExchanges, err := repo.GetByPeId(peID)

	// Assert that no errors occurred and all expectations were met
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), pointExchanges)

	assert.NoError(suite.T(), suite.mockSql.ExpectationsWereMet())
}
func (suite *TransactionRepositoryTestSuite) TestGetByPeId_Error() {
	expectedErr := errors.New("database error")

	row := suite.mockSql.ExpectQuery("SELECT pe_id, reward, price FROM mst_point_exchange").WillReturnError(expectedErr)

	repository := NewTxRepository(suite.mockDB)

	_, err := repository.GetByPeId(1)

	assert.Equal(suite.T(), expectedErr, err)
	assert.NotNil(suite.T(), row)

}

func (suite *TransactionRepositoryTestSuite) SetupTest() {
	mockDb, mockSql, err := sqlmock.New()
	if err != nil {
		log.Fatalln("Error database", err)
	}
	suite.mockDB = mockDb
	suite.mockSql = mockSql
}
func (suite *TransactionRepositoryTestSuite) TearDownTest() {
	suite.mockDB.Close()
}

func TestTxRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionRepositoryTestSuite))
}
