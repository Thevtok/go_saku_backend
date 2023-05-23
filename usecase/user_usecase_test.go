package usecase

import (
	"errors"
	"testing"

	"github.com/ReygaFitra/inc-final-project.git/model"
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

type userRepoMock struct {
	mock.Mock
}

func (r *userRepoMock) GetByEmailAndPassword(email string, password string, token string) (*model.User, error) {
	args := r.Called(email, password)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return &dummyUser[0], nil
}

func (r *userRepoMock) GetAll() any {
	args := r.Called()
	if args[0] != nil {
		return nil
	}
	return dummyUser
}

func (r *userRepoMock) GetByUsername(username string) (*model.User, error) {
	args := r.Called(username)
	if args[0] != nil {
		return nil, args.Error(1)
	}
	return &dummyUser[0], nil
}
func (r *userRepoMock) GetByPhone(username string) (*model.User, error) {
	args := r.Called(username)
	if args[0] != nil {
		return nil, args.Error(1)
	}
	return &dummyUser[0], nil
}

func (r *userRepoMock) GetByiD(id string) (*model.User, error) {
	args := r.Called(id)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*model.User), args.Error(1)
}
func (r *userRepoMock) GetByIDToken(id string) (*model.User, error) {
	args := r.Called(id)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*model.User), args.Error(1)
}

func (r *userRepoMock) Create(user *model.User) (any, error) {
	args := r.Called(user)
	if args[0] != nil {
		return nil, args.Error(1)
	}
	return &dummyUser[0], nil
}
func (r *userRepoMock) SaveDeviceToken(userID string, token string) error {
	args := r.Called(userID, token)
	if args[0] != nil {
		return args.Error(1)
	}
	return nil
}

func (r *userRepoMock) UpdateProfile(user *model.User) string {
	args := r.Called(user)
	if args[0] != nil {
		return "Failed update profile"
	}
	return "Success update profile"
}

func (r *userRepoMock) UpdateEmailPassword(user *model.User) string {
	args := r.Called(user)
	if args[0] == nil {
		return "Failed update email password"
	}
	return ""
}

func (r *userRepoMock) Delete(user *model.User) string {
	args := r.Called(user)
	if args[0] != nil {
		return "Failed Delete user"
	}
	return "Success Delete user"
}

func (r *userRepoMock) UpdateBalance(userID string, newBalance int) error {
	args := r.Called(userID, newBalance)
	if args[0] != nil {
		return args.Error(0)
	}
	return nil
}

func (r *userRepoMock) UpdatePoint(userID string, newPoint int) error {
	args := r.Called(userID, newPoint)
	if args[0] != nil {
		return args.Error(0)
	}
	return nil
}

type UserUseCaseTestSuite struct {
	userRepoMock *userRepoMock
	suite.Suite
}

// Test FindById
func (suite *UserUseCaseTestSuite) TestFindById_Success() {
	user := &dummyUser[0]
	userUC := NewUserUseCase(suite.userRepoMock)
	suite.userRepoMock.On("GetByiD", user.ID).Return(user, nil)

	res, err := userUC.FindById(user.ID)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user, res)
}

func (suite *UserUseCaseTestSuite) TestFindById_Failed() {
	user := &dummyUser[0]
	expectedErr := errors.New("user not found")
	userUC := NewUserUseCase(suite.userRepoMock)
	suite.userRepoMock.On("GetByiD", user.ID).Return(nil, expectedErr)
	res, err := userUC.FindById(user.ID)
	assert.Nil(suite.T(), res)
	assert.NotNil(suite.T(), err)
}

// Test FindByUsername
func (suite *UserUseCaseTestSuite) TestFindByUsername_Success() {
	user := &dummyUser[0]
	userUC := NewUserUseCase(suite.userRepoMock)
	suite.userRepoMock.On("GetByUsername", user.Username).Return(user, nil)
	res, err := userUC.FindByUsername(user.Username)
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), res)
}
func (suite *UserUseCaseTestSuite) TestFindByUsername_Failed() {
	user := &dummyUser[0]
	expectedErr := errors.New("user not found")
	userUC := NewUserUseCase(suite.userRepoMock)
	suite.userRepoMock.On("GetByUsername", user.Username).Return(nil, expectedErr)
	res, err := userUC.FindByUsername(user.Username)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
}

// Test GetAll
func (suite *UserUseCaseTestSuite) TestFindUsers_Success() {
	user := &dummyUser
	userUC := NewUserUseCase(suite.userRepoMock)
	suite.userRepoMock.On("GetAll").Return(user)
	res := userUC.FindUsers()
	assert.Nil(suite.T(), res)
}
func (suite *UserUseCaseTestSuite) TestFindUsers_Failed() {
	userUC := NewUserUseCase(suite.userRepoMock)
	suite.userRepoMock.On("GetAll").Return(nil)
	res := userUC.FindUsers()
	assert.NotNil(suite.T(), res)
}

