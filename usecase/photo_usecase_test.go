package usecase

import (
	"errors"
	"testing"

	"github.com/ReygaFitra/inc-final-project.git/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyPhoto = []model.PhotoUrl{
	{
		Photo_ID: 1,
		UserID:   1,
		Url:      "/Developments/Golang/src/final-project-inc/file/avatar1.jpg",
	},
	{
		Photo_ID: 2,
		UserID:   2,
		Url:      "/Developments/Golang/src/final-project-inc/file/avatar2.jpg",
	},
}

type photoRepoMock struct {
	mock.Mock
}

func (r *photoRepoMock) Create(photo *model.PhotoUrl) error {
	args := r.Called(photo)
	if args[0] != nil {
		return args.Error(0)
	}
	return nil
}

func (r *photoRepoMock) GetByID(id uint) (*model.PhotoUrl, error) {
	args := r.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PhotoUrl), args.Error(1)
}

func (r *photoRepoMock) Update(photo *model.PhotoUrl) error {
	args := r.Called(photo)
	if args[0] != nil {
		return args.Error(0)
	}
	return nil
}

func (r *photoRepoMock) Delete(id uint) string {
	args := r.Called(id)
	if args[0] != nil {
		return "delete failed"
	}
	return "delete success"
}

type PhotoUseCaseTestSuite struct {
	photoRepoMock *photoRepoMock
	suite.Suite
}

// Test Upload
func (suite *PhotoUseCaseTestSuite) TestUpload_Success() {
	photoUC := NewPhotoUseCase(suite.photoRepoMock)
	suite.photoRepoMock.On("Create", &dummyPhoto[0]).Return(nil)
	err := photoUC.Upload(&dummyPhoto[0])
	assert.Nil(suite.T(), err)
}
func (suite *PhotoUseCaseTestSuite) TestUpload_Failed() {
	photoUC := NewPhotoUseCase(suite.photoRepoMock)
	suite.photoRepoMock.On("Create", &dummyPhoto[0]).Return(errors.New("Failed"))
	err := photoUC.Upload(&dummyPhoto[0])
	assert.NotNil(suite.T(), err)
}

// Test Download
func (suite *PhotoUseCaseTestSuite) TestDownload_Success() {
	id := dummyPhoto[0].UserID
	photoUC := NewPhotoUseCase(suite.photoRepoMock)
	suite.photoRepoMock.On("GetByID", id).Return(&dummyPhoto[0], nil)
	res, err := photoUC.Download(id)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), res)
	assert.Equal(suite.T(), &dummyPhoto[0], res)
}
func (suite *PhotoUseCaseTestSuite) TestDownload_Failed() {
	id := dummyPhoto[0].UserID
	photoUC := NewPhotoUseCase(suite.photoRepoMock)
	suite.photoRepoMock.On("GetByID", id).Return(nil, errors.New("failed to get photo"))
	res, err := photoUC.Download(id)
	assert.Nil(suite.T(), res)
	assert.NotNil(suite.T(), err)
}

// Test Edit
func (suite *PhotoUseCaseTestSuite) TestEdit_Success() {
	photoUC := NewPhotoUseCase(suite.photoRepoMock)
	suite.photoRepoMock.On("Update", &dummyPhoto[0]).Return(nil)
	err := photoUC.Edit(&dummyPhoto[0])
	assert.Nil(suite.T(), err)
}
func (suite *PhotoUseCaseTestSuite) TestEdit_Failed() {
	photoUC := NewPhotoUseCase(suite.photoRepoMock)
	suite.photoRepoMock.On("Update", &dummyPhoto[0]).Return(errors.New("Failed"))
	err := photoUC.Edit(&dummyPhoto[0])
	assert.NotNil(suite.T(), err)
}

// Test Remove
func (suite *PhotoUseCaseTestSuite) TestRemove_Success() {
	id := dummyPhoto[0].UserID
	photoUC := NewPhotoUseCase(suite.photoRepoMock)
	suite.photoRepoMock.On("Delete", id).Return(nil)
	res := photoUC.Remove(id)
	assert.Equal(suite.T(), "delete success", res)
	assert.NotNil(suite.T(), res)
}
func (suite *PhotoUseCaseTestSuite) TestRemove_Failed() {
	id := dummyPhoto[0].UserID
	photoUC := NewPhotoUseCase(suite.photoRepoMock)
	suite.photoRepoMock.On("Delete", id).Return("delete failed")
	res := photoUC.Remove(id)
	assert.Equal(suite.T(), "delete failed", res)
}

func (suite *PhotoUseCaseTestSuite) SetupTest() {
	suite.photoRepoMock = new(photoRepoMock)
}

func TestPhotoUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(PhotoUseCaseTestSuite))
}
