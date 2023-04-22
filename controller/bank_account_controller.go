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
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to get user Bank Account")
		return
	}

	response.JSONSuccess(ctx.Writer, http.StatusOK, result)
}

func (c *BankAccController) FindBankAccByUsername(ctx *gin.Context) {
	username := ctx.Param("username")
	existingUser, _ := c.bankAccUsecase.FindBankAccByUsername(username)
	if existingUser == nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "Bank not found")
		return
	}

	response.JSONSuccess(ctx.Writer, http.StatusOK, existingUser)
}

func (c *BankAccController) FindBankAccByAccountID(ctx *gin.Context) {
	user_id_str := ctx.Param("account_id")
	user_id, err := strconv.ParseUint(user_id_str, 10, 64)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Account ID")
		return
	}

	existingUser, _ := c.bankAccUsecase.FindBankAccByAccountID(uint(user_id))
	if existingUser == nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "Bank not found")
		return
	}

	response.JSONSuccess(ctx.Writer, http.StatusOK, existingUser)
}

func (c *BankAccController) CreateBankAccount(ctx *gin.Context) {
	username := ctx.GetString("username")
	var newBankAcc model.BankAccResponse
	err := ctx.BindJSON(&newBankAcc)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := c.bankAccUsecase.Register(username, &newBankAcc)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to create Bank Account")
		return
	}

	response.JSONSuccess(ctx.Writer, http.StatusOK, result)
}

func (c *BankAccController) Edit(ctx *gin.Context) {
	account_id_str := ctx.Param("account_id")
	account_id, err := strconv.ParseUint(account_id_str, 10, 64)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Account ID")
		return
	}

	existingUser, _ := c.bankAccUsecase.FindBankAccByAccountID(uint(account_id))
	if existingUser == nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "Bank not found")
		return
	}
	user := &model.BankAcc{}
	if err := mapstructure.Decode(existingUser, user); err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to edit Bank")
		return
	}

	if err := ctx.BindJSON(user); err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid input")
		return
	}
	updateBank := c.bankAccUsecase.Edit(user)
	response.JSONSuccess(ctx.Writer, http.StatusOK, updateBank)
}

func (c *BankAccController) UnregAll(ctx *gin.Context) {
	username := ctx.Param("username")
	user := &model.BankAcc{
		Username: username,
	}
	res := c.bankAccUsecase.UnregAll(user)

	response.JSONSuccess(ctx.Writer, http.StatusOK, res)
}

func (c *BankAccController) UnregByAccountId(ctx *gin.Context) {
	accountIDStr := ctx.Param("account_id")
	accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid account ID")
		return
	}

	err = c.bankAccUsecase.UnregByAccountId(uint(accountID))
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to delete bank account")
		return
	}

	response.JSONSuccess(ctx.Writer, http.StatusOK, "Bank account deleted successfully")
}

func NewBankAccController(u usecase.BankAccUsecase) *BankAccController {
	controller := BankAccController{
		bankAccUsecase: u,
	}
	return &controller
}
