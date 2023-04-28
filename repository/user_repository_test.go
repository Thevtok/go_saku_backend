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

var dummyCredentials = []model.Credentials{
	{
		Email: "email1@mail.com",
		Password: "password1",
		UserID: 1,
		Username: "username1",
		Role: "user",
	},
	{
		Email: "email2@mail.com",
		Password: "password2",
		UserID: 2,
		Username: "username2",
		Role: "admin",
	},
}

var dummyUserRespons = []model.UserResponse{
	{
		Name: "name1",
		Username: "username1",
		Email: "email1@mail.com",
		Phone_Number: "08111111",
		Address: "address1",
		Balance: 100000,
		Point: 20,
	},
	{
		Name: "name2",
		Username: "username2",
		Email: "email2@mail.com",
		Phone_Number: "08111111",
		Address: "address2",
		Balance: 100000,
		Point: 40,
	},
}

var dummyUser = []model.User{
	{
		ID: 1,
		Name: "name1",
		Username: "username1",
		Email: "email1@mail.com",
		Password: "password1",
		Phone_Number: "08111111",
		Address: "address1",
		Balance: 100000,
		Role: "user",
		Point: 10,
	},
	{
		ID: 2,
		Name: "name2",
		Username: "username2",
		Email: "email2@mail.com",
		Password: "password2",
		Phone_Number: "08111111",
		Address: "address2",
		Balance: 50000,
		Role: "user",
		Point: 10,
	},
}

var dummyUserCreate = []model.UserCreate{
	{
		Name: "name1",
		Username: "username1",
		Email: "email1@mail.com",
		Password: "password1",
		Phone_Number: "08111111",
		Address: "address1",
		Balance: 100000,
	},
	{
		Name: "name2",
		Username: "username2",
		Email: "email2@mail.com",
		Password: "password2",
		Phone_Number: "082222",
		Address: "address2",
		Balance: 100000,
	},
}

type UserRepositoryTestSuite struct {
	suite.Suite
	mockDB *sql.DB
	mockSql sqlmock.Sqlmock
}

// Test UpdateBalance
func (suite *UserRepositoryTestSuite) TestUpdateBalance_Success() {
	// userID := dummyUser[0].ID
	// newBalance := dummyUser[0].Balance
	// suite.mockSql.ExpectQuery("SELECT name, user_id, email, phone_number, address, balance, point FROM mst_users WHERE user_id = $1").WithArgs(userID).WillReturnRows(sqlmock.NewRows([]string{"name", "user_id", "email", "phone_number", "address", "balance", "point"}))
	// suite.mockSql.ExpectQuery("UPDATE mst_users SET balance = \\$1 WHERE user_id = \\$2").WithArgs(newBalance, userID).WillReturnRows(sqlmock.NewRows([]string{"balance", "user_id"}).AddRow(newBalance, userID))
	// userRepository := NewUserRepository(suite.mockDB)
	// err := userRepository.UpdateBalance(userID, newBalance)
	// assert.NotNil(suite.T(), err)
	user := dummyUser[0]
    suite.mockSql.ExpectQuery("SELECT name, user_id, email, phone_number, address, balance, point FROM mst_users WHERE user_id = \\$1").WithArgs(user.ID).WillReturnRows(sqlmock.NewRows([]string{"name", "user_id", "email", "phone_number", "address", "balance", "point"}))
    suite.mockSql.ExpectExec("UPDATE mst_users SET balance = \\$1 WHERE user_id = \\$2").WithArgs(user.Balance, user.ID).WillReturnResult(sqlmock.NewResult(1,1))
    userRepository := NewUserRepository(suite.mockDB)
    err := userRepository.UpdateBalance(user.ID, user.Balance)
    assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), nil)
}
func (suite *UserRepositoryTestSuite) TestUpdateBalance_Failed() {
	// userID := dummyUser[0].ID
    // newBalance := dummyUser[0].Balance
    // expectedErr := fmt.Errorf("failed to update balance")
    // suite.mockSql.ExpectQuery("UPDATE mst_users SET balance = \\$1 WHERE user_id = \\$2").
    //     WithArgs(newBalance, userID).
    //     WillReturnError(expectedErr)
    // userRepository := NewUserRepository(suite.mockDB)
    // err := userRepository.UpdateBalance(userID, newBalance)
    // assert.NotNil(suite.T(), err)
	user := dummyUser[0]
    expectedError := fmt.Errorf("failed to update balance")
    suite.mockSql.ExpectQuery("SELECT name, user_id, email, phone_number, address, balance, point FROM mst_users WHERE user_id = \\$1").WithArgs(user.ID).WillReturnRows(sqlmock.NewRows([]string{"name", "user_id", "email", "phone_number", "address", "balance", "point"}))
    suite.mockSql.ExpectExec("UPDATE mst_users SET balance = \\$1 WHERE user_id = \\$2").WithArgs(user.Balance, user.ID).WillReturnError(expectedError)
    userRepository := NewUserRepository(suite.mockDB)
    err := userRepository.UpdateBalance(user.ID, user.Balance)
    assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), nil)
}

