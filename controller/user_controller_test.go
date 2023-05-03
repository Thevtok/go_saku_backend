package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyCredentials = []model.Credentials{
	{
		Email:    "email1@mail.com",
		Password: "password1",
		UserID:   1,
		Username: "username1",
		Role:     "user",
	},
	{
		Email:    "email2@mail.com",
		Password: "password2",
		UserID:   2,
		Username: "username2",
		Role:     "admin",
	},
}

var dummyUserRespons = []model.UserResponse{
	{
		Name:         "name1",
		Username:     "username1",
		Email:        "email1@mail.com",
		Phone_Number: "08111111",
		Address:      "address1",
		Balance:      100000,
		Point:        20,
	},
	{
		Name:         "name2",
		Username:     "username2",
		Email:        "email2@mail.com",
		Phone_Number: "08111111",
		Address:      "address2",
		Balance:      100000,
		Point:        40,
	},
}

var dummyUser = []model.User{
	{
		ID:           1,
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
		ID:           2,
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

var dummyUserCreate = []model.UserCreate{
	{
		Name:         "name1",
		Username:     "username1",
		Email:        "email1@mail.com",
		Password:     "password1",
		Phone_Number: "08111111",
		Address:      "address1",
		Balance:      100000,
	},
	{
		Name:         "name2",
		Username:     "username2",
		Email:        "email2@mail.com",
		Password:     "password2",
		Phone_Number: "082222",
		Address:      "address2",
		Balance:      100000,
	},
}

type UserUseCaseMock struct {
	mock.Mock
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func (c *UserUseCaseMock) Register(user *model.UserCreate) (any, error) {
	args := c.Called(user)
	if args.Get(0) == nil {
		return nil, args.Get(0).(error)
	}
	return &dummyUserCreate[0], nil
}

func (r *UserUseCaseMock) Login(email string, password string) (*model.Credentials, error) {
	args := r.Called(email, password)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return &dummyCredentials[0], nil
}

func (r *UserUseCaseMock) FindUsers() any {
	args := r.Called()
	if args[0] == nil {
		return nil
	}
	return dummyUserRespons
}

func (r *UserUseCaseMock) FindByUsername(username string) (*model.UserResponse, error) {
	args := r.Called(username)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return &dummyUserRespons[0], nil
}

func (r *UserUseCaseMock) FindById(id uint) (*model.User, error) {
	args := r.Called(id)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return &dummyUser[0], nil
}

func (r *UserUseCaseMock) EditProfile(user *model.User) string {
	args := r.Called(user)
	if args[0] == nil {
		return "Failed update profile"
	}
	return "Success update profile"
}

func (r *UserUseCaseMock)EditEmailPassword(user *model.User) string {
	args := r.Called(user)
	if args[0] == nil {
		return "Failed update email password"
	}
	return ""
}

func (r *UserUseCaseMock) Unreg(user *model.User) string {
	args := r.Called(user)
	if args[0] == nil {
		return "Failed Delete user"
	}
	return "Success Delete user"
}

func (r *UserUseCaseMock) UpdateBalance(userID uint, newBalance uint) error {
	args := r.Called(userID, newBalance)
	if args[0] == nil {
		return args.Error(0)
	}
	return nil
}

func (r *UserUseCaseMock) UpdatePoint(userID uint, newPoint int) error {
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
}

func(suite *UserControllerTestSuite) TestFindUsers_Success() {
	Users := &dummyUserRespons
  	controller := NewUserController(suite.useCaseMock)
	router := setupRouter()
	router.GET("/user", controller.FindUsers)

	suite.useCaseMock.On("FindUsers").Return(Users)
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	suite.useCaseMock.AssertExpectations(suite.T())
}
func(suite *UserControllerTestSuite) TestFindUsers_Failed() {
  	controller := NewUserController(suite.useCaseMock)
	router := setupRouter()
	router.GET("/user", controller.FindUsers)

	suite.useCaseMock.On("FindUsers").Return(nil)
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
	suite.useCaseMock.AssertExpectations(suite.T())
}


func(suite *UserControllerTestSuite) TestFindByUsername_Success() {
	Users := &dummyUserRespons[0]
	controller := NewUserController(suite.useCaseMock)
	router := setupRouter()
	router.GET("/user/:username", controller.FindUserByUsername)

	suite.useCaseMock.On("FindByUsername", Users.Username).Return(Users, "Success to get user")
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/username1", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	suite.useCaseMock.AssertExpectations(suite.T())
}

// //Test Register
// func (suite *UserControllerTestSuite) TestRegister_Success() {
// 	newUser := &dummyUserCreate[0]
//     NewUserController(suite.useCaseMock)
//     suite.useCaseMock.On("Register", newUser).Return(newUser, nil)
//     r := httptest.NewRecorder()
//     reqBody, _ := json.Marshal(newUser)
//     request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
//     suite.routerMock.ServeHTTP(r, request)
//     response := r.Body.String()
//     var actualUser model.UserCreate
//     json.Unmarshal([]byte(response), &actualUser)
//     assert.Equal(suite.T(), http.StatusCreated, r.Code)
// }
// func (suite *UserControllerTestSuite) TestRegister_Success() {
// 	newUser := &dummyUserCreate[0]
//     controller := NewUserController(suite.useCaseMock)
//     controller.SetRouter(suite.routerMock)
//     suite.useCaseMock.On("Register", newUser).Return(newUser, nil)
//     r := httptest.NewRecorder()
//     reqBody, _ := json.Marshal(newUser)
//     request, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
//     suite.routerMock.ServeHTTP(r, request)
//     response := r.Body.String()
//     var actualUser model.UserCreate
//     json.Unmarshal([]byte(response), &actualUser)
//     assert.Equal(suite.T(), http.StatusCreated, r.Code)
//     assert.Equal(suite.T(), newUser, &actualUser)
// }


func (suite *UserControllerTestSuite) SetupTest() {
	suite.routerMock = gin.Default()
	suite.useCaseMock = new(UserUseCaseMock)
}
func TestUserController(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}