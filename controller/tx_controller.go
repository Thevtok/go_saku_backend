package controller

import (
	"encoding/json"
	"io"
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

var depositToken string
var userIDdepo string
var depoAmount int
var deviceToken string

type TransactionController struct {
	txUsecase   usecase.TransactionUseCase
	userUsecase usecase.UserUseCase
	bankUsecase usecase.BankAccUsecase
}

func (c *TransactionController) HandlePaymentNotification(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	if ctx.Request.Method != http.MethodPost {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}
	defer ctx.Request.Body.Close()

	logrus.Println("Received payment notification:")
	logrus.Println(string(body))

	// Mengambil nomor virtual account (VA)
	var notification response.PaymentNotification
	err = json.Unmarshal(body, &notification)
	if err != nil {
		logrus.Errorf("Failed to decode notification payload: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to decode notification payload"})
		return
	}

	if notification.TransactionStatus == "settlement" {
		if len(notification.VANumbers) > 0 {
			vaNumber := notification.VANumbers[0].VANumber

			err = c.userUsecase.UpdateBalance(userIDdepo, depoAmount)
			if err != nil {
				logrus.Errorf("Failed to update balance user: %v", err)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance user"})
				return
			}
			logrus.Infof("useridDEPO: %s", userIDdepo)
			logrus.Infof("depoAmount: %d", depoAmount)

			err = c.txUsecase.UpdateDepositStatus(vaNumber, depositToken)
			if err != nil {
				logrus.Errorf("Failed to update deposit status: %v", err)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update deposit status"})
				return
			}
			amount := float64(depoAmount) / 1000                               //
			formattedAmount := "Rp " + strconv.FormatFloat(amount, 'f', 3, 64) //

			err = model.SendFCMNotification(deviceToken, "Deposit Berhasil", "Anda telah melakukan deposit sebesar "+formattedAmount)

			if err != nil {
				logrus.Errorf("failed to send FCM notification: %v", err)

				return
			}

			logrus.Infof("Deposit status updated for VA Number: %s", vaNumber)

		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Notification received"})
}

func (c *TransactionController) CreateDepositBank(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	// Parse user_id parameter
	userID := ctx.Param("user_id")

	// Parse bank_account_id parameter
	bankAccID, err := strconv.Atoi(ctx.Param("bank_account_id"))
	if err != nil {
		logrus.Errorf("Invalid Bank AccountID: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "invalid bank_account_id")
		return
	}
	user, err := c.userUsecase.FindByiDToken(userID)
	if err != nil {
		logrus.Errorf("user_id not found: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "user_id not found")
		return
	}
	// Retrieve bank account by bank_account_id
	bankAcc, err := c.bankUsecase.FindBankAccByAccountID(uint(bankAccID))
	if err != nil {
		logrus.Errorf("Bank_account_id not found: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Bank_account_id not found")
		return
	}

	// Parse request body
	var reqBody model.Deposit
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		logrus.Errorf("Incorrect request body: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Incorrect request body")
		return
	}
	depoAmount = reqBody.Amount
	userIDdepo = userID
	reqBody.UserID = userID
	reqBody.BankName = bankAcc.BankName
	reqBody.AccountHolderName = bankAcc.AccountHolderName
	reqBody.AccountNumber = bankAcc.AccountNumber

	if reqBody.Amount < 10000 {
		logrus.Errorf("Minimum deposit 10.000: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Minimum deposit 10.000")
		return
	}
	token, err := model.CreateMidtransTransactionFromDeposit(&reqBody, user)
	if err != nil {
		logrus.Errorf("Failed to create Midtrans transaction: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to create Midtrans transaction")
		return
	}

	reqBody.Token = token
	depositToken = token
	deviceToken = user.Token

	// Create the deposit transaction
	if err := c.txUsecase.CreateDepositBank(&reqBody); err != nil {
		logrus.Errorf("Failed to create Deposit Transaction: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to create Deposit Transaction")
		return
	}

	amount := float64(reqBody.Amount) / 1000                           //
	formattedAmount := "Rp " + strconv.FormatFloat(amount, 'f', 3, 64) //

	err = model.SendFCMNotification(user.Token, "Deposit Pending", "Silahkan selesaikan transaksi anda terlebih dahulu sebesar "+formattedAmount)

	if err != nil {
		logrus.Errorf("failed to send FCM notification: %v", err)

		return
	}

	// Kirim respons sukses
	logrus.Info("Deposit Transaction created Successfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusCreated, gin.H{
		"body_midtrans": token,
	})

}

func (c *TransactionController) CreateWithdrawal(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	// Parse user_id parameter
	userID := ctx.Param("user_id")

	bankAccID, err := strconv.Atoi(ctx.Param("bank_account_id"))
	if err != nil {
		logrus.Errorf("Invalid Bank AccountID: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "invalid bank_account_id")
		return
	}
	user, err := c.userUsecase.FindByiDToken(userID)
	if err != nil {
		logrus.Errorf("user_id not found: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "user_id not found")
		return
	}

	// Retrieve bank account by bank_account_id
	bankAcc, err := c.bankUsecase.FindBankAccByAccountID(uint(bankAccID))
	if err != nil {
		logrus.Errorf("Bank_account_id not found: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Bank_account_id not found")
		return
	}
	if bankAcc.UserID != userID {
		logrus.Errorf("Bank Account doesn't belong to the given UserID: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Bank Account doesn't belong to the given UserID")
		return
	}

	// Parse request body
	var reqBody model.Withdraw
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		logrus.Errorf("Incorrect request body: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Incorrect request body")
		return
	}
	if reqBody.Amount < 20000 {
		logrus.Errorf("Minimum withdraw 20.000: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Minimum withdraw 20.000")
		return
	}

	// Set the sender ID to the user ID
	reqBody.UserID = userID
	reqBody.BankName = bankAcc.BankName
	reqBody.AccountHolderName = bankAcc.AccountHolderName
	reqBody.AccountNumber = bankAcc.AccountNumber

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
	amount := float64(reqBody.Amount) / 1000                           // Mengonversi nilai amount ke dalam format yang diinginkan
	formattedAmount := "Rp " + strconv.FormatFloat(amount, 'f', 3, 64) // Mengformat nilai amount menjadi format mata uang Rupiah dengan 3 digit di belakang koma

	err = model.SendFCMNotification(user.Token, "Withdraw Berhasil", "Anda telah menarik uang sebesar "+formattedAmount)

	if err != nil {
		logrus.Errorf("failed to send FCM notification: %v", err)

		return
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
	var newTransfer model.Transfer
	if err := ctx.BindJSON(&newTransfer); err != nil {
		logrus.Errorf("Failed to parse transfer data: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Failed to parse transfer data: invalid JSON format")
		return
	}

	userID := ctx.Param("user_id")

	// Get sender by ID
	sender, err := c.userUsecase.FindByiDToken(userID)
	if err != nil {
		logrus.Errorf("Failed to get Sender User: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Failed to get Sender User")
		return
	}

	recipient, err := c.userUsecase.FindByPhone(newTransfer.RecipientPhoneNumber)
	if err != nil {
		logrus.Errorf("Failed to get Recipient User: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Failed to get Recipient User")
		return
	}

	if sender.Balance < newTransfer.Amount {
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Insufficient balance")
		return
	}
	if sender.Phone_Number == recipient.Phone_Number {
		response.JSONErrorResponse(ctx.Writer, false, http.StatusForbidden, "Input the recipient correctly")
		return
	}

	if newTransfer.Amount < 10000 {
		logrus.Errorf("Minimum transfer amount is 10,000")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Minimum transfer amount is 10,000")
		return
	}
	newTransfer.SenderName = sender.Name
	newTransfer.RecipientName = recipient.Name
	newTransfer.SenderPhoneNumber = sender.Phone_Number
	newTransfer.RecipientPhoneNumber = recipient.Phone_Number
	// Create transfer transaction in use case layer
	err = c.txUsecase.CreateTransfer(sender, recipient, newTransfer.Amount)
	logrus.Info("Processing transfer transaction...")
	logrus.Infof("Sender: %s, Recipient: %s, Amount: %d", sender.Name, recipient.Name, newTransfer.Amount)

	if err != nil {
		logrus.Errorf("Failed to create Transfer Transaction: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to create Transfer Transaction")
		return
	}
	err = c.txUsecase.AssignBadge(sender)
	if err != nil {
		logrus.Errorf("Failed to assign badge: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to asign badge")
	}
	if sender.BadgeID == 2 {
		err = model.SendFCMNotification(sender.Token, "Selamat", "Anda telah naik level menjadi Silver ")

		if err != nil {
			logrus.Errorf("failed to send FCM notification: %v", err)

			return
		}
	}
	if sender.BadgeID == 3 {
		err = model.SendFCMNotification(sender.Token, "Selamat", "Anda telah naik level menjadi Gold ")

		if err != nil {
			logrus.Errorf("failed to send FCM notification: %v", err)

			return
		}
	}
	if sender.BadgeID == 4 {
		err = model.SendFCMNotification(sender.Token, "Selamat", "Anda telah naik level menjadi Platinum ")

		if err != nil {
			logrus.Errorf("failed to send FCM notification: %v", err)

			return
		}
	}
	if sender.BadgeID == 5 {
		err = model.SendFCMNotification(sender.Token, "Selamat", "Anda telah naik level menjadi Diamond ")

		if err != nil {
			logrus.Errorf("failed to send FCM notification: %v", err)

			return
		}
	}

	amount := float64(newTransfer.Amount) / 1000
	formattedAmount := "Rp " + strconv.FormatFloat(amount, 'f', 3, 64)

	err = model.SendFCMNotification(sender.Token, "Transfer Berhasil", "Anda telah mengirim uang ke "+recipient.Name+" sebesar "+formattedAmount)

	if err != nil {
		logrus.Errorf("failed to send FCM notification: %v", err)

		return
	}
	err = model.SendFCMNotification(recipient.Token, "Receive Berhasil", "Anda telah menerima uang dari "+sender.Name+" sebesar "+formattedAmount)
	logrus.Info(recipient.Token)
	if err != nil {
		logrus.Errorf("failed to send FCM notification: %v", err)

		return
	}

	logrus.Info("Transfer Transaction created Successfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusCreated, "Transfer Successfully")
}

func (c *TransactionController) CreateRedeemTransaction(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		logrus.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	// Parse user_id from URL parameter
	userID := ctx.Param("user_id")

	peID, err := strconv.Atoi(ctx.Param("pe_id"))
	if err != nil {
		logrus.Errorf("Invalid pe_id: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid pe_id")
		return
	}

	user, err := c.userUsecase.FindById(userID)

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
	var txData model.Redeem
	if err := ctx.ShouldBindJSON(&txData); err != nil {
		logrus.Info(txData)
		logrus.Errorf("Invalid input: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid input")
		return
	}
	price := pointExchange.Price
	if txData.Amount != price {
		logrus.Info(txData)
		logrus.Errorf("Reward or price on point exchange data doesn't match with the transaction data")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Reward or price on point exchange data doesn't match with the transaction data")
		return
	}
	if user.Point < txData.Amount {
		logrus.Info(txData)
		logrus.Errorf("your point is not enough to redeem")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "your point is not enough to redeem")
		return
	}

	txData.UserID = userID
	txData.PEID = peID

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

	userId := ctx.Param("user_id")

	// Get sender by ID
	_, err = c.userUsecase.FindById(userId)
	if err != nil {
		logrus.Errorf("Failed to get Sender User: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Failed to get Sender User")
		return
	}

	txs, err := c.txUsecase.FindTxById(userId)
	if err != nil {
		logrus.Errorf("Failed to get Transaction %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to get Transaction")
		return
	}
	if len(txs) == 0 {
		logrus.Errorf("Transaction not found")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Transaction not found")
		return
	}

	logrus.Info("Transaction Log loaded Successfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, txs)
}

func NewTransactionController(usecase usecase.TransactionUseCase, uc usecase.UserUseCase, bk usecase.BankAccUsecase) *TransactionController {
	controller := TransactionController{
		txUsecase:   usecase,
		userUsecase: uc,
		bankUsecase: bk,
	}
	return &controller
}
