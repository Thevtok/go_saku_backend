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

func NewTransactionController(usecase usecase.TransactionUseCase, uc usecase.UserUseCase) *TransactionController {
	con := TransactionController{
		txUsecase:   usecase,
		userUsecase: uc,
	}
	return &con
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
