package controller

import (
	"net/http"
	"strconv"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/model/response"
	"github.com/ReygaFitra/inc-final-project.git/usecase"
	"github.com/gin-gonic/gin"
)

type BankAccController struct {
	bankAccUsecase usecase.BankAccUsecase
}

func (c *BankAccController) FindAllBankAcc(ctx *gin.Context) {
	result := c.bankAccUsecase.FindAllBankAcc()
	if result == nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to ges user Bank Account")
		return
	}

	response.JSONSuccess(ctx.Writer, http.StatusOK, result)
}

func (c *BankAccController) FindBankAccByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Bank Account ID")
		return
	}

	result := c.bankAccUsecase.FindBankAccByID(uint(id))
	if result == nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "Bank Account not found")
		return
	}

	response.JSONSuccess(ctx.Writer, http.StatusOK, result)
}

func (c *BankAccController) Register(ctx *gin.Context) {
	var newUser model.BankAcc
	err := ctx.BindJSON(&newUser)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Input")
		return
	}

	result := c.bankAccUsecase.Register(&newUser)
	response.JSONSuccess(ctx.Writer, http.StatusCreated, result)
}

func (c *BankAccController) Edit(ctx *gin.Context) {
	var user model.BankAcc
	err := ctx.BindJSON(&user)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Input")
		return
	}

	result := c.bankAccUsecase.Edit(&user)
	response.JSONSuccess(ctx.Writer, http.StatusOK, result)
}

func (c *BankAccController) Unreg(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Bank Account ID")
		return
	}

	result := c.bankAccUsecase.Unreg(uint(id))
	response.JSONSuccess(ctx.Writer, http.StatusOK, result)
}

func NewBankAccController(u usecase.BankAccUsecase) *BankAccController {
	controller := BankAccController{
		bankAccUsecase: u,
	}

	return &controller
}
