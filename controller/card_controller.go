package controller

import (
	"net/http"
	"strconv"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/usecase"
	"github.com/gin-gonic/gin"
)

type CardController struct {
	usecase usecase.CardUsecase
}

func (c *CardController) FindCardById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid card ID")
		return
	}

	res, err := c.usecase.FindByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *CardController) FindCardsByUserID(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid user ID")
		return
	}

	res, err := c.usecase.FindByUserID(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *CardController) Register(ctx *gin.Context) {
	var newCard model.Card

	if err := ctx.BindJSON(&newCard); err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid Input")
		return
	}

	err := c.usecase.Register(&newCard)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, newCard)
}

func (c *CardController) Edit(ctx *gin.Context) {
	var card model.Card

	if err := ctx.BindJSON(&card); err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid Input")
		return
	}

	err := c.usecase.Edit(&card)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, card)
}

func (c *CardController) Unregister(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid card ID")
		return
	}

	card, err := c.usecase.FindByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, err.Error())
		return
	}

	err = c.usecase.Unreg(card)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, card)
}

func NewCardController(u usecase.CardUsecase) *CardController {
	controller := CardController{
		usecase: u,
	}

	return &controller
}
