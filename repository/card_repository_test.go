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

var dummyCard = []model.Card{
	{
		CardID:         1,
		UserID:         1,
		CardType:       "BRI",
		CardNumber:     "123456789",
		ExpirationDate: "08/26",
		CVV:            "123",
	},
	{
		CardID:         2,
		UserID:         1,
		CardType:       "BCA",
		CardNumber:     "123456345",
		ExpirationDate: "02/27",
		CVV:            "012",
	},
	{
		CardID:         3,
		UserID:         2,
		CardType:       "Mandiri",
		CardNumber:     "987654321",
		ExpirationDate: "04/25",
		CVV:            "321",
	},
}

var dummyCardResponse = []model.CardResponse{
	{
		UserID:         1,
		CardID:         1,
		CardType:       "BRI",
		CardNumber:     "123456789",
		ExpirationDate: "08/26",
		CVV:            "123",
	},
	{
		UserID:         1,
		CardID:         2,
		CardType:       "BCA",
		CardNumber:     "123456345",
		ExpirationDate: "02/27",
		CVV:            "012",
	},
	{
		UserID:         2,
		CardID:         3,
		CardType:       "Mandiri",
		CardNumber:     "987654321",
		ExpirationDate: "04/25",
		CVV:            "321",
	},
}

var dummyCreateCard = []model.CreateCard{
	{
		UserID:         1,
		CardType:       "BRI",
		CardNumber:     "123456789",
		ExpirationDate: "08/26",
		CVV:            "123",
	},
	{
		UserID:         1,
		CardType:       "BCA",
		CardNumber:     "123456345",
		ExpirationDate: "02/27",
		CVV:            "012",
	},
	{
		UserID:         2,
		CardType:       "Mandiri",
		CardNumber:     "987654321",
		ExpirationDate: "04/25",
		CVV:            "321",
	},
}

type CardRepositoryTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSql sqlmock.Sqlmock
}

func (suite *CardRepositoryTestSuite) TestGetAll_Success() {
	var users = dummyCardResponse[:3]
	suite.mockSql.ExpectQuery("SELECT user_id, card_id, card_type, card_number, expiration_date, cvv FROM mst_card").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "card_id", "card_type", "card_number", "expiration_date", "CVV"}).
			AddRow(users[0].UserID, users[0].CardID, users[0].CardType, users[0].CardNumber, users[0].ExpirationDate, users[0].CVV).
			AddRow(users[1].UserID, users[1].CardID, users[1].CardType, users[1].CardNumber, users[1].ExpirationDate, users[1].CVV).
			AddRow(users[2].UserID, users[2].CardID, users[2].CardType, users[2].CardNumber, users[2].ExpirationDate, users[2].CVV))
	cardRepository := NewCardRepository(suite.mockDB)
	result := cardRepository.GetAll()
	assert.NotNil(suite.T(), result)
}

func (suite *CardRepositoryTestSuite) TestGetAllNilRows_Failed() {
	suite.mockSql.ExpectQuery("SELECT user_id, card_id, card_type, card_number, expiration_date, cvv FROM mst_card").
		WillReturnError(fmt.Errorf("no data"))
	cardRepository := NewCardRepository(suite.mockDB)
	err := cardRepository.GetAll()
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "no data", err)
}

func (suite *CardRepositoryTestSuite) TestGetAllScanRows_Failed() {
	var users = dummyCardResponse[:3]
	suite.mockSql.ExpectQuery("SELECT user_id, card_id, card_type, card_number, expiration_date, cvv FROM mst_card").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "card_id", "card_type", "card_number", "expiration_date", "CVV"}).
			AddRow(nil, users[0].CardID, users[0].CardType, users[0].CardNumber, users[0].ExpirationDate, users[0].CVV).
			AddRow(nil, users[1].CardID, users[1].CardType, users[1].CardNumber, users[1].ExpirationDate, users[1].CVV).
			AddRow(nil, users[2].CardID, users[2].CardType, users[2].CardNumber, users[2].ExpirationDate, users[2].CVV))
	cardRepository := NewCardRepository(suite.mockDB)
	err := cardRepository.GetAll()
	assert.NotNil(suite.T(), err)
}

func (suite *CardRepositoryTestSuite) TestGetByUserID_Success() {
	cards := dummyCardResponse[:2]
	suite.mockSql.ExpectQuery("SELECT user_id, card_id, card_type, card_number, expiration_date, cvv FROM mst_card WHERE user_id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "card_id", "card_type", "card_number", "expiration_date", "cvv"}).
			AddRow(cards[0].UserID, cards[0].CardID, cards[0].CardType, cards[0].CardNumber, cards[0].ExpirationDate, cards[0].CVV).
			AddRow(cards[1].UserID, cards[1].CardID, cards[1].CardType, cards[1].CardNumber, cards[1].ExpirationDate, cards[1].CVV))
	cardRepository := NewCardRepository(suite.mockDB)
	result, err := cardRepository.GetByUserID(1)
	assert.NotNil(suite.T(), result)
	assert.Nil(suite.T(), err)
}

