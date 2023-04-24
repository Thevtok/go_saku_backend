package usecase

import (
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
		Url: "url/file/photo.png",
	},
	{
		Photo_ID: 2,
		UserID: 2,
		Url: "url/file/photo2.jpg",
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

// Test Remove
func (suite *PhotoUseCaseTestSuite) TestRemove_Success() {
	id := dummyPhoto[0].UserID
	photoUC := NewPhotoUseCase(suite.repoMock)
	suite.repoMock.On("Delete", id).Return(nil)
	err := photoUC.Remove(id)
	assert.Equal(suite.T(), "delete success", err)
}
func (suite *PhotoUseCaseTestSuite) TestRemove_Failed() {
	id := dummyPhoto[0].UserID
	photoUC := NewPhotoUseCase(suite.repoMock)
	suite.repoMock.On("Delete", id).Return("delete failed")
	err := photoUC.Remove(id)
	assert.Equal(suite.T(), "delete failed", err)
}

func (suite *PhotoUseCaseTestSuite) SetupTest() {
	suite.repoMock = new(repoMock)
}

func TestPhotoUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(PhotoUseCaseTestSuite))
}