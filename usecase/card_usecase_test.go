package usecase

import (
	"errors"
	"testing"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyCard = []model.Card{
	{
		CardID:         1,
		UserID:         1,
		CardType:       "BRI",
		CardNumber:     "1234-5678-9101-3456",
		ExpirationDate: "07/25",
		CVV:            "123",
	},
	{
		CardID:         2,
		UserID:         2,
		CardType:       "BCA",
		CardNumber:     "4321-8765-9101-3456",
		ExpirationDate: "03/26",
		CVV:            "321",
	},
}
var dummyCardResponse = []model.CardResponse{
	{
		UserID:         1,
		CardType:       "BRI",
		CardNumber:     "1234-5678-9101-3456",
		ExpirationDate: "07/25",
		CVV:            "123",
	},
	{
		UserID:         2,
		CardType:       "BCA",
		CardNumber:     "4321-8765-9101-3456",
		ExpirationDate: "03/26",
		CVV:            "321",
	},
}
var dummyCardResponse1 = []*model.CardResponse{
	{
		UserID:         1,
		CardType:       "BRI",
		CardNumber:     "1234-5678-9101-3456",
		ExpirationDate: "07/25",
		CVV:            "123",
	},
	{
		UserID:         2,
		CardType:       "BCA",
		CardNumber:     "4321-8765-9101-3456",
		ExpirationDate: "03/26",
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

type cardRepoMock struct {
	mock.Mock
}

func (u *cardRepoMock) GetAll() any {
	args := u.Called()
	if args.Get(0) == nil {
		return nil
	}
	return dummyCard
}

func (u *cardRepoMock) GetByUserID(id uint) ([]*model.CardResponse, error) {
	args := u.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.CardResponse), nil
}

func (u *cardRepoMock) GetByCardID(id uint) (*model.Card, error) {
	args := u.Called(id)
	if args.Get(0) == nil {
		return nil, errors.New("card account not found")
	}
	return args.Get(0).(*model.Card), nil
}

func (u *cardRepoMock) Create(id uint, newCard *model.CreateCard) (any, error) {
	args := u.Called(id, newCard)
	if args.Get(0) == nil {
		return nil, errors.New("failed to create data")
	}
	return dummyCardResponse, nil
}

func (u *cardRepoMock) Update(card *model.Card) string {
	args := u.Called(card)
	if args.Get(0) == nil {
		return "failed to update Card Account"
	}
	return "Card Account updated Successfully"
}

func (u *cardRepoMock) DeleteByUserID(id uint) string {
	args := u.Called(id)
	if args.Get(0) == nil {
		return "failed to delete Card Account"
	}
	return "All Card Account deleted Successfully"
}

func (u *cardRepoMock) DeleteByCardID(cardID uint) error {
	args := u.Called(cardID)
	if args.Get(0) != nil {
		return errors.New("failed to delete data")
	}
	return nil
}

type CardUsecaseTestSuite struct {
	cardaccRepoMock *cardRepoMock
	suite.Suite
}

func (suite *CardUsecaseTestSuite) TestFindAllCard_Success() {
	cardUsecase := NewCardUsecase(suite.cardaccRepoMock)
	suite.cardaccRepoMock.On("GetAll").Return(dummyCard)
	res := cardUsecase.FindAllCard()
	assert.NotNil(suite.T(), res)
	assert.Equal(suite.T(), dummyCard, res)
}
func (suite *CardUsecaseTestSuite) TestFindAllCard_Failed() {
	cardUsecase := NewCardUsecase(suite.cardaccRepoMock)
	suite.cardaccRepoMock.On("GetAll").Return(nil)
	res := cardUsecase.FindAllCard()
	assert.Nil(suite.T(), res)
	assert.Empty(suite.T(), res)
}

func (suite *CardUsecaseTestSuite) TestFindCardByUserID_Success() {
	userID := uint(1)
	cardAcc := dummyCardResponse1
	cardUsecase := NewCardUsecase(suite.cardaccRepoMock)
	suite.cardaccRepoMock.On("GetByUserID", userID).Return(cardAcc, nil)
	result, err := cardUsecase.FindCardByUserID(userID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), cardAcc, result)
}
func (suite *CardUsecaseTestSuite) TestFindCardByUserID_Failed() {
	userID := uint(1)
	expectedError := errors.New("failed to get card account")
	cardUsecase := NewCardUsecase(suite.cardaccRepoMock)
	suite.cardaccRepoMock.On("GetByUserID", userID).Return(nil, expectedError)
	result, err := cardUsecase.FindCardByUserID(userID)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *CardUsecaseTestSuite) TestFindCardByCardID_Success() {
	userID := uint(1)
	cardAcc := &dummyCard[0]
	cardUsecase := NewCardUsecase(suite.cardaccRepoMock)
	suite.cardaccRepoMock.On("GetByCardID", userID).Return(cardAcc, nil)
	result, err := cardUsecase.FindCardByCardID(userID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), cardAcc, result)
}
func (suite *CardUsecaseTestSuite) TestFindCardByCardID_Failed() {
	userID := uint(1)
	expectedError := errors.New("card account not found")
	cardUsecase := NewCardUsecase(suite.cardaccRepoMock)
	suite.cardaccRepoMock.On("GetByCardID", userID).Return(nil, expectedError)
	result, err := cardUsecase.FindCardByCardID(userID)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *CardUsecaseTestSuite) TestRegister_Success() {
	userID := uint(1)
	cardUsecaseMock := NewCardUsecase(suite.cardaccRepoMock)
	suite.cardaccRepoMock.On("Create", userID, &dummyCardResponse[0]).Return(dummyCardResponse, nil)
	result, err := cardUsecaseMock.Register(userID, &dummyCreateCard[0])
	assert.NotNil(suite.T(), result)
	assert.Nil(suite.T(), err)
}
func (suite *CardUsecaseTestSuite) TestRegister_Failed() {
	userID := uint(1)
	expectedError := errors.New("failed to create data")
	cardUsecaseMock := NewCardUsecase(suite.cardaccRepoMock)
	suite.cardaccRepoMock.On("Create", userID, &dummyCardResponse[0]).Return(nil, expectedError)
	result, err := cardUsecaseMock.Register(userID, &dummyCreateCard[0])
	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), err)
}

