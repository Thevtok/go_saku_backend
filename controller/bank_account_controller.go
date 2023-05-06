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

type BankAccController struct {
	bankAccUsecase usecase.BankAccUsecase
}

func (c *BankAccController) FindAllBankAcc(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	result := c.bankAccUsecase.FindAllBankAcc()
	if result == nil {
		logrus.Errorf("Failed to get user Bank Account: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "Failed to get user Bank Account")
		return
	}

	logrus.Info("Bank Account loaded Successfully")
	response.JSONSuccess(ctx.Writer, http.StatusOK, result)
}

func (c *BankAccController) FindBankAccByUserID(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	if err != nil {
		logrus.Errorf("Invalid UserID: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid UserID")
		return
	}

	existingUser, err := c.bankAccUsecase.FindBankAccByUserID(uint(userID))
	if err != nil {
		logrus.Errorf("Bank Account not found: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "Bank Account not found")
		return
	}

	logrus.Info("Bank Account loaded Successfully")
	response.JSONSuccess(ctx.Writer, http.StatusOK, existingUser)
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
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid AccountID")
		return
	}

	existingUser, _ := c.bankAccUsecase.FindBankAccByAccountID(uint(userID))
	if existingUser == nil {
		logrus.Errorf("Bank Account not found: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "Bank Account not found")
		return
	}

	logrus.Info("Bank Account loaded Successfully")
	response.JSONSuccess(ctx.Writer, http.StatusOK, existingUser)
}

func (c *BankAccController) CreateBankAccount(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	userID, exists := ctx.Get("user_id")
	if !exists {
		logrus.Errorf("Failed to get UserID: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to get userID")
		return
	}

	var newBankAcc model.BankAccResponse
	err = ctx.BindJSON(&newBankAcc)
	if err != nil {
		logrus.Errorf("Invalid request body: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid request body")
		return
	}

	if newBankAcc.BankName == "" || newBankAcc.AccountNumber == "" || newBankAcc.AccountHolderName == "" {
		logrus.Errorf("Invalid Input: Required fields are empty")
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Input: Required fields are empty")
		return
	}

	result, err := c.bankAccUsecase.Register(userID.(uint), &newBankAcc)
	if err != nil {
		logrus.Errorf("Failed to create Bank Account: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to create Bank Account")
		return
	}

	logrus.Info("Bank Account created Successfully")
	response.JSONSuccess(ctx.Writer, http.StatusOK, result)
}

func (c *BankAccController) Edit(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	accountID, err := strconv.ParseUint(ctx.Param("account_id"), 10, 64)
	if err != nil {
		logrus.Errorf("Invalid AccountID: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid AccountID")
		return
	}

	existingUser, _ := c.bankAccUsecase.FindBankAccByAccountID(uint(accountID))
	if existingUser == nil {
		logrus.Errorf("Bank Account not found: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "Bank Account not found")
		return
	}

	user := &model.BankAcc{}
	if err := mapstructure.Decode(existingUser, user); err != nil {
		logrus.Errorf("Failed to edit Bank: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to edit Bank")
		return
	}

	if err := ctx.BindJSON(user); err != nil {
		logrus.Errorf("Invalid request body: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid request body")
		return
	}

	if user.BankName == "" || user.AccountNumber == "" || user.AccountHolderName == "" {
		logrus.Errorf("Invalid Input: Required fields are empty")
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Input: Required fields are empty")
		return
	}

	updateBank := c.bankAccUsecase.Edit(user)
	logrus.Info("Bank Account edited Successfully")
	response.JSONSuccess(ctx.Writer, http.StatusOK, updateBank)
}

func (c *BankAccController) UnregAll(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	if err != nil {
		logrus.Errorf("Invalid UserID: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid UserID")
		return
	}

	user := &model.BankAcc{
		UserID: uint(userID),
	}
	res := c.bankAccUsecase.UnregAll(user.UserID)
	logrus.Info("Bank Account deleted Successfully")
	response.JSONSuccess(ctx.Writer, http.StatusOK, res)
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
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid AccountID")
		return
	}

	err = c.bankAccUsecase.UnregByAccountID(uint(accountID))
	if err != nil {
		logrus.Errorf("Failed to delete Bank Account: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to delete Bank Account")
		return
	}
	logrus.Info("Bank Account deleted Successfully")
	response.JSONSuccess(ctx.Writer, http.StatusOK, "Bank account deleted successfully")
}

func NewBankAccController(u usecase.BankAccUsecase) *BankAccController {
	controller := BankAccController{
		bankAccUsecase: u,
	}
	return &controller
}
