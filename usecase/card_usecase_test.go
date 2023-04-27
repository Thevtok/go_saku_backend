package usecase

import (
	"github.com/ReygaFitra/inc-final-project.git/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyCard = []model.Card{
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

type repoMock struct {
	mock.Mock
}

func (r *repoMock) GetAll() []model.Card {
	args := r.Called()
	return args.Get(0).([]model.Card)
}

func (r *repoMock) GetByUserID(id uint) ([]model.Card, error) {
	args := r.Called(id)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Card), nil
}

func (r *repoMock) GetByCardID(id uint) (*model.Card, error) {
	args := r.Called(id)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Card), nil
}

func (r *repoMock) Create(id uint, newCard *model.CardResponse) ([]model.Card, error) {
	args := r.Called(id, newCard)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Card), nil
}

func (r *repoMock) Update(card *model.Card) error {
	args := r.Called(card)
	return args.Error(0)
}

func (r *repoMock) DeleteByUserID(userID uint) string {
	args := r.Called(userID)
	return args.String(0)
}

func (r *repoMock) DeleteByCardID(cardID uint) error {
	args := r.Called(cardID)
	return args.Error(0)
}

type CardUseCaseTestSuite struct {
	repoMock *repoMock
	suite.Suite
}

func (suite *CardUseCaseTestSuite) SetupTest() {
	suite.repoMock = new(repoMock)
}

func (suite *CardUseCaseTestSuite) TestFindAllCard() {
	cards := []model.Card{dummyCard[0], dummyCard[1]}
	suite.repoMock.On("GetAll").Return(cards)
	cardUC := NewCardUsecase(suite.repoMock)
	res := cardUC.FindAllCard()
	assert.Equal(suite.T(), cards, res)
}

func (suite *CardUseCaseTestSuite) TestFindCardByUserID_Success() {
	userID := dummyCard[0].UserID
	cards := []model.Card{dummyCard[0]}
	suite.repoMock.On("GetByUserID", userID).Return(cards, nil)
	cardUC := NewCardUsecase(suite.cardRepoMock)
	res, err := cardUC.FindCardByUserID(userID)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), cards, res)
}