func (suite *CardRepositoryTestSuite) TestGetByUserIDScanRows_Failed() {
	cards := dummyCardResponse[:2]
	suite.mockSql.ExpectQuery("SELECT user_id, card_id, card_type, card_number, expiration_date, cvv FROM mst_card WHERE user_id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "card_id", "card_type", "card_number", "expiration_date", "cvv"}).
			AddRow(nil, cards[0].CardID, cards[0].CardType, cards[0].CardNumber, cards[0].ExpirationDate, cards[0].CVV).
			AddRow(cards[1].UserID, cards[1].CardID, cards[1].CardType, cards[1].CardNumber, cards[1].ExpirationDate, cards[1].CVV))
	cardRepository := NewCardRepository(suite.mockDB)
	result, err := cardRepository.GetByUserID(1)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *CardRepositoryTestSuite) TestGetByCardID_Success() {
	card := dummyCard[0]
	suite.mockSql.ExpectQuery("SELECT card_id, card_type, card_number, expiration_date, cvv, user_id FROM mst_card WHERE card_id = \\$1").
		WithArgs(card.CardID).
		WillReturnRows(sqlmock.NewRows([]string{"card_id", "card_type", "card_number", "expiration_date", "cvv", "user_id"}).
			AddRow(card.CardID, card.CardType, card.CardNumber, card.ExpirationDate, card.CVV, card.UserID))
	cardRepository := NewCardRepository(suite.mockDB)
	result, err := cardRepository.GetByCardID(card.CardID)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), card, *result)
}

func (suite *CardRepositoryTestSuite) TestGetByCardID_Failed() {
	card := dummyCard[0]
	suite.mockSql.ExpectQuery("SELECT card_id, card_type, card_number, expiration_date, cvv, user_id FROM mst_card WHERE card_id = \\$1").
		WithArgs(card.CardID).
		WillReturnError(errors.New("card not found"))
	cardRepository := NewCardRepository(suite.mockDB)
	result, err := cardRepository.GetByCardID(card.CardID)
	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "card not found", err.Error())
}

func (suite *CardRepositoryTestSuite) TestCreateCard_Success() {
	newCard := dummyCreateCard[0]
	userID := uint(1)
	suite.mockSql.ExpectExec("INSERT INTO mst_card \\(user_id, card_type, card_number, expiration_date, cvv\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\)").
		WithArgs(userID, newCard.CardType, newCard.CardNumber, newCard.ExpirationDate, newCard.CVV).
		WillReturnResult(sqlmock.NewResult(1, 1))
	cardRepository := NewCardRepository(suite.mockDB)
	result, err := cardRepository.Create(userID, &newCard)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.EqualValues(suite.T(), &newCard, result)
}

func (suite *CardRepositoryTestSuite) TestCreate_Failed() {
	newCard := dummyCreateCard[0]
	userID := uint(1)
	suite.mockSql.ExpectExec("INSERT INTO mst_card \\(user_id, card_type, card_number, expiration_date, cvv\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\)").
		WithArgs(userID, newCard.CardType, newCard.CardNumber, newCard.ExpirationDate, newCard.CVV).WillReturnError(errors.New("failed to create data"))
	cardRepository := NewCardRepository(suite.mockDB)
	result, err := cardRepository.Create(userID, &newCard)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *CardRepositoryTestSuite) TestUpdate_Success() {
	card := dummyCard[0]
	suite.mockSql.ExpectQuery("SELECT user_id, card_id, card_type, card_number, expiration_date, cvv FROM mst_card WHERE user_id = \\$1").
		WithArgs(card.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "card_id", "card_type", "card_number", "expiration_date", "cvv"}))
	suite.mockSql.ExpectExec("UPDATE mst_card SET card_type = \\$1, card_number = \\$2, expiration_date = \\$3, cvv = \\$4 WHERE card_id = \\$5").
		WithArgs(card.CardType, card.CardNumber, card.ExpirationDate, card.CVV, card.CardID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	cardRepository := NewCardRepository(suite.mockDB)
	str := cardRepository.Update(&card)
	assert.NotNil(suite.T(), str)
}

func (suite *CardRepositoryTestSuite) TestUpdateScanID_Failed() {
	card := dummyCard[0]
	suite.mockSql.ExpectQuery("SELECT user_id, card_id, card_type, card_number, expiration_date, cvv FROM mst_card WHERE user_id = \\$1").
		WithArgs(card.UserID).
		WillReturnError(errors.New("user not found"))
	cardRepository := NewCardRepository(suite.mockDB)
	str := cardRepository.Update(&card)
	assert.NotNil(suite.T(), str)
}

