package controller

import (
	"net/http"
	"strconv"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/model/response"
	"github.com/ReygaFitra/inc-final-project.git/usecase"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

type UserController struct {
	usecase usecase.UserUseCase
}

func (c *UserController) FindUsers(ctx *gin.Context) {
	res := c.usecase.FindUsers()
	if res == nil {

		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to get users")
		return
	}
	response.JSONSuccess(ctx.Writer, http.StatusOK, res)
}

func (c *UserController) FindUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")
	res, _ := c.usecase.FindByUsername(username)
	if res == nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "User not found")
		return
	}
	response.JSONSuccess(ctx.Writer, http.StatusOK, res)
}

func (c *UserController) Register(ctx *gin.Context) {
	newUser := model.UserCreate{}

	if err := ctx.BindJSON(&newUser); err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Input")
		return
	}

	res, _ := c.usecase.Register(&newUser)
	response.JSONSuccess(ctx.Writer, http.StatusCreated, res)
}

func (c *UserController) Edit(ctx *gin.Context) {
	// Retrieve the user_id parameter from the request
	user_id_str := ctx.Param("user_id")
	user_id, err := strconv.ParseUint(user_id_str, 10, 64)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Retrieve the existing user
	existingUser, _ := c.usecase.FindById(uint(user_id))
	if existingUser == nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "User not found")
		return
	}

	// Create a new User instance and map the properties from existingUser
	user := &model.User{}
	if err := mapstructure.Decode(existingUser, user); err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to edit user")
		return
	}

	// Parse the request body to update the user
	if err := ctx.BindJSON(user); err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid input")
		return
	}

	// Update the user and save changes
	updatedUser := c.usecase.Edit(user)
	if updatedUser == "" {
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to edit user")
		return
	}

	response.JSONSuccess(ctx.Writer, http.StatusOK, updatedUser)
}

func (c *UserController) Unreg(ctx *gin.Context) {
	username := ctx.Param("username")

	user := &model.User{
		Username: username,
	}

	res := c.usecase.Unreg(user)
	response.JSONSuccess(ctx.Writer, http.StatusOK, res)
}

func NewUserController(u usecase.UserUseCase) *UserController {
	controller := UserController{
		usecase: u,
	}

	return &controller
}
