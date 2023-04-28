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
	"github.com/sirupsen/logrus"
)

type TransactionController struct {
	txUsecase   usecase.TransactionUseCase
	userUsecase usecase.UserUseCase
	bankUsecase usecase.BankAccUsecase
	cardUsecase usecase.CardUsecase
}

func (c *TransactionController) CreateDepositBank(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}
	defer logger.Close()
	logrus.SetOutput(logger)

	// Parse user_id parameter
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		logrus.Errorf("Invalid UserID: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID parameter"})
		return
	}

	// Parse bank_account_id parameter
	bankAccID, err := strconv.Atoi(ctx.Param("bank_account_id"))
	if err != nil {
		logrus.Errorf("Invalid Bank AccountID: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Bank AccountID parameter"})
		return
	}

	// Retrieve bank account by bank_account_id
	bankAcc, err := c.bankUsecase.FindBankAccByAccountID(uint(bankAccID))
	if err != nil {
		logrus.Errorf("Failed to create Deposit Transaction: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to create Deposit Transaction")
		return
	}

	// Check if bank account belongs to the given user_id
	if bankAcc.UserID != uint(userID) {
		logrus.Errorf("Bank Account doesn't belong to the given UserID: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bank Account doesn't belong to the given UserID"})
		return
	}

	// Check if bank account belongs to the given user_id
	if bankAcc.UserID != uint(userID) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bank account does not belong to the given user_id"})
		return
	}

	// Parse request body
	var reqBody model.TransactionBank
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		logrus.Errorf("Incorrect request body: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the sender ID to the user ID
	reqBody.SenderID = uint(userID)
	reqBody.BankAccountID = uint(bankAccID)

	// Create the deposit transaction
	if err := c.txUsecase.CreateDepositBank(&reqBody); err != nil {
		logrus.Errorf("Failed to create Deposit Transaction: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	logrus.Info("Deposit Transaction created Succesfully")
	ctx.JSON(http.StatusCreated, gin.H{"message": "Deposit Transaction created Successfully"})
}

func (c *TransactionController) CreateDepositCard(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}
	defer logger.Close()
	logrus.SetOutput(logger)

	// Parse user_id parameter
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		logrus.Errorf("Invalid UserID: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID parameter"})
		return
	}
	cardID, err := strconv.Atoi(ctx.Param("card_id"))
	if err != nil {
		logrus.Errorf("Invalid CardID: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid CardID parameter"})
		return
	}
	cardAcc, err := c.cardUsecase.FindCardByCardID(uint(cardID))
	if err != nil {
		logrus.Errorf("Failed to create Deposit Transaction: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to create Deposit Transaction")
		return
	}

	if cardAcc.UserID != uint(userID) {
		logrus.Errorf("CardID doesn't belong to the given UserID: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "CardID doesn't belong to the given UsreID"})
		return
	}

	if cardAcc.UserID != uint(userID) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "card account does not belong to the given user_id"})
		return
	}

	// Parse request body
	var reqBody model.TransactionCard
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		logrus.Errorf("Incorrect request body: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the sender ID to the user ID
	reqBody.SenderID = uint(userID)
	reqBody.CardID = uint(cardID)

	// Create the deposit transaction
	if err := c.txUsecase.CreateDepositCard(&reqBody); err != nil {
		logrus.Errorf("Failed to create Deposit Transaction: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	logrus.Info("Deposit Transaction created Succesfully")
	ctx.JSON(http.StatusCreated, gin.H{"message": "Deposit transaction created successfully"})
}

func (c *TransactionController) CreateWithdrawal(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}
	defer logger.Close()
	logrus.SetOutput(logger)

	// Parse user_id parameter
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		logrus.Errorf("Invalid UserID: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID parameter"})
		return
	}
	bankAccID, err := strconv.Atoi(ctx.Param("bank_account_id"))
	if err != nil {
		logrus.Errorf("Invalid Bank AccountID: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Bank AccountID parameter"})
		return
	}

	// Retrieve bank account by bank_account_id
	bankAcc, err := c.bankUsecase.FindBankAccByAccountID(uint(bankAccID))
	if err != nil {
		logrus.Errorf("Failed to create Withdrawal Transaction: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to create Withdrawal Transaction")
		return
	}
	if bankAcc.UserID != uint(userID) {
		logrus.Errorf("Bank Account doesn't belong to the given UserID: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bank Account doesn't belong to the given UserID"})
		return
	}
	if bankAcc.UserID != uint(userID) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bank account does not belong to the given user_id"})
		return
	}

	// Parse request body
	var reqBody model.TransactionWithdraw
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		logrus.Errorf("Incorrect request body: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the sender ID to the user ID
	reqBody.SenderID = uint(userID)

	// Create the withdrawal transaction
	if err := c.txUsecase.CreateWithdrawal(&reqBody); err != nil {
		if err.Error() == "insufficient balance" {
			logrus.Errorf("Failed to create Withdrawal Transaction: %v", err)
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		} else {
			logrus.Errorf("Failed to create Withdrawal Transaction: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	logrus.Info("Withdrawal Transaction created Succesfully")
	ctx.JSON(http.StatusCreated, gin.H{"message": "Withdrawal transaction created successfully"})
}

func (c *TransactionController) CreateTransferTransaction(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}
	defer logger.Close()
	logrus.SetOutput(logger)

	// Parse transfer data from request body
	newTransfer := model.TransactionTransfer{}
	if err := ctx.BindJSON(&newTransfer); err != nil {
		logrus.Errorf("Failed to parse transfer data: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Input")
		return
	}

	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		logrus.Errorf("Invalid Input: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Input")
		return
	}

	// Get sender by ID
	sender, err := c.userUsecase.FindById(uint(userID))
	if err != nil {
		logrus.Errorf("Failed to get Sender User: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to create transfer transaction")
		return
	}

	recipient, err := c.userUsecase.FindById(newTransfer.RecipientID)
	if err != nil {
		logrus.Errorf("Failed to get Recipient User: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to create transfer transaction")
		return
	}

	// Create transfer transaction in use case layer
	result, err := c.txUsecase.CreateTransfer(sender, recipient, newTransfer.Amount)
	if err != nil {
		logrus.Errorf("Failed to create Transfer Transaction: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to create Transfer Transaction")
		return
	}

	logrus.Info("Transfer Transaction created Succesfully")
	response.JSONSuccess(ctx.Writer, http.StatusCreated, result)
}

func (c *TransactionController) CreateRedeemTransaction(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}
	defer logger.Close()
	logrus.SetOutput(logger)

	// Parse user_id from URL parameter
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		logrus.Errorf("Invalid Input: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Input")
		return
	}
	peID, err := strconv.Atoi(ctx.Param("pe_id"))
	if err != nil {

		logrus.Errorf("Invalid Input: %v", err)

		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Input")
		return
	}
	_, err = c.txUsecase.FindByPeId(uint(peID))
	if err != nil {

		logrus.Errorf("Invalid Bank: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to create Redeem Transaction")

		return
	}

	// Parse redeem data from request body
	var txData model.TransactionPoint
	if err := ctx.ShouldBindJSON(&txData); err != nil {
		logrus.Errorf("Failed to parse redeem data: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Input")
		return
	}
	txData.SenderID = uint(userID)
	txData.PointExchangeID = peID

	// Create redeem transaction in use case layer
	err = c.txUsecase.CreateRedeem(&txData)
	if err != nil {
		logrus.Errorf("Failed to create Redeem Transaction: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, err.Error())
		return
	}

	logrus.Info("Redeem Transaction created Succesfully")
	response.JSONSuccess(ctx.Writer, http.StatusCreated, "Redeem Transaction created successfully")
}

func (c *TransactionController) GetTxBySenderId(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}
	defer logger.Close()
	logrus.SetOutput(logger)

	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		logrus.Errorf("Invalid Input: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "Failed to get UserID")
		return
	}

	// Get sender by ID
	_, err = c.userUsecase.FindById(uint(userId))
	if err != nil {
		logrus.Errorf("Failed to get Sender User: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to get Sender User")
		return
	}

	txs, err := c.txUsecase.FindTxById(uint(userId))
	if err != nil {
		logrus.Errorf("Failed to get Transaction Log")
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, err)
		return
	}

	logrus.Info("Transaction Log loaded Successfully")
	response.JSONSuccess(ctx.Writer, http.StatusOK, txs)
}

func NewTransactionController(usecase usecase.TransactionUseCase, uc usecase.UserUseCase, bk usecase.BankAccUsecase, cd usecase.CardUsecase) *TransactionController {
	controller := TransactionController{
		txUsecase:   usecase,
		userUsecase: uc,
		bankUsecase: bk,
		cardUsecase: cd,
	}
	return &controller
}
