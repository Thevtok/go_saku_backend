package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)



func TestAuthMiddlewareRole_Success(t *testing.T) {
	// Set up the test
	r := setupTest()

	// Generate a valid token for a master role user
	token, err := generateToken(&model.Credentials{
		Email:    "master@myapp.com",
		Password: "password",
		UserID:   1,
		Username: "masteruser",
		Role:     "master",
	})
	require.NoError(t, err, "Failed to generate a valid token")
	headers := map[string]string{
		"Authorization": token,
	}

	req, _ := http.NewRequest(http.MethodGet, "/user/bank", nil)
	req.Header.Set("Authorization", headers["Authorization"])
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code, "Response status code should be 200 OK")
	assert.JSONEq(t, `{"message":"Hello, World!"}`, w.Body.String(), "Response body should be a JSON success message")
}
func TestAuthMiddleware_Success(t *testing.T) {
    // Set up the test
    r := setupTest()

    // Generate a valid token
    token, err := generateToken(&model.Credentials{
        Email:    "user@test.com",
        Password: "password",
        UserID:   1,
        Username: "user1",
        Role:     "user",
    })
    if err != nil {
        t.Errorf("Failed to generate a valid token: %v", err)
    }
    headers := map[string]string{
        "Authorization": token,
    }

    req, _ := http.NewRequest(http.MethodGet, "/user/user1", nil)
    req.Header.Set("Authorization", headers["Authorization"])
    w := httptest.NewRecorder()

    // Perform the request
    r.ServeHTTP(w, req)

    // Check the response
    assert.Equal(t, http.StatusOK, w.Code, "Response status code should be 200 OK")
    assert.JSONEq(t, `{"message":"Hello, World!"}`, w.Body.String(), "Response body should be a JSON success message")
}
func TestAuthMiddleware_MissingToken(t *testing.T) {
    // Set up the test
    r := setupTest()

    req, _ := http.NewRequest(http.MethodGet, "/user/:username", nil)
    w := httptest.NewRecorder()

    // Perform the request
    r.ServeHTTP(w, req)

    // Check the response
    assert.Equal(t, http.StatusUnauthorized, w.Code, "Response status code should be 401 Unauthorized")
    assert.JSONEq(t, `{"status":false,"statusCode":401,"result":"unauthorized","message":"request failed"}`, w.Body.String(), "Response body should be a JSON unauthorized message")
}
func TestAuthMiddleware_InvalidToken(t *testing.T) {
    // Set up the test
    r := setupTest()

    // Generate an invalid token
    token := "invalid_token"
    headers := map[string]string{
        "Authorization": token,
    }

    req, _ := http.NewRequest(http.MethodGet, "/user/:username", nil)
    req.Header.Set("Authorization", headers["Authorization"])
    w := httptest.NewRecorder()

    // Perform the request
    r.ServeHTTP(w, req)

    // Check the response
    assert.Equal(t, http.StatusUnauthorized, w.Code, "Response status code should be 401 Unauthorized")
    assert.JSONEq(t, `{"status":false,"statusCode":401,"result":"unauthorized","message":"request failed"}`, w.Body.String(), "Response body should be a JSON unauthorized message")
}




func setupTest() *gin.Engine {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	r := gin.New()

	authMiddleware := AuthMiddleware()
	authMiddlewareRole := AuthMiddlewareRole()

	r.GET("/user/bank", authMiddlewareRole, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})
	r.GET("/user/:username", authMiddleware, func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
    })

	return r
}