// Test EditProfile
func (suite *UserUseCaseTestSuite) TestEditProfile_Success() {
	user := &dummyUser[0]
	userUC := NewUserUseCase(suite.userRepoMock)
	suite.userRepoMock.On("UpdateProfile", user).Return(nil)
	res := userUC.EditProfile(user)
	assert.NotNil(suite.T(), res)
	assert.Equal(suite.T(), "Success update profile", res)
}
func (suite *UserUseCaseTestSuite) TestEditProfile_Failed() {
	user := &dummyUser[0]
	userUC := NewUserUseCase(suite.userRepoMock)
	suite.userRepoMock.On("UpdateProfile", user).Return("Failed update profile")
	res := userUC.EditProfile(user)
	assert.Equal(suite.T(), "Failed update profile", res)
}

// Test Unreg
func (suite *UserUseCaseTestSuite) TestUnreg_Success() {
	user := &dummyUser[0]
	userUC := NewUserUseCase(suite.userRepoMock)
	suite.userRepoMock.On("Delete", user).Return(nil)
	res := userUC.Unreg(user)
	assert.NotNil(suite.T(), res)
	assert.Equal(suite.T(), "Success Delete user", res)
}
func (suite *UserUseCaseTestSuite) TestUnreg_Failed() {
	user := &dummyUser[0]
	userUC := NewUserUseCase(suite.userRepoMock)
	suite.userRepoMock.On("Delete", user).Return("Failed Delete user")
	res := userUC.Unreg(user)
	assert.Equal(suite.T(), "Failed Delete user", res)
}

// Test Register
func (suite *UserUseCaseTestSuite) TestRegister_Success() {
	user := &dummyUser[0]
	userUC := NewUserUseCase(suite.userRepoMock)
	suite.userRepoMock.On("Create", user).Return(user, nil)
	res, err := userUC.Register(user)
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), res)
}
func (suite *UserUseCaseTestSuite) TestRegister_Failed() {
	user := &dummyUser[0]
	userUC := NewUserUseCase(suite.userRepoMock)
	suite.userRepoMock.On("Create", user).Return(nil, errors.New("Failed create user"))
	res, err := userUC.Register(user)
	assert.NotNil(suite.T(), res)
	assert.Nil(suite.T(), err)
}

// Test EditEmailPassword
func (suite *UserUseCaseTestSuite) TestEditEmailPassword_Success() {
	user := &dummyUser[0]
	userUC := NewUserUseCase(suite.userRepoMock)
	suite.userRepoMock.On("UpdateEmailPassword", user).Return("", nil)
	res := userUC.EditEmailPassword(user)
	assert.NotNil(suite.T(), res)
}
func (suite *UserUseCaseTestSuite) TestEditEmailPassword_Failed() {
	user := &dummyUser[0]
	userUC := NewUserUseCase(suite.userRepoMock)
	suite.userRepoMock.On("UpdateEmailPassword", user).Return(errors.New("Failed update email password"))
	res := userUC.EditEmailPassword(user)
	assert.NotNil(suite.T(), res)
}
func (suite *UserUseCaseTestSuite) TestEditEmailPassword_Error() {
	user := &dummyUser[0]
	userUC := NewUserUseCase(suite.userRepoMock)
	suite.userRepoMock.On("UpdateEmailPassword", user).Return("", errors.New("Some error"))
	res := userUC.EditEmailPassword(user)
	assert.Empty(suite.T(), res)
}

// Test Login
func (suite *UserUseCaseTestSuite) TestLogin_Success() {
	user := &dummyUser[0]
	userUC := NewUserUseCase(suite.userRepoMock)
	suite.userRepoMock.On("GetByEmailAndPassword", user.Email, user.Password).Return(user.Password, user.Username, user.ID, user.Role, nil)
	res, err := userUC.Login(user.Email, user.Password, "")
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), res)

	// user := &dummyCredentials[0]
	// userUC := NewUserUseCase(suite.userRepoMock)
	// expected := &model.Credentials{Email: "", Password: "password1", UserID: 0x1, Username: "username1", Role: "user"}
	// suite.userRepoMock.On("GetByEmailAndPassword", user.Email, user.Password).Return(expected, nil)
	// res, err := userUC.Login(user.Email, user.Password)
	// assert.Nil(suite.T(), res)
	// assert.Equal(suite.T(), expected, err)
}
func (suite *UserUseCaseTestSuite) TestLogin_Failed() {
	user := &dummyUser[0]
	userUC := NewUserUseCase(suite.userRepoMock)
	suite.userRepoMock.On("GetByEmailAndPassword", user.Email, user.Password).Return(nil, errors.New("Failed login"))
	res, err := userUC.Login(user.Email, user.Password, "")
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), res)
}

func (suite *UserUseCaseTestSuite) SetupTest() {
	suite.userRepoMock = new(userRepoMock)
}

func TestUserUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUseCaseTestSuite))
}
