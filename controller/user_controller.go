package controller

import (
	"net/http"
	"strconv"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/model/response"
	"github.com/ReygaFitra/inc-final-project.git/usecase"
	"github.com/gin-gonic/gin"
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

func (c *UserController) FindUserByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid user ID")
		return
	}
	res := c.usecase.FindByID(uint(id))
	if res == nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "User not found")
		return
	}
	response.JSONSuccess(ctx.Writer, http.StatusOK, res)
}

func (c *UserController) Register(ctx *gin.Context) {
	newUser := model.User{}

	if err := ctx.BindJSON(&newUser); err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Input")
		return
	}

	res := c.usecase.Register(&newUser)
	response.JSONSuccess(ctx.Writer, http.StatusCreated, res)
}

func (c *UserController) Edit(ctx *gin.Context) {
	var user model.User

	if err := ctx.BindJSON(&user); err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Input")
		return
	}

	res := c.usecase.Edit(&user)
	response.JSONSuccess(ctx.Writer, http.StatusOK, res)
}

func (c *UserController) Unreg(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid user ID")
		return
	}

	res := c.usecase.Unreg(uint(id))
	response.JSONSuccess(ctx.Writer, http.StatusOK, res)
}

func NewUserController(u usecase.UserUseCase) *UserController {
	controller := UserController{
		usecase: u,
	}

	return &controller
}
