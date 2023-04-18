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
	"golang.org/x/crypto/bcrypt"
)

type LoginAuth struct {
	usecase usecase.UserUseCase
	jwtKey  []byte
}

var jwtKey = []byte(utils.DotEnv("KEY"))

func generateToken(user *model.User) (string, error) {
	// Set token claims
	claims := jwt.MapClaims{}
	claims["email"] = user.Email

	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	// Create token with claims and secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AuthMiddleware(jwtKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		email := claims["email"].(string)

		c.Set("email", email)

		c.Next()
	}
}
func (l *LoginAuth) Login(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	User, err := l.usecase.Login(user.Email, string(hashedPassword))
	if err != nil {
		log.Println(err) // log the error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	if User == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Print the hashed password stored in the database

	token, err := generateToken(User)
	if err != nil {
		log.Println(err) // log the error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
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
