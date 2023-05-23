package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyPhoto = []model.PhotoUrl{
	{
		Photo_ID: 1,
		UserID:   "1",
		Url:      utils.DotEnv("FILE_LOCATION_DUMMY"),
	},
	{
		Photo_ID: 2,
		UserID:   "2",
		Url:      "/Developments/Golang/src/final-project-inc/file/avatar2.jpg",
	},
}

type PhotoUseCaseMock struct {
	mock.Mock
}

func setupRouterPhoto() *gin.Engine {
	r := gin.Default()
	return r
}

func (r *PhotoUseCaseMock) Upload(photo *model.PhotoUrl) error {
	args := r.Called(photo)
	if args[0] == nil {
		return args.Error(1)
	}
	return nil
}

func (r *PhotoUseCaseMock) Download(id string) (*model.PhotoUrl, error) {
	args := r.Called(id)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	// return args.Get(0).(*model.PhotoUrl), nil
	return &dummyPhoto[0], nil
}

func (u *PhotoUseCaseMock) Edit(photo *model.PhotoUrl) error {
	args := u.Called(photo)
	if args[0] == nil {
		return args.Error(0)
	}
	return nil
}

func (r *PhotoUseCaseMock) Remove(id string) string {
	args := r.Called(id)
	if args[0] == nil {
		return "Failed remove"
	}
	return "Success remove"
}

type PhotoControllerTestSuite struct {
	suite.Suite
	routerMock  *gin.Engine
	useCaseMock *PhotoUseCaseMock
}

func (suite *PhotoControllerTestSuite) TestDownload_Success() {
	Photo := &dummyPhoto[0]
	controller := NewPhotoController(suite.useCaseMock)
	router := setupRouterPhoto()
	router.GET("/user/photo/:user_id", controller.Download)

	suite.useCaseMock.On("Download", Photo.UserID).Return(Photo, nil)
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/photo/1", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	suite.useCaseMock.AssertExpectations(suite.T())
}
func (suite *PhotoControllerTestSuite) TestDownload_Failed() {
	Photo := &dummyPhoto[0]
	controller := NewPhotoController(suite.useCaseMock)
	router := setupRouterPhoto()
	router.GET("/user/photo/:user_id", controller.Download)

	suite.useCaseMock.On("Download", Photo.UserID).Return(nil, errors.New("Error file not found"))
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/user/photo/1", nil)
	router.ServeHTTP(r, request)
	assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)
	suite.useCaseMock.AssertExpectations(suite.T())
}
func (suite *PhotoControllerTestSuite) TestDownload_UserNotFound() {
	controller := NewPhotoController(suite.useCaseMock)
	router := setupRouterPhoto()
	router.GET("/user/photo/:user_id", controller.Download)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/photo/abc", nil)
	router.ServeHTTP(r, req)
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	suite.useCaseMock.AssertNotCalled(suite.T(), "Download")
}

func (suite *PhotoControllerTestSuite) TestRemove_Failed() {
	controller := NewPhotoController(suite.useCaseMock)
	router := setupRouterPhoto()
	router.DELETE("/user/photo/:user_id", controller.Remove)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/user/photo/abc", nil)
	router.ServeHTTP(r, req)
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	suite.useCaseMock.AssertNotCalled(suite.T(), "Remove")
}

func (suite *PhotoControllerTestSuite) SetupTest() {
	suite.routerMock = gin.Default()
	suite.useCaseMock = new(PhotoUseCaseMock)
}
func TestPhotoController(t *testing.T) {
	suite.Run(t, new(PhotoControllerTestSuite))
}