//Test UpdatePoint
func (suite *UserRepositoryTestSuite) TestUpdatePoint_Success() {
	userID := dummyUser[0].ID
	newPoint := dummyUser[0].Point
	suite.mockSql.ExpectQuery("UPDATE mst_users SET point = \\$1 WHERE user_id = \\$2").WithArgs(newPoint, userID).WillReturnRows(sqlmock.NewRows([]string{"point", "user_id"}).AddRow(newPoint, userID))
	userRepository := NewUserRepository(suite.mockDB)
	err := userRepository.UpdatePoint(userID, newPoint)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), nil)
}
func (suite *UserRepositoryTestSuite) TestUpdatePoint_Failed() {
	userID := dummyUser[0].ID
    newPoint := dummyUser[0].Point
    expectedErr := fmt.Errorf("failed to update point")
    suite.mockSql.ExpectQuery("UPDATE mst_users SET point = \\$1 WHERE user_id = \\$2").
        WithArgs(newPoint, userID).
        WillReturnError(expectedErr)
    userRepository := NewUserRepository(suite.mockDB)
    err := userRepository.UpdatePoint(userID, newPoint)
    assert.NotNil(suite.T(), err)
}

// Test GetAll
func(suite *UserRepositoryTestSuite) TestGetAll_Success() {
	var users = dummyUserRespons[0]
	suite.mockSql.ExpectQuery("SELECT name, username, email, phone_number, address, balance, point from mst_users").WillReturnRows(sqlmock.NewRows([]string{"name", "username", "email", "phone_number", "address", "balance", "point"}).AddRow(users.Name, users.Username, users.Email, users.Phone_Number, users.Address, users.Balance, users.Point))
	userRepository := NewUserRepository(suite.mockDB)
	res := userRepository.GetAll()
	assert.NotNil(suite.T(), res)
}
func(suite *UserRepositoryTestSuite) TestGetAll_Failed() {
	suite.mockSql.ExpectQuery("SELECT name, username, email, phone_number, address, balance, point from mst_users").WillReturnError(errors.New("no data"))
	userRepository := NewUserRepository(suite.mockDB)
	res := userRepository.GetAll()
	assert.NotNil(suite.T(), res)
	assert.Equal(suite.T(), "no data", res)
}
func(suite *UserRepositoryTestSuite) TestGetAllScan_Failed() {
	var users = dummyUserRespons[0]
	rows := sqlmock.NewRows([]string{"name", "username", "email", "phone_number", "address", "balance", "point"})
	rows.AddRow(users.Name, users.Username, users.Email, users.Phone_Number, users.Address, users.Balance, users.Point)
	suite.mockSql.ExpectQuery("SELECT name, username, email, phone_number, address, balance, point from mst_users").WillReturnRows(rows)
	userRepository := NewUserRepository(suite.mockDB)
	res := userRepository.GetAll()
	assert.NotNil(suite.T(), res)
	assert.Error(suite.T(), errors.New("no data"))
}

