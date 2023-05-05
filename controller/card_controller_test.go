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

var dummyCard = []model.Card{
	{
		CardID: 		1,
		UserID:         1,
		CardType:       "BRI",
		CardNumber:     "1234-5678-9101-3456",
		ExpirationDate: "07/25",
		CVV:            "123",
	},
	{
		CardID: 		2,
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

type CardUsecaseMock struct {
	mock.Mock
}

func setupRouterCard() *gin.Engine {
	r := gin.Default()
	return r
}

func (c *CardUsecaseMock) FindAllCard() any {
	args := c.Called()
	if args.Get(0) == nil {
		return nil
	}
	return dummyCard
}

func (c *CardUsecaseMock) FindCardByUserID(id uint) (any, error) {
	args := c.Called(id)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return dummyCardResponse1[0], nil
}

func (c *CardUsecaseMock) FindCardByCardID(id uint) (*model.Card, error) {
	args := c.Called(id)
	// if args.Get(0) == nil {
	// 	return nil, errors.New("card account not found")
	// }
	// return args.Get(0).(*model.Card), nil

	if args[0] == nil {
		return nil, args.Error(1)
	}
	return &dummyCard[0], nil
}

func (c *CardUsecaseMock) Register(id uint, newCard *model.CardResponse) (any, error) {
	args := c.Called(id, newCard)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return &dummyCardResponse[0], nil
}

func (c *CardUsecaseMock) Edit(card *model.Card) string {
	args := c.Called(card)
	if args.Get(0) == nil {
		return "failed to update Card Account"
	}
	return "Card Account updated Successfully"
}

func (c *CardUsecaseMock) UnregALL(card *model.Card) string {
	args := c.Called(card)
	if args[0] == nil {
		return "failed to delete Card Account"
	}
	return "All Card Account deleted Successfully"
}

func (c *CardUsecaseMock) UnregByCardID(cardID uint) error {
	args := c.Called(cardID)
	if args.Get(0) != nil {
		return errors.New("failed to delete card id")
	}
	return nil
}

type CardControllerTestSuite struct {
	suite.Suite
	routerMock  *gin.Engine
	usecaseMock *CardUsecaseMock
}

func(suite *CardControllerTestSuite) TestFindAllCard_Success() {
	Cards := &dummyCardResponse
  	controller := NewCardController(suite.usecaseMock)
	router := setupRouterCard()
	router.GET("/user/card", controller.FindAllCard)

	suite.usecaseMock.On("FindAllCard").Return(Cards)
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/card", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}
func(suite *CardControllerTestSuite) TestFindAllCard_Failed() {
  	controller := NewCardController(suite.usecaseMock)
	router := setupRouterCard()
	router.GET("/user/card", controller.FindAllCard)

	suite.usecaseMock.On("FindAllCard").Return(nil)
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/card", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}

func(suite *CardControllerTestSuite) TestFindCardByUserID_Success() {
	Card := dummyCardResponse1[0]
	controller := NewCardController(suite.usecaseMock)
	router := setupRouterCard()
	router.GET("/user/card/:user_id", controller.FindCardByUserID)

	suite.usecaseMock.On("FindCardByUserID", Card.UserID).Return(Card, nil)
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/card/1", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}
func (suite *CardControllerTestSuite) TestFindCardByUserID_InvalidUserID() {
	controller := NewCardController(suite.usecaseMock)
    router := setupRouterCard()
    router.GET("/user/card/:user_id", controller.FindCardByUserID)

    userID := "invalid_id"
    suite.usecaseMock.On("FindCardByUserID", mock.Anything).Return(nil, errors.New("Invalid user ID"))

    r := httptest.NewRecorder()
    request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user/card/%s", userID), nil)
    router.ServeHTTP(r, request)
    assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
    suite.usecaseMock.AssertExpectations(suite.T())
}
func (suite *CardControllerTestSuite) TestFindCardByUserID_UserNotFound() {
	controller := NewCardController(suite.usecaseMock)
	router := setupRouterCard()
	router.GET("/user/card/:user_id", controller.FindCardByUserID)

	suite.usecaseMock.On("FindCardByUserID", uint(1)).Return(nil, errors.New("user not found"))
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/card/1", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}

func(suite *CardControllerTestSuite) TestFindCardByCardID_Success() {
	Card := &dummyCard[0]
	controller := NewCardController(suite.usecaseMock)
	router := setupRouterCard()
	router.GET("/user/card/:user_id/:card_id", controller.FindCardByCardID)

	suite.usecaseMock.On("FindCardByCardID", Card.CardID).Return(Card, nil)
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/card/1/1", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}
func (suite *CardControllerTestSuite) TestFindCardByCardID_InvalidCardID() {
	controller := NewCardController(suite.usecaseMock)
	router := setupRouterCard()
	router.GET("/user/card/:user_id/:card_id", controller.FindCardByCardID)

    r := httptest.NewRecorder()
    request, _ := http.NewRequest(http.MethodGet, "/user/card/1/invalid_ID", nil)
    router.ServeHTTP(r, request)
    assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
    suite.usecaseMock.AssertExpectations(suite.T())
}
func (suite *CardControllerTestSuite) TestFindCardByCardID_UserNotFound() {
	controller := NewCardController(suite.usecaseMock)
	router := setupRouterCard()
	router.GET("/user/card/:user_id/:card_id", controller.FindCardByCardID)

	suite.usecaseMock.On("FindCardByCardID", uint(1)).Return(nil, errors.New("Card ID not found"))
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/card/1/1", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}

// func (suite *CardControllerTestSuite) TestCreateCardID_Success() {
// 	newCard := dummyCardResponse1[0]
// 	controller := NewCardController(suite.usecaseMock)
// 	router := setupRouterCard()
// 	router.POST("/user/card/add/:user_id", controller.CreateCardID)

// 	suite.usecaseMock.On("FindCardByUserID", uint(1)).Return(newCard, nil)
// 	suite.usecaseMock.On("Register", uint(1), newCard).Return(newCard, nil)
// 	r := httptest.NewRecorder()
//     reqBody, _ := json.Marshal(newCard)
// 	request, _ := http.NewRequest(http.MethodPost, "/user/card/add/1", bytes.NewBuffer(reqBody))
// 	router.ServeHTTP(r, request)
// 	response := r.Body.String()
// 	var actualCard *model.CardResponse
// 	json.Unmarshal([]byte(response), actualCard)
// 	assert.Equal(suite.T(), http.StatusCreated, r.Code)
// }

// func (suite *CardControllerTestSuite) TestCreateCardID_Success() {
// 	newCard := dummyCardResponse1[0]
//     controller := NewCardController(suite.usecaseMock)
//     router := setupRouterCard()
//     router.POST("/user/card/add/:user_id", controller.CreateCardID)

//     suite.usecaseMock.On("FindCardByUserID", uint(1)).Return(newCard, nil)
//     suite.usecaseMock.On("Register", uint(1), &newCard).Return(&newCard, nil)
//     r := httptest.NewRecorder()
//     reqBody, _ := json.Marshal(newCard)
//     request, _ := http.NewRequest(http.MethodPost, "/user/card/add/1", bytes.NewBuffer(reqBody))
//     token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
//     request.Header.Set("Authorization", "Bearer "+token)
//     router.ServeHTTP(r, request)
//     response := r.Body.Bytes()
//     var actualCard *model.CardResponse = &model.CardResponse{}
//     err := json.Unmarshal(response, actualCard)
//     if err != nil {
//         fmt.Println("Error:", err)
//     }
//     assert.Equal(suite.T(), http.StatusCreated, r.Code)
// }

// func (suite *CardControllerTestSuite) TestCreateCardID_Success() {
// 	newCard := dummyCardResponse1[0]
// 	controller := NewCardController(suite.usecaseMock)
// 	router := setupRouterCard()
// 	router.POST("/user/card/add/:user_id", controller.CreateCardID)

// 	suite.usecaseMock.On("FindCardByUserID", uint(1)).Return(newCard, nil)
// 	suite.usecaseMock.On("Register", uint(1), &newCard).Return(&newCard, nil)
// 	r := httptest.NewRecorder()
// 	reqBody, _ := json.Marshal(newCard)
// 	request, _ := http.NewRequest(http.MethodPost, "/user/card/add/1", bytes.NewBuffer(reqBody))
// 	token := "dummy_token"
// 	request.Header.Set("Authorization", "Bearer "+token)
// 	router.ServeHTTP(r, request)
// 	response := r.Body.Bytes()
// 	var actualCard *model.CardResponse = &model.CardResponse{}
// 	err := json.Unmarshal(response, actualCard)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 	}
// 	assert.Equal(suite.T(), http.StatusCreated, r.Code)
// }

func (suite *CardControllerTestSuite) TestEdit_Success() {
	card := dummyCard[0]
	controller := NewCardController(suite.usecaseMock)
	router := setupRouterCard()
	router.PUT("/user/card/update/:user_id/:card_id", controller.Edit)

	suite.usecaseMock.On("FindCardByCardID", uint(1)).Return(&card, nil)
    suite.usecaseMock.On("Edit", mock.Anything).Return("Card Account updated Successfully")
    r := httptest.NewRecorder()
    reqBody, _ := json.Marshal(card)
    request, _ := http.NewRequest(http.MethodPut, "/user/card/update/1/1", bytes.NewBuffer(reqBody))
    router.ServeHTTP(r, request)
    response := r.Body.String()
    var result model.Card
    json.Unmarshal([]byte(response), &result)
    assert.Equal(suite.T(), http.StatusOK, r.Code)
    suite.usecaseMock.AssertExpectations(suite.T())
}
func (suite *CardControllerTestSuite) TestEdit_CardNotFound() {
	card := dummyCard[0]
	controller := NewCardController(suite.usecaseMock)
	router := setupRouterCard()
	router.PUT("/user/card/update/:user_id/:card_id", controller.Edit)

	suite.usecaseMock.On("FindCardByCardID", uint(1)).Return(nil, errors.New("Card not found"))
    r := httptest.NewRecorder()
    reqBody, _ := json.Marshal(card)
    request, _ := http.NewRequest(http.MethodPut, "/user/card/update/1/1", bytes.NewBuffer(reqBody))
    router.ServeHTTP(r, request)
    response := r.Body.String()
    var result model.Card
    json.Unmarshal([]byte(response), &result)
    assert.Equal(suite.T(), http.StatusNotFound, r.Code)
    suite.usecaseMock.AssertExpectations(suite.T())
}
func (suite *CardControllerTestSuite) TestEdit_InvalidCardID() {
	card := dummyCard[0]
	controller := NewCardController(suite.usecaseMock)
	router := setupRouterCard()
	router.PUT("/user/card/update/:user_id/:card_id", controller.Edit)

    r := httptest.NewRecorder()
    reqBody, _ := json.Marshal(card)
    request, _ := http.NewRequest(http.MethodPut, "/user/card/update/1/invalid_id", bytes.NewBuffer(reqBody))
    router.ServeHTTP(r, request)
    response := r.Body.String()
    var result model.Card
    json.Unmarshal([]byte(response), &result)
    assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
    suite.usecaseMock.AssertExpectations(suite.T())
}
// func (suite *CardControllerTestSuite) TestEditBindJSON_Failed() {
// 	controller := NewCardController(suite.usecaseMock)
// 	router := setupRouterCard()
// 	router.PUT("/user/card/update/:user_id/:card_id", controller.Edit)

//     r := httptest.NewRecorder()
//     request, _ := http.NewRequest(http.MethodPut, "/user/card/update/1/1", nil)
//     router.ServeHTTP(r, request)
//     assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
//     suite.usecaseMock.AssertExpectations(suite.T())
// }


func (suite *CardControllerTestSuite) TestUnregAll_Success() {
	card := &model.Card{
		UserID: 1,
	}
	controller := NewCardController(suite.usecaseMock)
	router := setupRouterCard()
	router.DELETE("/user/card/:user_id", controller.UnregAll)

	suite.usecaseMock.On("UnregALL", card).Return("All Card Account deleted Successfully")
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, "/user/card/1", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}
func (suite *CardControllerTestSuite) TestUnregAll_InvalidUserID() {
	controller := NewCardController(suite.usecaseMock)
	router := setupRouterCard()
	router.DELETE("/user/card/:user_id", controller.UnregAll)

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, "/user/card/not_a_number", nil)
	suite.usecaseMock.On("UnregALL", mock.AnythingOfType("*model.Card")).Return("failed to delete Card Account").Times(0)

	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}

func (suite *CardControllerTestSuite) TestUnregByCardId_Success() {
	controller := NewCardController(suite.usecaseMock)
	router := setupRouterCard()
	router.DELETE("/user/card/:user_id/:card_id", controller.UnregByCardId)

	cardID := "1"
	userID := "1"
	url := fmt.Sprintf("/user/card/%s/%s", userID, cardID)

	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	resp := httptest.NewRecorder()

	suite.usecaseMock.On("UnregByCardID", uint(1)).Return(nil)

	router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.usecaseMock.AssertExpectations(suite.T())
}
func (suite *CardControllerTestSuite) TestUnregByCardId_InvalidCardID() {
	controller := NewCardController(suite.usecaseMock)
	router := setupRouterCard()
	router.DELETE("/user/card/:user_id/:card_id", controller.UnregByCardId)

	cardID := "invalid"
	userID := "1"
	url := fmt.Sprintf("/user/card/%s/%s", userID, cardID)

	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
	suite.usecaseMock.AssertNotCalled(suite.T(), "UnregByCardID")
}
func (suite *CardControllerTestSuite) TestUnregByCardId_Failed() {
	controller := NewCardController(suite.usecaseMock)
	router := setupRouterCard()
	router.DELETE("/user/card/:user_id/:card_id", controller.UnregByCardId)

    CardId := uint(1)
	id := "1"
    expectedErr := errors.New("failed to delete card ID")
    suite.usecaseMock.On("UnregByCardID", CardId).Return(expectedErr)

    r := httptest.NewRecorder()
    request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/user/card/1/%s", id), nil)
    router.ServeHTTP(r, request)
    assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)
    suite.usecaseMock.AssertExpectations(suite.T())
} 



func (suite *CardControllerTestSuite) SetupTest() {
	suite.routerMock = gin.Default()
	suite.usecaseMock = new(CardUsecaseMock)
}
func TestCardController(t *testing.T) {
	suite.Run(t, new(CardControllerTestSuite))
}