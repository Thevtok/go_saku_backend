package controller

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

func DotEnv(key string) string {
	// load .env file
	if err := godotenv.Load("config.env"); err != nil {
		log.Fatalln("error saat load .env file")
	}

	return os.Getenv(key)
}

type PhotoUseCaseMock struct {
	mock.Mock
}

func (r *PhotoUseCaseMock) Upload(photo *model.PhotoUrl) error {
	args := r.Called(photo)
	if args[0] != nil {
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
	if args[0] != nil {
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
// // Set environment variable for configuration
// err := os.Setenv("FILE_LOCATION", "/tmp/%s")
// require.NoError(suite.T(), err)

// // Initialize controller and mock usecase
// usecaseMock := new(PhotoUseCaseMock)
// controller := NewPhotoController(usecaseMock)

// // Set up request
// photo := dummyPhoto[0]
// reqBody, err := json.Marshal(photo)
// require.NoError(suite.T(), err)
// req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/user/photo/%d", photo.UserID), bytes.NewBuffer(reqBody))
// require.NoError(suite.T(), err)

// // Set up dependencies for Upload method
// usecaseMock.On("Upload", mock.Anything).Return(nil)
// ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
// ctx.Request = req

// // Call Upload method and check response
// controller.Upload(ctx)
// response := ctx.Writer.Body.String()
// var actualPhoto model.PhotoUrl
// err = json.Unmarshal([]byte(response), &actualPhoto)
// require.NoError(suite.T(), err)
// assert.Equal(suite.T(), http.StatusCreated, ctx.Writer.Status())
// assert.Equal(suite.T(), photo.Url, actualPhoto.Url)

// // Clean up environment variable
// err = os.Unsetenv("FILE_LOCATION")
// require.NoError(suite.T(), err)
// }

func (suite *PhotoControllerTestSuite) TestRemove_Success() {
   // Load environment variables
   err := godotenv.Load("../config.env")
   if err != nil {
	   suite.T().Fatal("Error loading .env file")
   }

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