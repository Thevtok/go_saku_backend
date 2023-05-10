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

	logrus.SetOutput(logger)

	// Parse user_id parameter
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		logrus.Errorf("Invalid UserID: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "invalid user_id")
		return
	}

	// Parse bank_account_id parameter
	bankAccID, err := strconv.Atoi(ctx.Param("bank_account_id"))
	if err != nil {
		logrus.Errorf("Invalid Bank AccountID: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "invalid bank_account_id")
		return
	}

	// Retrieve bank account by bank_account_id
	bankAcc, err := c.bankUsecase.FindBankAccByAccountID(uint(bankAccID))
	if err != nil {
		logrus.Errorf("Bank_account_id not found: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Bank_account_id not found")
		return
	}

	// Check if bank account belongs to the given user_id
	if bankAcc.UserID != uint(userID) {
		logrus.Errorf("Bank Account doesn't belong to the given UserID: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Bank Account doesn't belong to the given UserID")

		return
	}

	// Parse request body
	var reqBody model.TransactionBank
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		logrus.Errorf("Incorrect request body: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Incorrect request body")
		return
	}

	// Set the sender ID to the user ID
	reqBody.SenderID = uint(userID)
	reqBody.BankAccountID = uint(bankAccID)

	// Create the deposit transaction
	if err := c.txUsecase.CreateDepositBank(&reqBody); err != nil {
		logrus.Errorf("Failed to create Deposit Transaction: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to create Deposit Transaction")
		return
	}

	// Return success response
	logrus.Info("Deposit Transaction created Succesfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusCreated, "Deposit Transaction created Succesfully")
}

