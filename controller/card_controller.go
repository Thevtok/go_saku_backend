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

func (c *CardController) FindCardByUsername(ctx *gin.Context) {
	username := ctx.Param("username")
	existingUser, _ := c.cardUsecase.FindCardByUsername(username)
	if existingUser == nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "Card not found")
		return
	}

	response.JSONSuccess(ctx.Writer, http.StatusOK, existingUser)
}

func (c *CardController) FindCardByCardID(ctx *gin.Context) {
	user_id_str := ctx.Param("card_id")
	user_id, err := strconv.ParseUint(user_id_str, 10, 64)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Card ID")
		return
	}

	existingUser, _ := c.cardUsecase.FindCardByCardID(uint(user_id))
	if existingUser == nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "Bank not found")
		return
	}

	response.JSONSuccess(ctx.Writer, http.StatusOK, existingUser)
}

func (c *CardController) CreateCardID(ctx *gin.Context) {
	username := ctx.GetString("username")
	var newCardID model.CardResponse
	err := ctx.BindJSON(&newCardID)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := c.cardUsecase.Register(username, &newCardID)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to create Card")
		return
	}

	response.JSONSuccess(ctx.Writer, http.StatusOK, result)
}

func (c *CardController) Edit(ctx *gin.Context) {
	card_id_str := ctx.Param("card_id")
	card_id, err := strconv.ParseUint(card_id_str, 10, 64)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Card ID")
		return
	}

	existingUser, _ := c.cardUsecase.FindCardByCardID(uint(card_id))
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
	username := ctx.Param("username")
	user := &model.Card{
		Username: username,
	}
	result := c.cardUsecase.UnregALL(user)

	response.JSONSuccess(ctx.Writer, http.StatusOK, result)
}

func (c *CardController) UnregByCardId(ctx *gin.Context) {
	cardIDStr := ctx.Param("card_id")
	cardID, err := strconv.ParseUint(cardIDStr, 10, 64)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid account ID")
		return
	}

	err = c.cardUsecase.UnregByCardId(uint(cardID))
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