func (suite *CardRepositoryTestSuite) TestUpdate_Failed() {
	card := dummyCard[0]
	expectedError := fmt.Errorf("failed to update Card ID")
	suite.mockSql.ExpectQuery("SELECT user_id, card_id, card_type, card_number, expiration_date, cvv FROM mst_card WHERE user_id = \\$1").
		WithArgs(card.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "card_id", "card_type", "card_number", "expiration_date", "cvv"}))
	suite.mockSql.ExpectExec("UPDATE mst_card SET card_type = \\$1, card_number = \\$2, expiration_date = \\$3, cvv = \\$4 WHERE card_id = \\$5").
		WithArgs(card.CardType, card.CardNumber, card.ExpirationDate, card.CVV, card.CardID).WillReturnError(expectedError)
	cardRepository := NewCardRepository(suite.mockDB)
	str := cardRepository.Update(&card)
	assert.NotNil(suite.T(), str)
}

func (suite *CardRepositoryTestSuite) TestDeleteByUserID_Success() {
	userID := uint(1)
	suite.mockSql.ExpectExec("DELETE FROM mst_card WHERE user_id = \\$1").WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	cardRepository := NewCardRepository(suite.mockDB)
	str := cardRepository.DeleteByUserID(userID)
	assert.EqualValues(suite.T(), "Deleted All Card ID Successfully", str)
}

func (suite *CardRepositoryTestSuite) TestDeleteByUserID_Failed() {
	userID := uint(1)
	expectedError := fmt.Errorf("failed to delete card")
	suite.mockSql.ExpectExec("DELETE FROM mst_card WHERE user_id = $1").
		WithArgs(userID).
		WillReturnError(expectedError)
	cardRepository := NewCardRepository(suite.mockDB)
	str := cardRepository.DeleteByUserID(userID)
	assert.NotNil(suite.T(), str)
}

func (suite *CardRepositoryTestSuite) TestDeleteByCardID_Success() {
	card := dummyCard[0]
	suite.mockSql.ExpectQuery("SELECT card_id, card_type, card_number, expiration_date, cvv, user_id FROM mst_card WHERE card_id = \\$1").
		WithArgs(card.CardID).
		WillReturnRows(sqlmock.NewRows([]string{"card_id", "user_id", "card_type", "card_number", "expiration_date", "cvv"}).AddRow(card.CardID, card.CardType, card.CardNumber, card.ExpirationDate, card.CVV, card.UserID))
	cardID := uint(1)
	suite.mockSql.ExpectExec("DELETE FROM mst_card WHERE card_id = \\$1").
		WithArgs(cardID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	cardRepository := NewCardRepository(suite.mockDB)
	err := cardRepository.DeleteByCardID(card.CardID)
	assert.Nil(suite.T(), err)
}

func (suite *CardRepositoryTestSuite) TestDeleteCardID_Failed() {
	card := dummyCard[0]
	suite.mockSql.ExpectQuery("SELECT card_id, card_type, card_number, expiration_date, cvv, user_id FROM mst_card WHERE card_id = \\$1").
		WithArgs(card.CardID).
		WillReturnRows(sqlmock.NewRows([]string{"card_id", "user_id", "card_type", "card_number", "expiration_date", "cvv"}).AddRow(card.CardID, card.CardType, card.CardNumber, card.ExpirationDate, card.CVV, card.UserID))
	cardID := uint(1)
	suite.mockSql.ExpectExec("DELETE FROM mst_card WHERE card_id = \\$1").
		WithArgs(cardID).
		WillReturnError(fmt.Errorf("failed to delete card with ID %d", cardID))
	cardRepository := NewCardRepository(suite.mockDB)
	err := cardRepository.DeleteByCardID(card.CardID)
	assert.NotNil(suite.T(), err)
}

func (suite *CardRepositoryTestSuite) TestDeleteByCardIDScan_Failed() {
	cardID := uint(4)
	suite.mockSql.ExpectQuery("SELECT card_id, card_type, card_number, expiration_date, cvv, user_id FROM mst_card WHERE card_id = \\$1").
		WithArgs(cardID).
		WillReturnError(sql.ErrNoRows)
	cardRepository := NewCardRepository(suite.mockDB)
	err := cardRepository.DeleteByCardID(cardID)
	assert.NotNil(suite.T(), err)
}

// Setup test
func (suite *CardRepositoryTestSuite) SetupTest() {
	mockDb, mockSql, err := sqlmock.New()
	if err != nil {
		log.Println("Error database", err)
	}
	suite.mockDB = mockDb
	suite.mockSql = mockSql
}

func (suite *CardRepositoryTestSuite) TearDownTest() {
	suite.mockDB.Close()
}

func TestCardRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CardRepositoryTestSuite))
}
