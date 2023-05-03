package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/model/response"
	"github.com/ReygaFitra/inc-final-project.git/usecase"
	"github.com/ReygaFitra/inc-final-project.git/utils"
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
	defer logger.Close()
	logrus.SetOutput(logger)

	res := c.usecase.FindUsers()
	if res == nil {
		logrus.Error("Failed to get users")
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "Users not found")
		return
	}

	logrus.Info("Success to get users")
	response.JSONSuccess(ctx.Writer, http.StatusOK, res)
}

func (c *UserController) FindUserByUsername(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logger.Close()
	logrus.SetOutput(logger)
	username := ctx.Param("username")
	res, _ := c.usecase.FindByUsername(username)
	if res == nil {
		logrus.Error("User not found")
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "User not found")
		return
	}
	logrus.Info("Succes to get user")
	response.JSONSuccess(ctx.Writer, http.StatusOK, res)
}

func (c *UserController) Register(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logger.Close()
	logrus.SetOutput(logger)
	newUser := model.UserCreate{}
	if err := ctx.BindJSON(&newUser); err != nil {
		logrus.Errorf("Invalid Input : %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Input")
		return
	}

	res, _ := c.usecase.Register(&newUser)
	logrus.Info("Success Register User")
	response.JSONSuccess(ctx.Writer, http.StatusCreated, res)
}

func (c *UserController) EditEmailPassword(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logger.Close()
	logrus.SetOutput(logger)
	// Retrieve the user_id parameter from the request
	user_id_str := ctx.Param("user_id")
	user_id, err := strconv.ParseUint(user_id_str, 10, 64)
	if err != nil {
		logrus.Errorf("Invalid user ID : %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Retrieve the existing user
	existingUser, _ := c.usecase.FindById(uint(user_id))
	if existingUser == nil {
		logrus.Error("User not found")
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "User not found")
		return
	}

	// Create a new User instance and map the properties from existingUser
	user := &model.User{}
	if err := mapstructure.Decode(existingUser, user); err != nil {
		logrus.Errorf("Failed Edit user : %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to edit user")
		return
	}

	// Parse the request body to update the user
	if err := ctx.BindJSON(user); err != nil {
		logrus.Errorf("Invalid input : %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid input")
		return
	}

	// Update the user and save changes
	updatedUser := c.usecase.EditEmailPassword(user)
	if updatedUser == "" {
		logrus.Error("Failed Edit user")
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Failed to edit user")
		return
	}
	logrus.Info("Edit Successfully")
	response.JSONSuccess(ctx.Writer, http.StatusOK, updatedUser)
}

func (c *UserController) EditProfile(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logger.Close()
	logrus.SetOutput(logger)
	// Retrieve the user_id parameter from the request
	user_id_str := ctx.Param("user_id")
	user_id, err := strconv.ParseUint(user_id_str, 10, 64)
	if err != nil {
		logrus.Errorf("Invalid user ID : %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Retrieve the existing user
	existingUser, _ := c.usecase.FindById(uint(user_id))
	if existingUser == nil {
		logrus.Error("User not found")
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "User not found")
		return
	}

	// Create a new User instance and map the properties from existingUser
	user := &model.User{}
	if err := mapstructure.Decode(existingUser, user); err != nil {
		logrus.Errorf("Failed to edit user : %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to edit user")
		return
	}

	// Parse the request body to update the user
	if err := ctx.BindJSON(user); err != nil {
		logrus.Errorf("Invalid input : %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid input")
		return
	}

	// Update the user and save changes
	updatedUser := c.usecase.EditProfile(user)
	if updatedUser == "" {
		logrus.Error("Failed to edit user")
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to edit user")
		return
	}

	logrus.Info("Edit Profile Succesfully")
	response.JSONSuccess(ctx.Writer, http.StatusOK, updatedUser)
}

func (c *UserController) Unreg(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logger.Close()
	logrus.SetOutput(logger)
	username := ctx.Param("username")
	user := &model.User{
		Username: username,
	}

	res := c.usecase.Unreg(user)
	logrus.Info("Delete Successfully")
	response.JSONSuccess(ctx.Writer, http.StatusOK, res)
}

func NewUserController(usercase usecase.UserUseCase) *UserController {
	controller := UserController{
		usecase: usercase,
	}
	return &controller
}
