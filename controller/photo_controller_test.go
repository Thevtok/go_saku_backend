package controller

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/gin-gonic/gin"
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

type PhotoUseCaseMock struct {
	mock.Mock
}

func (c *PhotoUseCaseMock) Upload(ctx context.Context, id uint, file multipart.File, header *multipart.FileHeader) error {
	args := c.Called(ctx, id, file, header)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (c *PhotoUseCaseMock) Download(id uint) (*model.PhotoUrl, error) {
	args := c.Called(id)
	if args.Get(0) != nil {
		return &model.PhotoUrl{}, args.Get(0).(error)
	}
	return args.Get(0).(*model.PhotoUrl), nil
}

func (c *PhotoUseCaseMock) Edit(photo *model.PhotoUrl, id uint, file multipart.File, header *multipart.FileHeader) error {
	args := c.Called(photo, id, file, header)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (c *PhotoUseCaseMock) Remove(id uint) string {
	args := c.Called(id)
	if args.Get(0) != nil {
		return "Failed remove"
	}
	return "Success remove"
}

type PhotoControllerTestSuite struct {
	suite.Suite
	routerMock *gin.Engine
	useCaseMock *PhotoUseCaseMock
}

// Test Upload
// func (suite *PhotoControllerTestSuite) TestUpload_Success() {
// 	photo := dummyPhoto[0]
// 	NewPhotoController(suite.useCaseMock)
// 	r := httptest.NewRecorder()
// 	reqBody, _ := json.Marshal(photo)
// 	request, _ := http.NewRequest(http.MethodPost, "/user/photo/:user_id", bytes.NewBuffer(reqBody))
// 	suite.routerMock.ServeHTTP(r, request)
// 	response := r.Body.String()
// 	var actualPhoto model.PhotoUrl
// 	json.Unmarshal([]byte(response), &actualPhoto)
// 	assert.Equal(suite.T(), http.StatusOK, r.Code)
// 	assert.Equal(suite.T(), photo.Url, actualPhoto.Url)
// }

func (suite *PhotoControllerTestSuite) TestRemove_Success() {
    // prepare mock expectation
	id := dummyPhoto[0].UserID
    suite.useCaseMock.On("Remove", id).Return("Success remove")

    // perform request
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("DELETE", "/user/photo/:user_id", nil)
    suite.routerMock.ServeHTTP(w, req)

    // assert response status code and body
    assert.Equal(suite.T(), http.StatusOK, w.Code)
    assert.Equal(suite.T(), "{\"message\":\"Success remove\"}", w.Body.String())
}


func TestPhotoController(t *testing.T) {
    suite.Run(t, new(PhotoControllerTestSuite))
}