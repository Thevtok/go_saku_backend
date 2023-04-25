package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/model/response"
	"github.com/ReygaFitra/inc-final-project.git/usecase"
	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	txUsecase   usecase.TransactionUseCase
	userUsecase usecase.UserUseCase
}

func (c *TransactionController) CreateTransferTransaction(ctx *gin.Context) {
	// Parse transfer data from request body
	newTransfer := model.TransactionTransfer{}
	if err := ctx.BindJSON(&newTransfer); err != nil {
		log.Printf("Failed to parse transfer data: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Input")
		return
	}

	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		log.Printf("Failed to convert user_id to uint: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Input")
		return
	}

	// Get sender by ID
	sender, err := c.userUsecase.FindById(uint(userID))
	if err != nil {
		log.Printf("Failed to get sender user: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to create transfer transaction")
		return
	}

	recipient, err := c.userUsecase.FindById(newTransfer.RecipientID)
	if err != nil {
		log.Printf("Failed to get recipient user: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to create transfer transaction")
		return
	}

	// Create transfer transaction in use case layer
	result, err := c.txUsecase.CreateTransfer(sender, recipient, newTransfer.Amount)
	if err != nil {
		log.Printf("Failed to create transfer transaction: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to create transfer transaction")
		return
	}

	response.JSONSuccess(ctx.Writer, http.StatusCreated, result)
}

func (c *TransactionController) CreateDepositBank(ctx *gin.Context) {
	// Parse user_id parameter
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid user_id parameter"})
		return
	}

	// Parse request body
	var reqBody model.TransactionBank
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the sender ID to the user ID
	reqBody.SenderID = uint(userID)

	// Create the deposit transaction
	if err := c.txUsecase.CreateDepositBank(&reqBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	ctx.JSON(http.StatusCreated, gin.H{"message": "Deposit transaction created successfully"})
}

func (c *TransactionController) CreateDepositCard(ctx *gin.Context) {
	// Parse user_id parameter
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid user_id parameter"})
		return
	}

	// Parse request body
	var reqBody model.TransactionCard
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the sender ID to the user ID
	reqBody.SenderID = uint(userID)

	// Create the deposit transaction
	if err := c.txUsecase.CreateDepositCard(&reqBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	ctx.JSON(http.StatusCreated, gin.H{"message": "Deposit transaction created successfully"})
}

func (c *TransactionController) CreateWithdrawal(ctx *gin.Context) {
	// Parse user_id parameter
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid user_id parameter"})
		return
	}

	// Parse request body
	var reqBody model.TransactionWithdraw
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the sender ID to the user ID
	reqBody.SenderID = uint(userID)

	// Create the withdrawal transaction
	if err := c.txUsecase.CreateWithdrawal(&reqBody); err != nil {
		if err.Error() == "insufficient balance" {
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (c *TransactionController) CreateRedeemTransaction(ctx *gin.Context) {
	// Parse user_id from URL parameter
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		log.Printf("Failed to convert user_id to uint: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Input")
		return
	}

	// Parse redeem data from request body
	var txData model.TransactionPoint
	if err := ctx.ShouldBindJSON(&txData); err != nil {
		log.Printf("Failed to parse redeem data: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Input")
		return
	}
	txData.SenderID = uint(userID)

	// Create redeem transaction in use case layer
	err = c.txUsecase.CreateRedeem(&txData)
	if err != nil {
		log.Printf("Failed to create redeem transaction: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, err.Error())
		return
	}
	response.JSONSuccess(ctx.Writer, http.StatusCreated, "Redeem transaction created successfully")
}
func (c *TransactionController) GetTxBySenderId(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "Failed to get user_id")
		return
	}

	// Get sender by ID
	_, err = c.userUsecase.FindById(uint(userId))
	if err != nil {
		log.Printf("Failed to get sender user: %v", err)
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to create transfer transaction")
		return
	}

	txs, err := c.txUsecase.FindTxById(uint(userId))
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, err)
		return
	}

	response.JSONSuccess(ctx.Writer, http.StatusOK, txs)
}

func NewTransactionController(usecase usecase.TransactionUseCase, uc usecase.UserUseCase) *TransactionController {
	controller := TransactionController{
		txUsecase:   usecase,
		userUsecase: uc,
	}
	return &controller
}
