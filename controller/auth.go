package controller

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/model/response"
	"github.com/ReygaFitra/inc-final-project.git/usecase"
	"github.com/ReygaFitra/inc-final-project.git/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LoginAuth struct {
	usecase usecase.UserUseCase
	jwtKey  []byte
}

var jwtKey = []byte(utils.DotEnv("KEY"))

func generateToken(user *model.Credentials) (string, error) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	// Set token claims
	claims := jwt.MapClaims{}
	claims["email"] = user.Email
	claims["password"] = user.Password
	claims["username"] = user.Username
	claims["user_id"] = uint(user.UserID)
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	// Print out the claims

	// Create token with claims and secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		logrus.Errorf("Failed Claim: %v", err)
		return "", err
	}

	logrus.Info("Claim Token Succesfully")

	return tokenString, nil
}

func AuthMiddleware() gin.HandlerFunc {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			logrus.Errorf("unauthorized %v", err)

			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			logrus.Info("Claim Token Succesfully")
			return jwtKey, nil

		})

		if err != nil || !token.Valid {
			logrus.Errorf("unauthorized %v", err)
			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		claims := token.Claims.(*jwt.MapClaims)
		email, ok := (*claims)["email"].(string)
		if !ok {
			logrus.Errorf("invalid claim email")
			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "invalid email claim")
			c.Abort()
			return
		}
		password, ok := (*claims)["password"].(string)
		if !ok {
			logrus.Errorf("invalid claim password ")
			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "invalid password claim")
			c.Abort()
			return
		}
		username, ok := (*claims)["username"].(string)
		if !ok {
			logrus.Errorf("invalid claim username")
			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "invalid userId claim")
			c.Abort()
			return
		}
		requestedID := c.Param("username")
		if username != requestedID {
			response.JSONErrorResponse(c.Writer, false, http.StatusForbidden, "you do not have permission to access this resource")
			c.Abort()
			return
		}

		c.Set("email", email)
		c.Set("password", password)
		c.Set("username", username)

		c.Next()
		logrus.Info("Success parsing midleware")
	}
}

func (l *LoginAuth) Login(c *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	var user model.Credentials

	if err := c.ShouldBindJSON(&user); err != nil {
		logrus.Errorf("invalid json")
		response.JSONErrorResponse(c.Writer, false, http.StatusBadRequest, "invalid json")

		return
	}

	// Retrieve the user by email
	foundUser, err := l.usecase.Login(user.Email, user.Password)
	if err != nil {
		logrus.Errorf("invalid email or password")
		response.JSONErrorResponse(c.Writer, false, http.StatusBadRequest, "invalid email or password")

		return
	}

	// Verify that the provided password matches the stored hashed password
	err = utils.CheckPasswordHash(user.Password, foundUser.Password)
	if err != nil {
		logrus.Errorf("invalid credentials")

		response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// Generate a token and return it to the client
	token, err := generateToken(foundUser)
	if err != nil {
		logrus.Errorf("failed generate token")
		response.JSONErrorResponse(c.Writer, false, http.StatusInternalServerError, "failed generate token")
		return
	}

	response.JSONSuccess(c.Writer, true, http.StatusOK, gin.H{"token": token})
	// Return the token and the user ID in the context

	logrus.Info("Success getting token")
}

func NewUserAuth(u usecase.UserUseCase) *LoginAuth {
	loginauth := LoginAuth{
		usecase: u,
		jwtKey:  jwtKey,
	}
	return &loginauth
}

func AuthMiddlewareID() gin.HandlerFunc {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	return func(c *gin.Context) {
		// Add log statement here

		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			logrus.Errorf("unauthorized %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			logrus.Info("Claim Token Succesfully")
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			logrus.Errorf("failed generate token")

			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		claims := token.Claims.(*jwt.MapClaims)
		email, ok := (*claims)["email"].(string)
		if !ok {
			logrus.Errorf("invalid claim email")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email claim"})
			c.Abort()
			return
		}
		password, ok := (*claims)["password"].(string)
		if !ok {
			logrus.Errorf("invalid claim password")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password claim"})
			c.Abort()
			return
		}
		user_id, ok := (*claims)["user_id"].(float64)
		if !ok {
			logrus.Errorf("invalid claim user_id")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id claim"})
			c.Abort()
			return
		}
		requestedID, err := strconv.ParseFloat(c.Param("user_id"), 64)
		if err != nil {
			logrus.Errorf("invalid user_id")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			c.Abort()
			return
		}
		if user_id != requestedID {
			logrus.Errorf("you do not have permission to access this resource")
			c.JSON(http.StatusForbidden, gin.H{"error": "you do not have permission to access this resource"})
			c.Abort()
			return
		}

		c.Set("email", email)
		c.Set("password", password)
		c.Set("user_id", uint(requestedID))

		c.Next()
		logrus.Info("Success parsing midleware")
	}
}

func AuthMiddlewareRole() gin.HandlerFunc {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			logrus.Errorf("unauthorized %v", err)
			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "unauthorized")
			c.Abort()
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			logrus.Info("Claim Token Succesfully")
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			logrus.Errorf("failed generate token")
			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "unauthorized")
			c.Abort()
		}

		claims := token.Claims.(*jwt.MapClaims)
		email, ok := (*claims)["email"].(string)
		if !ok {
			logrus.Errorf("invalid claim email")
			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "invalid email claim")
			c.Abort()
			return
		}
		password, ok := (*claims)["password"].(string)
		if !ok {
			logrus.Errorf("invalid claim password")
			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "invalid password claim")
			c.Abort()
			return
		}
		role, ok := (*claims)["role"].(string)
		if !ok {
			logrus.Errorf("invalid claim role")
			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "invalid role claim")
			c.Abort()
			return
		}

		if role != "master" {
			logrus.Errorf("you do not have permission to access this resource")
			response.JSONErrorResponse(c.Writer, false, http.StatusForbidden, "you do not have permission to access this resource")
			c.Abort()
			return
		}

		c.Set("email", email)
		c.Set("password", password)
		c.Set("role", role)

		c.Next()
		logrus.Info("Success parsing midleware")
	}
}
