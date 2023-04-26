package usecase

// import (
// 	"github.com/ReygaFitra/inc-final-project.git/model"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/suite"
// )

// var dummyCard = []model.Card{
// 	{
// 		UserID:         1,
// 		CardType:       "BRI",
// 		CardNumber:     "1234-5678-9101-3456",
// 		ExpirationDate: "07/25",
// 		CVV:            "123",
// 	},
// 	{
// 		UserID:         2,
// 		CardType:       "BCA",
// 		CardNumber:     "4321-8765-9101-3456",
// 		ExpirationDate: "03/26",
// 		CVV:            "321",
// 	},
// }

// type cardRepoMock struct {
// 	mock.Mock
// }

// func (m *cardRepoMock) GetAll() []model.Card {
// 	args := m.Called()
// 	return args.Get(0).([]model.Card)
// }

// func (m *cardRepoMock) FindCardByUserID(id uint) ([]model.Card, error) {
// 	args := m.Called(id)
// 	return args.Get(0).([]model.Card), args.Error(1)
// }

// func (m *cardRepoMock) FindCardByCardID(id uint) (*model.Card, error) {
// 	args := m.Called(id)
// 	return args.Get(0).(*model.Card), args.Error(1)
// }

// func (m *cardRepoMock) Create(id uint, newCard *model.CardResponse) (*model.Card, error) {
// 	args := m.Called(id, newCard)
// 	return args.Get(0).(*model.Card), args.Error(1)
// }

// func (m *cardRepoMock) Update(card *model.Card) error {
// 	args := m.Called(card)
// 	return args.Error(0)
// }

// func (m *cardRepoMock) DeleteByUserID(userID uint) error {
// 	args := m.Called(userID)
// 	return args.Error(0)
// }

// func (m *cardRepoMock) DeleteByCardID(cardID uint) error {
// 	args := m.Called(cardID)
// 	return args.Error(0)
// }
