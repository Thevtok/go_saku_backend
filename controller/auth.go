package controller

import (
	"log"
	"net/http"
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

	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	// Create token with claims and secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

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
		email := (*claims)["email"].(string)
		password := (*claims)["password"].(string)

		c.Set("email", email)
		c.Set("password", password)

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

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func NewUserAuth(u usecase.UserUseCase) *LoginAuth {
	loginauth := LoginAuth{
		usecase: u,
		jwtKey:  jwtKey,
	}
	return &loginauth
}
