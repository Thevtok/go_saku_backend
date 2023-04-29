package controller

import (
	"testing"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/gin-gonic/gin"
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

func (r *PhotoUseCaseMock) Upload(photo *model.PhotoUrl) error {
	args := r.Called(photo)
	if args[0] == nil {
		return args.Error(0)
	}
	return nil
}

func (r *PhotoUseCaseMock) Download(id uint) (*model.PhotoUrl, error) {
	args := r.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PhotoUrl), args.Error(1)
}

func (u *PhotoUseCaseMock) Edit(photo *model.PhotoUrl) error {
	args := u.Called(photo)
	if args[0] == nil {
		return args.Error(0)
	}
	return nil
}

func (r *PhotoUseCaseMock) Remove(id uint) string {
	args := r.Called(id)
	if args[0] != nil {
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
// 	newPhoto := &dummyPhoto[0]
// 	NewPhotoController(suite.useCaseMock)
// 	r := httptest.NewRecorder()
// 	reqBody, _ := json.Marshal(newPhoto)
// 	request, _ := http.NewRequest(http.MethodPost, "/user/photo/:user_id", bytes.NewBuffer(reqBody))
// 	suite.routerMock.ServeHTTP(r, request)
// 	response := r.Body.String()
// 	var actualPhoto model.PhotoUrl
// 	json.Unmarshal([]byte(response), &actualPhoto)
// 	assert.Equal(suite.T(), http.StatusCreated, r.Code)
// }

// func (suite *PhotoControllerTestSuite) TestRemove_Success() {
//    // Load environment variables
//    err := godotenv.Load("../config.env")
//    if err != nil {
// 	   suite.T().Fatal("Error loading .env file")
//    }

//    // prepare mock expectation
//    id := dummyPhoto[0].UserID
//    suite.useCaseMock.On("Remove", id).Return("Success remove")

//    // perform request
//    w := httptest.NewRecorder()
//    req, _ := http.NewRequest("DELETE", "/user/photo/:user_id", nil)
//    suite.routerMock.ServeHTTP(w, req)

//    // assert response status code and body
//    assert.Equal(suite.T(), http.StatusOK, w.Code)
//    assert.Equal(suite.T(), "{\"message\":\"Success remove\"}", w.Body.String())
// }

func (suite *PhotoControllerTestSuite) SetupTest() {
	suite.routerMock = gin.Default()
	suite.useCaseMock = new(PhotoUseCaseMock)
}
func TestPhotoController(t *testing.T) {
    suite.Run(t, new(PhotoControllerTestSuite))
}