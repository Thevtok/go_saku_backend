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

type CardController struct {
	cardUsecase usecase.CardUsecase
}

func (c *CardController) FindAllCard(ctx *gin.Context) {
	result := c.cardUsecase.FindAllCard()
	if result == nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to get user Card")
		return
	}

	response.JSONSuccess(ctx.Writer, http.StatusOK, result)
}

func (c *CardController) FindCardByUserID(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid user ID")
	}
	existingUser, _ := c.cardUsecase.FindCardByUserID(uint(userID))
	if existingUser == nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "Card not found")
		return
	}

	response.JSONSuccess(ctx.Writer, http.StatusOK, existingUser)
}

func (c *CardController) FindCardByCardID(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("card_id"), 10, 64)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Card ID")
		return
	}

	existingUser, _ := c.cardUsecase.FindCardByCardID(uint(userID))
	if existingUser == nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "Card ID not found")
		return
	}

	response.JSONSuccess(ctx.Writer, http.StatusOK, existingUser)
}

func (c *CardController) CreateCardID(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Failed to get user ID")
		return
	}

	var newCard model.CardResponse
	err := ctx.BindJSON(&newCard)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := c.cardUsecase.Register(userID.(uint), &newCard)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to create Card ID")
		return
	}

	response.JSONSuccess(ctx.Writer, http.StatusOK, result)
}

func (c *CardController) Edit(ctx *gin.Context) {
	cardID, err := strconv.ParseUint(ctx.Param("card_id"), 10, 64)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Card ID")
		return
	}

	existingUser, _ := c.cardUsecase.FindCardByCardID(uint(cardID))
	if existingUser == nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "Card not found")
		return
	}
	user := &model.Card{}
	if err := mapstructure.Decode(existingUser, user); err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to edit Card")
		return
	}

	if err := ctx.BindJSON(user); err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid input")
		return
	}
	updateCard := c.cardUsecase.Edit(user)

	response.JSONSuccess(ctx.Writer, http.StatusOK, updateCard)
}

func (c *CardController) UnregAll(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid user ID")
	}

	user := &model.Card{
		UserID: uint(userID),
	}
	result := c.cardUsecase.UnregALL(user)

	response.JSONSuccess(ctx.Writer, http.StatusOK, result)
}

func (c *CardController) UnregByCardId(ctx *gin.Context) {
	cardID, err := strconv.ParseUint(ctx.Param("card_id"), 10, 64)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid account ID")
		return
	}

	err = c.cardUsecase.UnregByCardID(uint(cardID))
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to delete Card ID")
		return
	}

	response.JSONSuccess(ctx.Writer, http.StatusOK, "Card ID deleted Successfully")
}

func NewCardController(u usecase.CardUsecase) *CardController {
	controller := CardController{
		cardUsecase: u,
	}

	return &controller
}
