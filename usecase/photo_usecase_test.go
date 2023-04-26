package usecase

import (
	"errors"
	"testing"

	"github.com/ReygaFitra/inc-final-project.git/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyPhoto = []model.PhotoUrl {
	{
		Photo_ID: 1,
		UserID: 1,
		Url: "/Developments/Golang/src/final-project-inc/file/avatar1.jpg",
	},
	{
		Photo_ID: 2,
		UserID: 2,
		Url: "/Developments/Golang/src/final-project-inc/file/avatar2.jpg",
	},
}

type repoMock struct {
	mock.Mock
}

func (r *repoMock) Create(photo *model.PhotoUrl) error {
	args := r.Called(photo)
	if args[0] != nil {
		return args.Error(0)
	}
	return nil
}

func (r *repoMock) GetByID(id uint) (*model.PhotoUrl, error) {
	args := r.Called(id)
	if args[0] != nil {
		return &model.PhotoUrl{}, args.Error(0)
	}
	return &dummyPhoto[0], nil
}

func (r *repoMock) Update(photo *model.PhotoUrl) error {
	args := r.Called(photo)
	if args[0] != nil {
		return args.Error(0)
	}
	return nil
}

func (r *repoMock) Delete(id uint) string {
	args := r.Called(id)
	if args[0] != nil {
		return "delete failed"
	}
	return "delete success"
}

type PhotoUseCaseTestSuite struct {
	repoMock *repoMock
	suite.Suite
}

// Test Upload
// func (suite *PhotoUseCaseTestSuite) TestUpload_Success() {
// 	// id dummy
// 	id := dummyPhoto[0].UserID
// 	// file/url photo dummy
// 	file, err := os.Open(dummyPhoto[0].Url)
// 	if err != nil {
// 	panic(err)
// 	}
// 	defer file.Close()
// 	// header dummy
// 	header := &multipart.FileHeader{
// 		Filename: dummyPhoto[0].Url,
// 	}
// 	photoUC := NewPhotoUseCase(suite.repoMock)
// 	suite.repoMock.On("Create", mock.AnythingOfType("*model.PhotoUrl")).Return(nil)
// 	err = photoUC.Upload(context.Background(), id, file, header)
// 	assert.Nil(suite.T(), err)
// }
// func (suite *PhotoUseCaseTestSuite) TestUpload_Failed() {
// 	// id dummy
// 	id := dummyPhoto[0].UserID
// 	// file/url photo dummy
// 	file, err := os.Open(dummyPhoto[0].Url)
// 	if err != nil {
// 	panic(err)
// 	}
// 	defer file.Close()
// 	// header dummy
// 	header := &multipart.FileHeader{
// 		Filename: "/Developments/Golang/src/final-project-inc/file/%s",
// 	}
// 	photoUC := NewPhotoUseCase(suite.repoMock)
// 	suite.repoMock.On("Create", &dummyPhoto[0]).Return(errors.New("Failed to upload"))
// 	err = photoUC.Upload(context.Background(), id, file, header)
// 	assert.NotNil(suite.T(), err)
// 	assert.Equal(suite.T(), err, "Failed to upload")
// }

// Test Download
func (suite *PhotoUseCaseTestSuite) TestDownload_Success() {
	id := dummyPhoto[0].UserID
	photoUC := NewPhotoUseCase(suite.repoMock)
	suite.repoMock.On("GetByID", id).Return(&dummyPhoto[0], nil)
	res, err := photoUC.Download(id)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), &dummyPhoto[0], res)
}
func (suite *PhotoUseCaseTestSuite) TestDownload_Failed() {
	id := dummyPhoto[0].UserID
	photoUC := NewPhotoUseCase(suite.repoMock)
	suite.repoMock.On("GetByID", id).Return(&model.PhotoUrl{}, errors.New("Failed to get photo"))
	_, err := photoUC.Download(id)
	assert.Equal(suite.T(), &model.PhotoUrl{}, err)
}

// Test Remove
func (suite *PhotoUseCaseTestSuite) TestRemove_Success() {
	id := dummyPhoto[0].UserID
	photoUC := NewPhotoUseCase(suite.repoMock)
	suite.repoMock.On("Delete", id).Return(nil)
	res := photoUC.Remove(id)
	assert.Equal(suite.T(), "delete success", res)
	assert.NotNil(suite.T(), res)
}
func (suite *PhotoUseCaseTestSuite) TestRemove_Failed() {
	id := dummyPhoto[0].UserID
	photoUC := NewPhotoUseCase(suite.repoMock)
	suite.repoMock.On("Delete", id).Return("delete failed")
	res := photoUC.Remove(id)
	assert.Equal(suite.T(), "delete failed", res)
}

func (suite *PhotoUseCaseTestSuite) SetupTest() {
	suite.repoMock = new(repoMock)
}

func TestPhotoUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(PhotoUseCaseTestSuite))
}