func (c *TransactionController) CreateDepositCard(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	// Parse user_id parameter
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		logrus.Errorf("Invalid UserID: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "invalid user_id")
		return
	}
	cardID, err := strconv.Atoi(ctx.Param("card_id"))
	if err != nil {
		logrus.Errorf("Invalid CardID: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "invalid card_id")
		return
	}
	cardAcc, err := c.cardUsecase.FindCardByCardID(uint(cardID))
	if err != nil {
		logrus.Errorf("Card_id not found: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Card_id not found")
		return
	}

	if cardAcc.UserID != uint(userID) {
		logrus.Errorf("CardID doesn't belong to the given UserID: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "CardID doesn't belong to the given UserID")
		return
	}

	// Parse request body
	var reqBody model.TransactionCard
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		logrus.Errorf("Incorrect request body: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Incorrect request body")
		return
	}

	// Set the sender ID to the user ID
	reqBody.SenderID = uint(userID)
	reqBody.CardID = uint(cardID)

	// Create the deposit transaction
	if err := c.txUsecase.CreateDepositCard(&reqBody); err != nil {
		logrus.Errorf("Failed to create Deposit Transaction: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to create Deposit Transaction")
		return
	}

	// Return success response
	logrus.Info("Deposit Transaction created Succesfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusCreated, "Deposit Transaction created Succesfully")
}

func (c *TransactionController) CreateWithdrawal(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	// Parse user_id parameter
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		logrus.Errorf("Invalid UserID: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "invalid user_id")
		return
	}
	bankAccID, err := strconv.Atoi(ctx.Param("bank_account_id"))
	if err != nil {
		logrus.Errorf("Invalid Bank AccountID: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "invalid bank_account_id")
		return
	}

	// Retrieve bank account by bank_account_id
	bankAcc, err := c.bankUsecase.FindBankAccByAccountID(uint(bankAccID))
	if err != nil {
		logrus.Errorf("Bank_account_id not found: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Bank_account_id not found")
		return
	}
	if bankAcc.UserID != uint(userID) {
		logrus.Errorf("Bank Account doesn't belong to the given UserID: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Bank Account doesn't belong to the given UserID")
		return
	}

	// Parse request body
	var reqBody model.TransactionWithdraw
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		logrus.Errorf("Incorrect request body: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Incorrect request body")
		return
	}

	// Set the sender ID to the user ID
	reqBody.SenderID = uint(userID)

	// Create the withdrawal transaction
	if err := c.txUsecase.CreateWithdrawal(&reqBody); err != nil {
		if err.Error() == "insufficient balance" {
			logrus.Errorf("Failed to create Withdrawal Transaction: %v", err)
			response.JSONErrorResponse(ctx.Writer, false, http.StatusUnprocessableEntity, "insufficient balance")
			return
		} else {
			logrus.Errorf("Failed to create Withdrawal Transaction: %v", err)
			response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to create Withdrawal Transaction")
			return
		}
	}
	logrus.Info("Withdrawal Transaction created Succesfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusCreated, "Withdrawal Transaction created Succesfully")
}

func (c *TransactionController) CreateTransferTransaction(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	// Parse transfer data from request body
	var newTransfer model.TransactionTransferResponse
	if err := ctx.BindJSON(&newTransfer); err != nil {
		logrus.Errorf("Failed to parse transfer data: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "invalid input")
		return
	}

	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		logrus.Errorf("Invalid user_id: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "invalid user_id")
		return
	}

	// Get sender by ID
	sender, err := c.userUsecase.FindById(uint(userID))
	if err != nil {
		logrus.Errorf("Failed to get Sender User: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Failed to get Sender User")
		return
	}

	recipient, err := c.userUsecase.FindById(newTransfer.RecipientID)
	if err != nil {
		logrus.Errorf("Failed to get Recipient User: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Failed to get Recipient User")
		return
	}

	if sender.Balance < newTransfer.Amount {
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "insufficient balance")
		return
	}
	if sender.ID == recipient.ID {
		response.JSONErrorResponse(ctx.Writer, false, http.StatusForbidden, "Input the recipient correctly")
		return
	}

	// Create transfer transaction in use case layer
	result, err := c.txUsecase.CreateTransfer(sender, recipient, newTransfer.Amount)
	if err != nil {
		logrus.Errorf("Failed to create Transfer Transaction: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to create Transfer Transaction")
		return
	}

	logrus.Info("Transfer Transaction created Succesfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusCreated, result)
}
func (c *TransactionController) CreateRedeemTransaction(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		logrus.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	// Parse user_id from URL parameter
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		logrus.Errorf("Invalid user_id: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid user_id")
		return
	}
	peID, err := strconv.Atoi(ctx.Param("pe_id"))
	if err != nil {
		logrus.Errorf("Invalid pe_id: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid pe_id")
		return
	}

	user, err := c.userUsecase.FindById(uint(userID))

	if err != nil {
		logrus.Errorf("Failed to get user: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Failed to get user")
		return
	}

	pointExchange, err := c.txUsecase.FindByPeId(peID)
	if err != nil {
		logrus.Errorf("Failed to find point exchange: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Failed to find point exchange")
		return
	}

	// Parse redeem data from request body
	var txData model.TransactionPoint
	if err := ctx.ShouldBindJSON(&txData); err != nil {
		logrus.Info(txData)
		logrus.Errorf("Invalid input: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid input")
		return
	}
	price := pointExchange.Price
	if txData.Point != price {
		logrus.Info(txData)
		logrus.Errorf("Reward or price on point exchange data doesn't match with the transaction data")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Reward or price on point exchange data doesn't match with the transaction data")
		return
	}
	if user.Point < txData.Point {
		logrus.Info(txData)
		logrus.Errorf("your point is not enough to redeem")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "your point is not enough to redeem")
		return
	}

	txData.SenderID = uint(userID)
	txData.PointExchangeID = peID

	// Create redeem transaction in use case layer
	err = c.txUsecase.CreateRedeem(&txData)
	if err != nil {
		logrus.Errorf("Failed to create redeem transaction: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to create redeem transaction")
		return
	}

	logrus.Info("Redeem transaction created successfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusCreated, "Redeem transaction created successfully")
}

func (c *TransactionController) GetTxBySenderId(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		logrus.Errorf("Invalid user_id: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Failed to get UserID")
		return
	}

	// Get sender by ID
	_, err = c.userUsecase.FindById(uint(userId))
	if err != nil {
		logrus.Errorf("Failed to get Sender User: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Failed to get Sender User")
		return
	}

	txs, err := c.txUsecase.FindTxById(uint(userId), uint(userId))
	if err != nil {
		logrus.Errorf("Failed to get Transaction")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to get Transaction")
		return
	}

	logrus.Info("Transaction Log loaded Successfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, txs)
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