func (suite *CardUsecaseTestSuite) TestEdit_Success() {
	cardAcc := &dummyCard[0]
	cardUsecase := NewCardUsecase(suite.cardaccRepoMock)
	suite.cardaccRepoMock.On("Update", cardAcc).Return("Card Account updated Successfully")
	result := cardUsecase.Edit(cardAcc)
	assert.NotNil(suite.T(), result)
}
func (suite *CardUsecaseTestSuite) TestEdit_Failed() {
	cardAcc := &dummyCard[0]
	cardUsecase := NewCardUsecase(suite.cardaccRepoMock)
	suite.cardaccRepoMock.On("Update", cardAcc).Return("Failed to Update Card Account")
	err := cardUsecase.Edit(cardAcc)
	assert.NotNil(suite.T(), err)
}

func (suite *CardUsecaseTestSuite) TestUnregAll_Success() {
	userID := uint(1)
	cardUsecase := NewCardUsecase(suite.cardaccRepoMock)
	suite.cardaccRepoMock.On("DeleteByUserID", userID).Return("All Card Account deleted Successfully")
	result := cardUsecase.UnregALL(userID)
	assert.NotNil(suite.T(), result)
}
func (suite *CardUsecaseTestSuite) TestUnregAll_Failed() {
	userID := uint(1)
	cardUsecase := NewCardUsecase(suite.cardaccRepoMock)
	suite.cardaccRepoMock.On("DeleteByUserID", userID).Return("Failed to delete Card Account")
	err := cardUsecase.UnregALL(userID)
	assert.NotNil(suite.T(), err)
}

func (suite *CardUsecaseTestSuite) TestUnregByCardID_Success() {
	cardID := uint(1)
	cardUsecase := NewCardUsecase(suite.cardaccRepoMock)
	suite.cardaccRepoMock.On("DeleteByCardID", cardID).Return(nil)
	err := cardUsecase.UnregByCardID(cardID)
	assert.NoError(suite.T(), err)
}
func (suite *CardUsecaseTestSuite) TestUnregByCardID_Failed() {
	cardID := uint(1)
	cardUsecase := NewCardUsecase(suite.cardaccRepoMock)
	suite.cardaccRepoMock.On("DeleteByCardID", cardID).Return(errors.New("failed to delete data"))
	err := cardUsecase.UnregByCardID(cardID)
	assert.EqualError(suite.T(), err, "failed to delete data")
}

func (suite *CardUsecaseTestSuite) SetupTest() {
	suite.cardaccRepoMock = new(cardRepoMock)
}

func TestCardAccTestSuite(t *testing.T) {
	suite.Run(t, new(CardUsecaseTestSuite))
}
