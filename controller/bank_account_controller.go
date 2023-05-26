package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/model/response"
	"github.com/ReygaFitra/inc-final-project.git/usecase"
	"github.com/ReygaFitra/inc-final-project.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BankAccController struct {
	bankAccUsecase usecase.BankAccUsecase
}

func (c *BankAccController) FindBankAccByUserID(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	userID := ctx.Param("user_id")

	existingUser, err := c.bankAccUsecase.FindBankAccByUserID(userID)
	if err != nil {
		logrus.Errorf("Bank Account not found: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Bank Account not found")
		return
	}

	logrus.Info("Bank Account loaded Successfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, existingUser)
}

func (c *BankAccController) FindBankAccByAccountID(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	userID, err := strconv.ParseUint(ctx.Param("account_id"), 10, 64)
	if err != nil {
		logrus.Errorf("Invalid AccountID: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid AccountID")
		return
	}

	existingUser, _ := c.bankAccUsecase.FindBankAccByAccountID(uint(userID))
	if existingUser == nil {
		logrus.Errorf("Bank Account not found: %v", err)

		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Bank not found")
		return
	}
	logrus.Info("Bank Account loaded Successfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, existingUser)
}

func (c *BankAccController) CreateBankAccount(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	userID := ctx.Param("user_id")

	var newBankAcc model.BankAcc
	err = ctx.BindJSON(&newBankAcc)
	if err != nil {
		logrus.Errorf("Invalid request body: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid request body")
		return
	}
	newBankAcc.BankName = strings.ToLower(newBankAcc.BankName)

	if newBankAcc.BankName == "" || newBankAcc.AccountNumber == "" || newBankAcc.AccountHolderName == "" {
		logrus.Errorf("Invalid Input: Required fields are empty")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid Input: Required fields are empty")
		return
	}

	result, err := c.bankAccUsecase.Register(userID, &newBankAcc)
	if err != nil {
		logrus.Errorf("Failed to create Bank Account: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to create Bank Account")
		return
	}

	logrus.Info("Bank Account created Successfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusCreated, result)
}

func (c *BankAccController) UnregByAccountID(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	accountID, err := strconv.ParseUint(ctx.Param("account_id"), 10, 64)
	if err != nil {
		logrus.Errorf("Invalid AccountID: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid AccountID")
		return
	}

	err = c.bankAccUsecase.UnregByAccountID(uint(accountID))
	if err != nil {
		logrus.Errorf("Failed to delete Bank Account: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to delete Bank Account")
		return
	}
	logrus.Info("Bank Account deleted Successfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, "Bank account deleted successfully")
}

func NewBankAccController(u usecase.BankAccUsecase) *BankAccController {
	controller := BankAccController{
		bankAccUsecase: u,
	}
	return &controller
}
