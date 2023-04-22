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

func (c *CardController) FindCardByID(ctx *gin.Context) {
	user_id_str := ctx.Param("user_id")
	user_id, err := strconv.ParseUint(user_id_str, 10, 64)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid user ID")
		return
	}

	existingUser, _ := c.cardUsecase.FindCardByID(uint(user_id))
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
	userID := ctx.GetUint("user_id")

	var newCardID model.CardResponse
	err := ctx.BindJSON(&newCardID)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := c.cardUsecase.Register(userID, &newCardID)
	if err != nil {
		response.JSONSuccess(ctx.Writer, http.StatusOK, result)
	}
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
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid User ID")
		return
	}
	user := &model.Card{
		UserID: uint(userID),
	}

	result := c.cardUsecase.UnregALL(user)
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, err.Error())
		return
	}

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
