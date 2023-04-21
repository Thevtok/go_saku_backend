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
	user_id_str := ctx.Param("user_id")
	user_id, err := strconv.ParseUint(user_id_str, 10, 64)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid user ID")
		return
	}

	existingUser, err := c.bankAccUsecase.FindBankAccByID(uint(user_id))
	if existingUser == nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "Bank not found")
		return
	}

	response.JSONSuccess(ctx.Writer, http.StatusOK, existingUser)
}

func (c *BankAccController) Register(ctx *gin.Context) {
	var newBankAcc model.BankAccResponse

	if err := ctx.BindJSON(&newBankAcc); err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid input")
		return
	}

	result, _ := c.bankAccUsecase.Register(&newBankAcc)
	response.JSONSuccess(ctx.Writer, http.StatusCreated, result)
}

func (c *BankAccController) Edit(ctx *gin.Context) {
	user_id_str := ctx.Param("user_id")
	user_id, err := strconv.ParseUint(user_id_str, 10, 64)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid user ID")
		return
	}

	existingUser, err := c.bankAccUsecase.FindBankAccByID(uint(user_id))
	if existingUser == nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "Bank not found")
		return
	}
	user := &model.BankAcc{}
	if err := mapstructure.Decode(existingUser, user); err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to edit bank")
		return
	}

	if err := ctx.BindJSON(user); err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid input")
		return
	}
	updateBank := c.bankAccUsecase.Edit(user)
	response.JSONSuccess(ctx.Writer, http.StatusOK, updateBank)
}

func (c *BankAccController) Unreg(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid User ID")
		return
	}
	user := &model.BankAcc{
		UserID: uint(userID),
	}

	result := c.bankAccUsecase.Unreg(user)
	if result != "success" {
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, result)
		return
	}

	response.JSONSuccess(ctx.Writer, http.StatusOK, "Bank Account unregistered successfully")
}

func NewBankAccController(u usecase.BankAccUsecase) *BankAccController {
	controller := BankAccController{
		bankAccUsecase: u,
	}

	return &controller
}