// Test GetByUsername
func (suite *UserRepositoryTestSuite) TestGetByUsername_Success() {
	user := dummyUserRespons[0]
	suite.mockSql.ExpectQuery("SELECT name, username, email, phone_number, address, balance, point FROM mst_users WHERE username = $1").WithArgs(user.Username).WillReturnRows(sqlmock.NewRows([]string{"name", "username", "email", "phone_number", "address", "balance", "point"}).AddRow(user.Name, user.Username, user.Email, user.Phone_Number, user.Address, user.Balance, user.Point))
	userRepository := NewUserRepository(suite.mockDB)
	res, err := userRepository.GetByUsername(user.Username)
	assert.Nil(suite.T(), res)
	assert.NotNil(suite.T(), err)
}
func (suite *UserRepositoryTestSuite) TestGetByUsername_Failed() {
	user := dummyUserRespons[0]
	expectedErr := errors.New("some error")
	suite.mockSql.ExpectQuery("SELECT name, username, email, phone_number, address, balance, point, tx_count FROM mst_users WHERE username = $1").WithArgs(user.Username).WillReturnError(expectedErr)
	userRepository := NewUserRepository(suite.mockDB)
	res, err := userRepository.GetByUsername(user.Username)
	assert.Nil(suite.T(), res)
	assert.NotNil(suite.T(), err)
}

// Test GetByID
func (suite *UserRepositoryTestSuite) TestGetByiD_Success() {
	user := dummyUser[0]
	suite.mockSql.ExpectQuery("SELECT name, user_id, email, phone_number, address, balance, point FROM mst_users WHERE user_id = $1").WithArgs(user.ID).WillReturnRows(sqlmock.NewRows([]string{"name", "user_id", "email", "phone_number", "address", "balance", "point"}).AddRow(user.Name, user.ID, user.Email, user.Phone_Number, user.Address, user.Balance, user.Point))
	userRepository := NewUserRepository(suite.mockDB)
	res, err := userRepository.GetByiD(user.ID)
	assert.Nil(suite.T(), res)
	assert.NotNil(suite.T(), err)
}
func (suite *UserRepositoryTestSuite) TestGetByiD_Failed() {
	user := dummyUser[0]
	expectedErr := errors.New("some error")
	suite.mockSql.ExpectQuery("SELECT name, user_id, email, phone_number, address, balance, point FROM mst_users WHERE user_id = $1").WithArgs(user.ID).WillReturnError(expectedErr)
	userRepository := NewUserRepository(suite.mockDB)
	res, err := userRepository.GetByiD(user.ID)
	assert.Nil(suite.T(), res)
	assert.NotNil(suite.T(), err)
}

// TestUpdateProfile
func (suite *UserRepositoryTestSuite) TestUpdateProfile_Success() {
	user := dummyUser[0]
	suite.mockSql.ExpectExec("UPDATE mst_users SET name= \\$1,  phone_number= \\$2, address= \\$3, username= \\$4 WHERE user_id= \\$5").WithArgs(user).WillReturnResult(sqlmock.NewResult(1,1))
	userRepository := NewUserRepository(suite.mockDB)
	str := userRepository.UpdateProfile(&user)
	assert.NotNil(suite.T(), str)
}
func (suite *UserRepositoryTestSuite) TestUpdateProfile_Failed() {
	user := dummyUser[0]
	expectedError := fmt.Errorf("failed to update user")
	suite.mockSql.ExpectExec("UPDATE mst_users SET name= \\$1,  phone_number= \\$2, address= \\$3, username= \\$4 WHERE user_id= \\$5").WithArgs(user).WillReturnError(expectedError)
	userRepository := NewUserRepository(suite.mockDB)
	str := userRepository.UpdateProfile(&user)
	assert.NotNil(suite.T(), str)
}

// Test UpdateEmailPassword
func(suite *UserRepositoryTestSuite) TestUpdateEmailPassword_Success() {
	user := dummyUser[0]
	suite.mockSql.ExpectExec("UPDATE mst_users SET email= \\$1, password= \\$2  WHERE user_id= \\$3").WithArgs(user).WillReturnResult(sqlmock.NewResult(1,1))
	userRepository := NewUserRepository(suite.mockDB)
	str := userRepository.UpdateEmailPassword(&user)
	assert.NotNil(suite.T(), str)
}
func(suite *UserRepositoryTestSuite) TestUpdateEmailPassword_Failed() {
	user := dummyUser[0]
	expectedError := fmt.Errorf("failed to update user")
	suite.mockSql.ExpectExec("UPDATE mst_users SET email= \\$1, password= \\$2  WHERE user_id= \\$3").WithArgs(user).WillReturnError(expectedError)
	userRepository := NewUserRepository(suite.mockDB)
	str := userRepository.UpdateEmailPassword(&user)
	assert.NotNil(suite.T(), str)
}

