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

var dummyPhoto = []model.PhotoUrl{
	{
		Photo_ID: 1,
		UserID:   "1",
		Url:      "/Developments/Golang/src/final-project-inc/file/avatar1.jpg",
	},
	{
		Photo_ID: 2,
		UserID:   "2",
		Url:      "/Developments/Golang/src/final-project-inc/file/avatar2.jpg",
	},
}

type PhotoRepositoryTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSql sqlmock.Sqlmock
}

// Test Create
func (suite *PhotoRepositoryTestSuite) TestCreate_Success() {
	newPhoto := dummyPhoto[0]
	suite.mockSql.ExpectExec("INSERT INTO mst_photo_url \\(url_photo, user_id\\) VALUES").WithArgs(
		newPhoto.Url,
		newPhoto.UserID,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	photoRepository := NewPhotoRepository(suite.mockDB)
	err := photoRepository.Create(&newPhoto)
	assert.Nil(suite.T(), err)
}
func (suite *PhotoRepositoryTestSuite) TestCreate_Failed() {
	newPhoto := dummyPhoto[0]
	suite.mockSql.ExpectExec("INSERT INTO mst_photo_url").WillReturnError(errors.New("Failed"))
	photoRepository := NewPhotoRepository(suite.mockDB)
	err := photoRepository.Create(&newPhoto)
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, "Failed")
}

// Test GetByID
func (suite *PhotoRepositoryTestSuite) TestGetByID_Success() {
	photo := dummyPhoto[0]
	suite.mockSql.ExpectQuery("SELECT url_photo, user_id from mst_photo_url WHERE user_id = \\$1").WithArgs(photo.UserID).WillReturnRows(sqlmock.NewRows([]string{"url_photo", "user_id"}).AddRow(photo.Url, photo.UserID))
	photoRepository := NewPhotoRepository(suite.mockDB)
	res, err := photoRepository.GetByID(photo.UserID)
	assert.Nil(suite.T(), err)
	photo.Photo_ID = 0
	assert.Equal(suite.T(), &photo, res)
}
func (suite *PhotoRepositoryTestSuite) TestGetByID_Failed() {
	photo := dummyPhoto[0]
	expectedErrors := errors.New("some error")
	suite.mockSql.ExpectQuery("SELECT url_photo, user_id from mst_photo_url WHERE user_id = \\$1").WithArgs(photo.UserID).WillReturnError(expectedErrors)
	photoRepository := NewPhotoRepository(suite.mockDB)
	res, err := photoRepository.GetByID(photo.UserID)
	assert.NotNil(suite.T(), res)
	assert.Nil(suite.T(), err)

}

// Test Update
func (suite *PhotoRepositoryTestSuite) TestUpdate_Success() {
	photo := dummyPhoto[0]
	suite.mockSql.ExpectExec("UPDATE mst_photo_url SET url_photo = \\$1 WHERE user_id = \\$2").WithArgs(
		photo.Url,
		photo.UserID,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	photoRepository := NewPhotoRepository(suite.mockDB)
	err := photoRepository.Update(&photo)
	assert.Nil(suite.T(), err)
}
func (suite *PhotoRepositoryTestSuite) TestUpdate_Failed() {
	photo := dummyPhoto[0]
	suite.mockSql.ExpectQuery("UPDATE mst_photo_url SET url_photo = \\$1 WHERE user_id = \\$2").WithArgs(
		photo.Url,
		photo.UserID,
	).WillReturnError(fmt.Errorf("Update Failed"))
	photoRepository := NewPhotoRepository(suite.mockDB)
	err := photoRepository.Update(&photo)
	assert.NotNil(suite.T(), err)
}

// Test Delete
func (suite *PhotoRepositoryTestSuite) TestDelete_Succes() {
	id := dummyPhoto[0].UserID
	suite.mockSql.ExpectExec("DELETE FROM mst_photo_url WHERE user_id = \\$1").WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))
	photoRepository := NewPhotoRepository(suite.mockDB)
	res := photoRepository.Delete(id)
	assert.Equal(suite.T(), "Delete photo successfully", res)
}
func (suite *PhotoRepositoryTestSuite) TestDelete_Failed() {
	id := dummyPhoto[0].UserID
	suite.mockSql.ExpectExec("DELETE FROM mst_photo_url WHERE user_id = \\$1").WithArgs(id).WillReturnError(fmt.Errorf("Failed delete photo"))
	photoRepository := NewPhotoRepository(suite.mockDB)
	errString := photoRepository.Delete(id)
	assert.NotNil(suite.T(), errString)
}

// setup test
func (suite *PhotoRepositoryTestSuite) SetupTest() {
	mockDb, mockSql, err := sqlmock.New()
	if err != nil {
		log.Fatalln("Error database", err)
	}
	suite.mockDB = mockDb
	suite.mockSql = mockSql
}
func (suite *PhotoRepositoryTestSuite) TearDownTest() {
	suite.mockDB.Close()
}

func TestPhotoRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PhotoRepositoryTestSuite))
}
