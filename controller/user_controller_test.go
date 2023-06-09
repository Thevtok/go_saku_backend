package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyUser = []model.User{
	{
		ID:           "1",
		Name:         "name1",
		Username:     "username1",
		Email:        "email1@mail.com",
		Password:     "password1",
		Phone_Number: "08111111",
		Address:      "address1",
		Balance:      100000,
		Role:         "user",
		Point:        10,
	},
	{
		ID:           "2",
		Name:         "name2",
		Username:     "username2",
		Email:        "email2@mail.com",
		Password:     "password2",
		Phone_Number: "08111111",
		Address:      "address2",
		Balance:      50000,
		Role:         "user",
		Point:        10,
	},
}

type UserUseCaseMock struct {
	mock.Mock
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func (r *UserUseCaseMock) Register(user *model.User) (any, error) {
	args := r.Called(user)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return &dummyUser[0], nil
}

func (r *UserUseCaseMock) Login(email string, password string, token string) (*model.User, error) {
	args := r.Called(email, password)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return &dummyUser[0], nil
}

func (r *UserUseCaseMock) FindUsers() any {
	args := r.Called()
	if args[0] == nil {
		return nil
	}
	return dummyUser
}

func (r *UserUseCaseMock) SaveDeviceToken(userId string, token string) error {
	args := r.Called(userId, token)
	if args[0] == nil {
		return nil
	}
	return nil
}

func (r *UserUseCaseMock) FindByUsername(username string) (*model.User, error) {
	args := r.Called(username)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return &dummyUser[0], nil
}
func (r *UserUseCaseMock) FindByPhone(username string) (*model.User, error) {
	args := r.Called(username)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return &dummyUser[0], nil
}

func (r *UserUseCaseMock) FindById(id string) (*model.User, error) {
	args := r.Called(id)
	if user, ok := args.Get(0).(*model.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}
func (r *UserUseCaseMock) FindByiDToken(id string) (*model.User, error) {
	args := r.Called(id)
	if user, ok := args.Get(0).(*model.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *UserUseCaseMock) EditProfile(user *model.User) string {
	args := r.Called(user)
	if args[0] == nil {
		return "Failed update profile"
	}
	return "Success update profile"
}

func (r *UserUseCaseMock) EditEmailPassword(user *model.User) string {
	args := r.Called(user)
	if args[0] == nil {
		return "Failed update email password"
	}
	return "Edit Successfully"
}

func (r *UserUseCaseMock) Unreg(user *model.User) string {
	args := r.Called(user)
	if args[0] == nil {
		return "Failed Delete user"
	}
	return "Success Delete user"
}

func (r *UserUseCaseMock) UpdateBalance(userID string, newBalance int) error {
	args := r.Called(userID, newBalance)
	if args[0] == nil {
		return args.Error(0)
	}
	return nil
}

func (r *UserUseCaseMock) UpdatePoint(userID string, newPoint int) error {
	args := r.Called(userID, newPoint)
	if args[0] == nil {
		return args.Error(0)
	}
	return nil
}

type UserControllerTestSuite struct {
	suite.Suite
	routerMock  *gin.Engine
	useCaseMock *UserUseCaseMock
	bankMock    *BankAccUsecaseMock

	photoMock *PhotoUseCaseMock
}

func (suite *UserControllerTestSuite) TestFindUsers_Success() {
	Users := &dummyUser
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.GET("/user", controller.FindUsers)

	suite.useCaseMock.On("FindUsers").Return(Users)
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	suite.useCaseMock.AssertExpectations(suite.T())
}
func (suite *UserControllerTestSuite) TestFindUsers_Failed() {
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.GET("/user", controller.FindUsers)

	suite.useCaseMock.On("FindUsers").Return(nil)
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
	suite.useCaseMock.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestFindByUsername_Success() {
	Users := &dummyUser[0]
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.GET("/user/:username", controller.FindUserByUsername)

	suite.useCaseMock.On("FindByUsername", Users.Username).Return(Users, nil)
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/username1", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	suite.useCaseMock.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestRegister_Success() {
	newUser := &model.User{
		Name:         "John Doe",
		Username:     "johndoe",
		Email:        "johndoe@gmail.com",
		Password:     "password1A",
		Phone_Number: "081234567890",
		Address:      "Jakarta",
	}
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.POST("/register", controller.Register)

	expectedResult := &model.User{
		Name:         "John Doe",
		Username:     "johndoe",
		Email:        "johndoe@gmail.com",
		Phone_Number: "081234567890",
		Address:      "Jakarta",
	}

	suite.useCaseMock.On("Register", newUser).Return(expectedResult, nil)

	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(newUser)
	request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusCreated, r.Code)
}
func (suite *UserControllerTestSuite) TestRegister_Failed() {
	newUser := &model.User{}
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.POST("/register", controller.Register)

	suite.useCaseMock.On("Register", newUser).Return(nil, errors.New("Invalid Input"))

	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(newUser)
	request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}
func (suite *UserControllerTestSuite) TestRegisterBindJSON_Failed() {
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.POST("/register", controller.Register)

	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/register", nil)
	router.ServeHTTP(r, request)
	response := r.Body.String()
	var actualUser model.User
	json.Unmarshal([]byte(response), &actualUser)
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}
func (suite *UserControllerTestSuite) TestRegister_EmailInvalid() {
	newUser := &model.User{
		Name:         "Dummy",
		Username:     "dummy123",
		Email:        "dummyemail.com",
		Password:     "Password123",
		Phone_Number: "081234567890",
		Address:      "Dummy Address",
	}
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.POST("/register", controller.Register)
	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(newUser)
	request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	response := r.Body.String()
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	assert.Contains(suite.T(), response, "email must be a gmail address")
}
func (suite *UserControllerTestSuite) TestRegister_InvalidPassword() {
	newUser := &model.User{
		Name:         "Dummy User",
		Username:     "dummyuser",
		Email:        "dummyuser@gmail.com",
		Password:     "1234567", // password kurang dari 8 karakter
		Phone_Number: "081234567890",
		Address:      "Jalan Dummy No. 123",
	}
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.POST("/register", controller.Register)

	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(newUser)
	request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	response := r.Body.String()
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	assert.Contains(suite.T(), response, "Invalid Input: password must have at least 8 characters")
}
func (suite *UserControllerTestSuite) TestRegister_InvalidPasswordUppercase() {
	newUser := &model.User{
		Name:         "Dummy User",
		Username:     "dummyuser",
		Email:        "dummyuser@gmail.com",
		Password:     "dummypassword",
		Phone_Number: "081234567890",
		Address:      "Jalan Dummy No. 123",
	}
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.POST("/register", controller.Register)

	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(newUser)
	request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	response := r.Body.String()
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	assert.Contains(suite.T(), response, "password must contain at least one uppercase letter and one number")
}
func (suite *UserControllerTestSuite) TestRegister_InvalidPhoneNumber() {
	newUser := &model.User{
		Name:         "Dummy User",
		Username:     "dummyuser",
		Email:        "dummyuser@gmail.com",
		Password:     "Password123",
		Phone_Number: "0812345",
		Address:      "Jalan Dummy No. 123",
	}
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.POST("/register", controller.Register)

	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(newUser)
	request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	response := r.Body.String()
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	assert.Contains(suite.T(), response, "phone_number must be 11 - 13 digit")
}

func (suite *UserControllerTestSuite) TestUnreg_Success() {
	user := &model.User{
		ID: "uint(1)",
	}
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.DELETE("/user/:user_id", controller.Unreg)
	suite.bankMock.On("UnregAll", uint(1)).Return("s")

	suite.photoMock.On("Remove", uint(1)).Return("s")
	suite.useCaseMock.On("Unreg", user).Return("Success Delete user")
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, "/user/1", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	suite.useCaseMock.AssertExpectations(suite.T())
}
func (suite *UserControllerTestSuite) TestFindByUsername_Failed() {
	Users := &dummyUser[0]
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.GET("/user/:username", controller.FindUserByUsername)

	suite.useCaseMock.On("FindByUsername", Users.Username).Return(nil, errors.New("User not found"))
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/username1", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
	suite.useCaseMock.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestEditProfile_Success() {
	userID := "1"
	updatedUser := &model.User{
		Name:         "Updated User",
		Username:     "updateduser",
		Email:        "updateduser@gmail.com",
		Phone_Number: "081234567890",
		Address:      "Jalan Updated No. 123",
	}
	mockUseCase := &UserUseCaseMock{}
	mockUseCase.On("FindById", uint(1)).Return(&model.User{}, nil)
	mockUseCase.On("EditProfile", updatedUser).Return(updatedUser)
	controller := NewUserController(mockUseCase, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.PUT("/user/profile/:user_id", controller.EditProfile)

	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(updatedUser)
	request, _ := http.NewRequest(http.MethodPut, "/user/profile/"+userID, bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	assert.Equal(suite.T(), http.StatusOK, r.Code)
}
func (suite *UserControllerTestSuite) TestEditProfile_UserNotFound() {
	user := &dummyUser[0]
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.PUT("/user/profile/:user_id", controller.EditProfile)

	suite.useCaseMock.On("FindById", uint(1)).Return(nil, errors.New("User not found"))
	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(user)
	request, _ := http.NewRequest(http.MethodPut, "/user/profile/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)
	response := r.Body.String()
	var result map[string]string
	json.Unmarshal([]byte(response), &result)
	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
	suite.useCaseMock.AssertExpectations(suite.T())
}
func (suite *UserControllerTestSuite) TestEditProfile_InvalidUserID() {
	user := &dummyUser[0]
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.PUT("/user/profile/:user_id", controller.EditProfile)

	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(user)
	request, _ := http.NewRequest(http.MethodPut, "/user/profile/invalid_id", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)
	response := r.Body.String()
	var result map[string]string
	json.Unmarshal([]byte(response), &result)
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	suite.useCaseMock.AssertExpectations(suite.T())
}
func (suite *UserControllerTestSuite) TestEditProfile_InvalidPhoneNumber() {
	userID := "1"
	updatedUser := &model.User{
		Phone_Number: "0812345",
	}
	mockUseCase := &UserUseCaseMock{}
	mockUseCase.On("FindById", uint(1)).Return(&model.User{}, nil)
	controller := NewUserController(mockUseCase, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.PUT("/user/profile/:user_id", controller.EditProfile)

	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(updatedUser)
	request, _ := http.NewRequest(http.MethodPut, "/user/profile/"+userID, bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)

	response := r.Body.String()
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	assert.Contains(suite.T(), response, "phone_number must be 11 - 13 digit")
}

func (suite *UserControllerTestSuite) TestEditProfile_InvalidInput() {
	user := &dummyUser[0]
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.PUT("/user/profile/:user_id", controller.EditProfile)

	suite.useCaseMock.On("FindById", uint(1)).Return(user, nil)
	r := httptest.NewRecorder()
	reqBody := []byte(`invalid request body`)
	request, _ := http.NewRequest(http.MethodPut, "/user/profile/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)
	response := r.Body.String()
	var result map[string]string
	json.Unmarshal([]byte(response), &result)
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	suite.useCaseMock.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestEditEmailPassword_Success() {
	user := &model.User{
		ID:       "1",
		Name:     "John",
		Email:    "john@gmail.com",
		Password: "Password123",
	}
	updatedUser := &model.User{
		ID:       "1",
		Name:     "John",
		Email:    "john@gmail.com",
		Password: "Password123",
	}
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.PUT("/user/pass/:user_id", controller.EditEmailPassword)

	suite.useCaseMock.On("FindById", uint(1)).Return(user, nil)
	suite.useCaseMock.On("EditEmailPassword", user).Return(updatedUser)
	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(user)
	request, _ := http.NewRequest(http.MethodPut, "/user/pass/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)
	response := r.Body.String()
	var result map[string]interface{}
	json.Unmarshal([]byte(response), &result)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	assert.Equal(suite.T(), true, result["status"])
	suite.useCaseMock.AssertExpectations(suite.T())
}
func (suite *UserControllerTestSuite) TestEditEmailPassword_UserNotFound() {
	user := &dummyUser[0]
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.PUT("/user/pass/:user_id", controller.EditEmailPassword)

	suite.useCaseMock.On("FindById", uint(1)).Return(nil, errors.New("User not found"))
	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(user)
	request, _ := http.NewRequest(http.MethodPut, "/user/pass/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)
	response := r.Body.String()
	var result map[string]string
	json.Unmarshal([]byte(response), &result)
	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
	suite.useCaseMock.AssertExpectations(suite.T())
}
func (suite *UserControllerTestSuite) TestEditEmailPassword_InvalidUserID() {
	user := &dummyUser[0]
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.PUT("/user/pass/:user_id", controller.EditEmailPassword)

	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(user)
	request, _ := http.NewRequest(http.MethodPut, "/user/pass/invalid_id", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)
	response := r.Body.String()
	var result map[string]string
	json.Unmarshal([]byte(response), &result)
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	suite.useCaseMock.AssertExpectations(suite.T())
}
func (suite *UserControllerTestSuite) TestEditEmailPassword_InvalidInput() {
	user := &dummyUser[0]
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.PUT("/user/pass/:user_id", controller.EditEmailPassword)

	suite.useCaseMock.On("FindById", uint(1)).Return(user, nil)
	r := httptest.NewRecorder()
	reqBody := []byte(`invalid request body`)
	request, _ := http.NewRequest(http.MethodPut, "/user/pass/1", bytes.NewBuffer(reqBody))
	router.ServeHTTP(r, request)
	response := r.Body.String()
	var result map[string]string
	json.Unmarshal([]byte(response), &result)
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	suite.useCaseMock.AssertExpectations(suite.T())
}
func (suite *UserControllerTestSuite) TestEditEmailPassword_InvalidEmail() {
	user := &model.User{
		ID:       "1",
		Name:     "John Doe",
		Email:    "johndoe@yahoo.com",
		Password: "password123",
	}
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.PUT("/user/pass/:user_id", controller.EditEmailPassword)

	suite.useCaseMock.On("FindById", uint(1)).Return(user, nil)

	reqBody, _ := json.Marshal(map[string]string{
		"email":    "johndoe@yahoo.com",
		"password": "newpassword",
	})
	request, _ := http.NewRequest(http.MethodPut, "/user/pass/1", bytes.NewBuffer(reqBody))
	request.Header.Set("Content-Type", "application/json")

	r := httptest.NewRecorder()
	router.ServeHTTP(r, request)

	response := r.Body.String()
	var result map[string]interface{}
	json.Unmarshal([]byte(response), &result)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	assert.Contains(suite.T(), response, "email must be a gmail address")
}
func (suite *UserControllerTestSuite) TestEditEmailPassword_InvalidPassword() {
	user := &model.User{
		ID:       "1",
		Name:     "John Doe",
		Email:    "johndoe@gmail.com",
		Password: "Pas1",
	}
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.PUT("/user/pass/:user_id", controller.EditEmailPassword)

	suite.useCaseMock.On("FindById", uint(1)).Return(user, nil)

	reqBody, _ := json.Marshal(map[string]string{
		"email":    "johndoe@gmail.com",
		"password": "Pas1",
	})
	request, _ := http.NewRequest(http.MethodPut, "/user/pass/1", bytes.NewBuffer(reqBody))
	request.Header.Set("Content-Type", "application/json")

	r := httptest.NewRecorder()
	router.ServeHTTP(r, request)

	response := r.Body.String()
	var result map[string]interface{}
	json.Unmarshal([]byte(response), &result)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	assert.Contains(suite.T(), response, "Invalid Input: password must have at least 8 characters")
}
func (suite *UserControllerTestSuite) TestEditEmailPassword_InvalidPasswordFormat() {
	user := &model.User{
		ID:       "1",
		Name:     "John Doe",
		Email:    "johndoe@gmail.com",
		Password: "password",
	}
	controller := NewUserController(suite.useCaseMock, suite.bankMock, suite.photoMock)
	router := setupRouter()
	router.PUT("/user/pass/:user_id", controller.EditEmailPassword)

	suite.useCaseMock.On("FindById", uint(1)).Return(user, nil)

	reqBody, _ := json.Marshal(map[string]string{
		"email":    "johndoe@gmail.com",
		"password": "password",
	})
	request, _ := http.NewRequest(http.MethodPut, "/user/pass/1", bytes.NewBuffer(reqBody))
	request.Header.Set("Content-Type", "application/json")

	r := httptest.NewRecorder()
	router.ServeHTTP(r, request)

	response := r.Body.String()
	var result map[string]interface{}
	json.Unmarshal([]byte(response), &result)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	assert.Contains(suite.T(), response, "password must contain at least one uppercase letter and one number")
}

func TestIsValidPassword(t *testing.T) {
	cases := []struct {
		password string
		expected bool
	}{
		{"Abcdef123", true},
		{"Abcdef", false},
		{"abcdef123", false},
		{"ABCDEF123", true},
		{"123456", false},
		{"abcxyz", false},
		{"", false},
	}

	for _, c := range cases {
		actual := isValidPassword(c.password)
		if actual != c.expected {
			t.Errorf("isValidPassword(%q) == %t, expected %t", c.password, actual, c.expected)
		}
	}
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.routerMock = gin.Default()
	suite.useCaseMock = new(UserUseCaseMock)
	suite.bankMock = new(BankAccUsecaseMock)

	suite.photoMock = new(PhotoUseCaseMock)
}
func TestUserController(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