// Test Delete
func (suite *UserRepositoryTestSuite) TestDelete_Success() {
	user := dummyUser[0]
	suite.mockSql.ExpectQuery("SELECT name, user_id, email, phone_number, address, balance, point FROM mst_users WHERE user_id = $1").WithArgs(user.ID).WillReturnRows(sqlmock.NewRows([]string{"name", "user_id", "email", "phone_number", "address", "balance", "point"}))

	suite.mockSql.ExpectExec("DELETE FROM mst_users WHERE username = \\$1").WithArgs(user).WillReturnResult(sqlmock.NewResult(1,1))
	userRepository := NewUserRepository(suite.mockDB)
	str := userRepository.Delete(&user)
	assert.NotNil(suite.T(), str)
}
func (suite *UserRepositoryTestSuite) TestDelete_Failed() {
	user := dummyUser[0]
	expectedError := fmt.Errorf("failed to delete user")
	suite.mockSql.ExpectExec("DELETE FROM mst_users WHERE username = \\$1").WithArgs(user).WillReturnError(expectedError)
	userRepository := NewUserRepository(suite.mockDB)
	str := userRepository.Delete(&user)
	assert.NotNil(suite.T(), str)
}

// Test Create
func (suite *UserRepositoryTestSuite) TestCreate_Success() {
	newUser := dummyUserCreate[0]
	suite.mockSql.ExpectExec("INSERT INTO mst_users (name, username, email, password, phone_number, address, balance, role, point) VALUES (\\$[0-9]+, \\$[0-9]+, \\$[0-9]+, \\$[0-9]+, \\$[0-9]+, \\$[0-9]+, \\$[0-9]+, \\$[0-9]+, \\$[0-9]+)").WithArgs(newUser.Name, newUser.Username, newUser.Email, newUser.Password, newUser.Phone_Number, newUser.Address, newUser.Balance, 0, "user", 0).WillReturnResult(sqlmock.NewResult(1,1))
	userRepository := NewUserRepository(suite.mockDB)
	res, err := userRepository.Create(&newUser)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), res)
}
func (suite *UserRepositoryTestSuite) TestCreate_Failed() {
	newUser := dummyUserCreate[0]
	expectedError := fmt.Errorf("failed to Create user")
	suite.mockSql.ExpectExec("INSERT INTO mst_users (name, username, email, password, phone_number, address, balance, role, point) VALUES (\\$[0-9]+, \\$[0-9]+, \\$[0-9]+, \\$[0-9]+, \\$[0-9]+, \\$[0-9]+, \\$[0-9]+, \\$[0-9]+, \\$[0-9]+)").WithArgs(newUser).WillReturnError(expectedError)
	userRepository := NewUserRepository(suite.mockDB)
	res, err := userRepository.Create(&newUser)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), res)
}

func (suite *UserRepositoryTestSuite) TestGetByEmailAndPassword_Success() {
	var user = dummyCredentials[0]
	suite.mockSql.ExpectQuery("SELECT user_id, username, password, role FROM mst_users WHERE email = \\$1").WithArgs(user.Email, user.Password).WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "password", "role", "email"}).AddRow(user.UserID, user.Username, user.Password, user.Role, user.Email))
	userRepository := NewUserRepository(suite.mockDB)
	res, err := userRepository.GetByEmailAndPassword(user.Email, user.Password)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), res)
}
func (suite *UserRepositoryTestSuite) TestGetByEmailAndPassword_Failed() {
	var user = dummyCredentials[0]
	expectedError := errors.New("user not found")
	suite.mockSql.ExpectQuery("SELECT user_id, username, password, role FROM mst_users WHERE email = \\$1").WithArgs(user.Email, user.Password).WillReturnError(expectedError)
	userRepository := NewUserRepository(suite.mockDB)
	res, err := userRepository.GetByEmailAndPassword(user.Email, user.Password)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), res)
}

// setup test
func (suite *UserRepositoryTestSuite) SetupTest() {
	mockDb, mockSql, err := sqlmock.New()
	if err != nil {
		log.Fatalln("Error database", err)
	}
	suite.mockDB = mockDb
	suite.mockSql = mockSql
}
func (suite *UserRepositoryTestSuite) TearDownTest() {
	suite.mockDB.Close()
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}