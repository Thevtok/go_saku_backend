package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/model/response"
	"github.com/ReygaFitra/inc-final-project.git/usecase"
	"github.com/ReygaFitra/inc-final-project.git/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	usecase usecase.UserUseCase
}

func (c *UserController) FindUsers(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	res := c.usecase.FindUsers()
	if res == nil {
		logrus.Error("Failed to get users")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Users not found")
		return
	}

	logrus.Info("Success to get users")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, res)
}

func (c *UserController) FindUserByUsername(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	username := ctx.Param("username")
	res, _ := c.usecase.FindByUsername(username)
	if res == nil {
		logrus.Error("User not found")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "User not found")
		return
	}
	logrus.Info("Succes to get user")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, res)
}

func (c *UserController) Register(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	newUser := model.UserCreate{}
	if err := ctx.BindJSON(&newUser); err != nil {
		logrus.Errorf("Invalid Input : %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid Input")
		return
	}

	if newUser.Name == "" || newUser.Username == "" || newUser.Email == "" || newUser.Password == "" || newUser.Phone_Number == "" || newUser.Address == "" {
		logrus.Errorf("Invalid Input: Required fields are empty")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid Input: Required fields are empty")
		return
	}
	if !strings.HasSuffix(newUser.Email, "@gmail.com") {
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "email must be a gmail address")
		return
	}
	if len(newUser.Password) < 8 {
		logrus.Errorf("Invalid Input: password must have at least 8 characters")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid Input: password must have at least 8 characters")
		return
	}
	if !isValidPassword(newUser.Password) {
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "password must contain at least one uppercase letter and one number")
		return
	}
	if len(newUser.Phone_Number) < 11 || len(newUser.Phone_Number) > 13 {
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "phone_number must be 11 - 13 digit")
		return
	}

	res, err := c.usecase.Register(&newUser)
	if err != nil {
		logrus.Errorf("Failed to Register User : %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to Register User")
		return
	}

	logrus.Info("Success Register User")
	response.JSONSuccess(ctx.Writer, true, http.StatusCreated, res)
}

func (c *UserController) EditEmailPassword(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	// Retrieve the user_id parameter from the request
	user_id_str := ctx.Param("user_id")
	user_id, err := strconv.ParseUint(user_id_str, 10, 64)
	if err != nil {
		logrus.Errorf("Invalid user ID : %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Retrieve the existing user
	existingUser, _ := c.usecase.FindById(uint(user_id))
	if existingUser == nil {
		logrus.Error("User not found")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "User not found")
		return
	}

	// Create a new User instance and map the properties from existingUser
	user := &model.User{}
	if err := mapstructure.Decode(existingUser, user); err != nil {
		logrus.Errorf("Failed Edit user : %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to edit user")
		return
	}
	if err := ctx.BindJSON(user); err != nil {
		logrus.Errorf("Invalid input : %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid input")
		return
	}
	if !strings.HasSuffix(user.Email, "@gmail.com") {
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "email must be a gmail address")
		return
	}
	if len(user.Password) < 8 {
		logrus.Errorf("Invalid Input: password must have at least 8 characters")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid Input: password must have at least 8 characters")
		return
	}
	if !isValidPassword(user.Password) {
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "password must contain at least one uppercase letter and one number")
		return
	}

	// Parse the request body to update the user

	// Update the user and save changes
	updatedUser := c.usecase.EditEmailPassword(user)
	if updatedUser == "" {
		logrus.Error("Failed Edit user")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Failed to edit user")
		return
	}
	logrus.Info("Edit Successfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, updatedUser)
}

func (c *UserController) EditProfile(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	// Retrieve the user_id parameter from the request
	user_id_str := ctx.Param("user_id")
	user_id, err := strconv.ParseUint(user_id_str, 10, 64)
	if err != nil {
		logrus.Errorf("Invalid user ID : %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Retrieve the existing user
	existingUser, _ := c.usecase.FindById(uint(user_id))
	if existingUser == nil {
		logrus.Error("User not found")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "User not found")
		return
	}

	// Create a new User instance and map the properties from existingUser
	user := &model.User{}
	if err := mapstructure.Decode(existingUser, user); err != nil {
		logrus.Errorf("Failed to edit user : %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to edit user")
		return
	}

	// Parse the request body to update the user
	if err := ctx.BindJSON(user); err != nil {
		logrus.Errorf("Invalid input : %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid input")
		return
	}
	if len(user.Phone_Number) < 11 || len(user.Phone_Number) > 13 {
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "phone_number must be 11 - 13 digit")
		return
	}

	// Update the user and save changes
	updatedUser := c.usecase.EditProfile(user)
	if updatedUser == "" {
		logrus.Error("Failed to edit user")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to edit user")
		return
	}

	logrus.Info("Edit Profile Succesfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, updatedUser)
}

func (c *UserController) Unreg(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	username := ctx.Param("username")
	user := &model.User{
		Username: username,
	}

	res := c.usecase.Unreg(user)
	logrus.Info("Delete Successfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, res)
}

func NewUserController(usercase usecase.UserUseCase) *UserController {
	controller := UserController{
		usecase: usercase,
	}
	return &controller
}

func (ctx *UserController) AuthMiddlewareIDExist() gin.HandlerFunc {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	return func(c *gin.Context) {
		// Add log statement here
		logrus.Info("Authenticating user...")

		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			logrus.Errorf("unauthorized %v", err)
			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			logrus.Errorf("failed generate token")

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
		userID := uint((*claims)["user_id"].(float64))

		if !ok {
			logrus.Errorf("invalid claim user_id")
			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "invalid userId claim")
			c.Abort()
			return
		}

		requestedID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
		if err != nil {
			logrus.Errorf("invalid user_id")
			response.JSONErrorResponse(c.Writer, false, http.StatusBadRequest, "invalid user_id")
			c.Abort()
			return
		}

		exist, _ := ctx.usecase.FindById(uint(requestedID))
		if exist == nil {
			response.JSONErrorResponse(c.Writer, false, http.StatusNotFound, "User Not Found")
			c.Abort()
			return
		}

		if userID != uint(requestedID) {
			logrus.Errorf("you do not have permission to access this resource")

			response.JSONErrorResponse(c.Writer, false, http.StatusForbidden, "you do not have permission to access this resource")
			c.Abort()
			return
		}

		c.Set("email", email)
		c.Set("password", password)
		c.Set("user_id", userID)

		c.Next()
		logrus.Info("Success parsing middleware")
	}
}

func isValidPassword(password string) bool {
	hasNum := false
	hasUpper := false

	for _, char := range password {
		if unicode.IsNumber(char) {
			hasNum = true
		} else if unicode.IsUpper(char) {
			hasUpper = true
		}
	}

	return hasNum && hasUpper
}
