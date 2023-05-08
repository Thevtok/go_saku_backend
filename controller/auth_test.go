package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// func TestAuthMiddlewareID_Success(t *testing.T) {
// 	// Set up the test
// 	r := setupTest()

// 	// Generate a valid token for a master role user
// 	token, err := generateToken(&model.Credentials{
// 		Email:    "master@myapp.com",
// 		Password: "password",
// 		UserID:   1,
// 		Username: "masteruser",
// 		Role:     "master",
// 	})
// 	require.NoError(t, err, "Failed to generate a valid token")
// 	headers := map[string]string{
// 		"Authorization": token,
// 	}

// 	req, _ := http.NewRequest(http.MethodGet, "/user/bank/:user_id", nil)
// 	req.Header.Set("Authorization", headers["Authorization"])
// 	w := httptest.NewRecorder()

// 	// Perform the request
// 	r.ServeHTTP(w, req)

// 	// Check the response
// 	assert.Equal(t, http.StatusOK, w.Code, "Response status code should be 200 OK")

// }
// func TestAuthMiddlewareID_MissingAuthorizationHeader(t *testing.T) {
// 	// Set up the test
// 	r := setupTest()
// 	req, _ := http.NewRequest(http.MethodGet, "/user/bank/:user_id", nil)
// 	w := httptest.NewRecorder()

// 	// Perform the request
// 	r.ServeHTTP(w, req)

// 	// Check the response
// 	assert.Equal(t, http.StatusUnauthorized, w.Code, "Response status code should be 401 Unauthorized")
// 	assert.JSONEq(t, `{"error":"unauthorized"}`, w.Body.String(), "Response body should be a JSON error message")
//  }
//  func TestAuthMiddlewareID_InvalidToken(t *testing.T) {
// 	// Set up the test
// 	r := setupTest()
// 	req, _ := http.NewRequest(http.MethodGet, "/user/bank/:user_id", nil)
// 	req.Header.Set("Authorization", "Bearer invalidtoken")
// 	w := httptest.NewRecorder()

// 	// Perform the request
// 	r.ServeHTTP(w, req)

// 	// Check the response
// 	assert.Equal(t, http.StatusUnauthorized, w.Code, "Response status code should be 401 Unauthorized")
// 	assert.JSONEq(t, `{"error":"unauthorized"}`, w.Body.String(), "Response body should be a JSON error message")
//  }

// func TestAuthMiddlewareRole_Success(t *testing.T) {
// 	// Set up the test
// 	r := setupTest()

// 	// Generate a valid token for a master role user
// 	token, err := generateToken(&model.Credentials{
// 		Email:    "master@myapp.com",
// 		Password: "password",
// 		UserID:   1,
// 		Username: "masteruser",
// 		Role:     "master",
// 	})
// 	require.NoError(t, err, "Failed to generate a valid token")
// 	headers := map[string]string{
// 		"Authorization": token,
// 	}

// 	req, _ := http.NewRequest(http.MethodGet, "/user/bank", nil)
// 	req.Header.Set("Authorization", headers["Authorization"])
// 	w := httptest.NewRecorder()

// 	// Perform the request
// 	r.ServeHTTP(w, req)

// 	// Check the response
// 	assert.Equal(t, http.StatusOK, w.Code, "Response status code should be 200 OK")
// 	assert.JSONEq(t, `{"message":"Hello, World!"}`, w.Body.String(), "Response body should be a JSON success message")
// }
// func TestAuthMiddlewareRole_MissingAuthorizationHeader(t *testing.T) {
//    // Set up the test
//    r := setupTest()
//    req, _ := http.NewRequest(http.MethodGet, "/user/bank", nil)
//    w := httptest.NewRecorder()

//    // Perform the request
//    r.ServeHTTP(w, req)

//    // Check the response
//    assert.Equal(t, http.StatusUnauthorized, w.Code, "Response status code should be 401 Unauthorized")
//    assert.JSONEq(t, `{"error":"unauthorized"}`, w.Body.String(), "Response body should be a JSON error message")
// }
// func TestAuthMiddlewareRole_InvalidToken(t *testing.T) {
// 	// Set up the test
// 	r := setupTest()
// 	req, _ := http.NewRequest(http.MethodGet, "/user/bank", nil)
// 	req.Header.Set("Authorization", "Bearer invalidtoken")
// 	w := httptest.NewRecorder()

// 	// Perform the request
// 	r.ServeHTTP(w, req)

// 	// Check the response
// 	assert.Equal(t, http.StatusUnauthorized, w.Code, "Response status code should be 401 Unauthorized")
// 	assert.JSONEq(t, `{"error":"unauthorized"}`, w.Body.String(), "Response body should be a JSON error message")
//  }
//  func TestAuthMiddlewareRole_UnauthorizedRole(t *testing.T) {
// 	// Set up the test
// 	r := setupTest()
// 	token, _ := generateToken(&model.Credentials{Email: "admin@mail.com", Password: "admin123", UserID: 1, Username: "admin", Role: "admin"})
// 	req, _ := http.NewRequest(http.MethodGet, "/user/bank", nil)
// 	req.Header.Set("Authorization", headers["Authorization"])
// 	w := httptest.NewRecorder()

// 	// Perform the request
// 	r.ServeHTTP(w, req)

// 	// Check the response
// 	assert.Equal(t, http.StatusForbidden, w.Code, "Response status code should be 403 Forbidden")
// 	assert.JSONEq(t, `{"error":"you do not have permission to access this resource"}`, w.Body.String(), "Response body should be a JSON error message")
//  }

func setupTest() *gin.Engine {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	r := gin.New()

	authMiddlewareRole := AuthMiddlewareRole()
	authMiddlewareID := AuthMiddlewareID()

	r.GET("/user/bank", authMiddlewareRole, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})
	r.GET("/user/bank/:user_id", authMiddlewareID, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	return r
}
