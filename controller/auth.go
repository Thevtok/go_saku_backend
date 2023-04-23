package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/usecase"
	"github.com/ReygaFitra/inc-final-project.git/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type LoginAuth struct {
	usecase usecase.UserUseCase
	jwtKey  []byte
}

var jwtKey = []byte(utils.DotEnv("KEY"))

func generateToken(user *model.Credentials) (string, error) {
	// Set token claims

	claims := jwt.MapClaims{}
	claims["email"] = user.Email
	claims["password"] = user.Password
	claims["username"] = user.Username
	claims["user_id"] = uint(user.UserID)

	claims["role"] = user.Role

	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	// Print out the claims
	fmt.Printf("Claims: %+v\n", claims)

	// Create token with claims and secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	// Print out the token string
	fmt.Printf("Token: %s\n", tokenString)

	return tokenString, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		claims := token.Claims.(*jwt.MapClaims)
		email, ok := (*claims)["email"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email claim"})
			c.Abort()
			return
		}
		password, ok := (*claims)["password"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password claim"})
			c.Abort()
			return
		}
		username, ok := (*claims)["username"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username claim"})
			c.Abort()
			return
		}
		requestedID := c.Param("username")
		if username != requestedID {
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have permission to access this resource"})
			c.Abort()
			return
		}

		c.Set("email", email)
		c.Set("password", password)
		c.Set("username", username)

		c.Next()
	}
}

func (l *LoginAuth) Login(c *gin.Context) {
	var user model.Credentials

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	// Retrieve the user by email
	foundUser, err := l.usecase.Login(user.Email, user.Password)
	if err != nil {
		log.Println(err) // log the error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	// Verify that the provided password matches the stored hashed password
	err = utils.CheckPasswordHash(user.Password, foundUser.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate a token and return it to the client
	token, err := generateToken(foundUser)
	if err != nil {
		log.Println(err) // log the error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	// Return the token and the user ID in the context
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func NewUserAuth(u usecase.UserUseCase) *LoginAuth {
	loginauth := LoginAuth{
		usecase: u,
		jwtKey:  jwtKey,
	}
	return &loginauth
}
func AuthMiddlewareID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add log statement here
		log.Println("AuthMiddlewareID called")

		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Add log statement here
		log.Println("Token string:", tokenString)

		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Add log statement here
		log.Println("Token parsed successfully")

		claims := token.Claims.(*jwt.MapClaims)
		email, ok := (*claims)["email"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email claim"})
			c.Abort()
			return
		}
		password, ok := (*claims)["password"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password claim"})
			c.Abort()
			return
		}
		user_id, ok := (*claims)["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id claim"})
			c.Abort()
			return
		}
		requestedID, err := strconv.ParseFloat(c.Param("user_id"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			c.Abort()
			return
		}
		if user_id != requestedID {
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have permission to access this resource"})
			c.Abort()
			return
		}

		c.Set("email", email)
		c.Set("password", password)
		c.Set("user_id", uint(requestedID))

		c.Next()
	}
}

func AuthMiddlewareRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		claims := token.Claims.(*jwt.MapClaims)
		email, ok := (*claims)["email"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email claim"})
			c.Abort()
			return
		}
		password, ok := (*claims)["password"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password claim"})
			c.Abort()
			return
		}
		role, ok := (*claims)["role"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid role claim"})
			c.Abort()
			return
		}

		if role != "master" {
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have permission to access this resource"})
			c.Abort()
			return
		}

		c.Set("email", email)
		c.Set("password", password)
		c.Set("role", role)

		c.Next()
	}
}
