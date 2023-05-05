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

type CardController struct {
	cardUsecase usecase.CardUsecase
}

func (c *CardController) FindAllCard(ctx *gin.Context) {
	//logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logger.Close()
	logrus.SetOutput(logger)
	result := c.cardUsecase.FindAllCard()
	if result == nil {
		logrus.Errorf("Failed to get user card: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to get user Card")
		return
	}

	logrus.Infof("Success to get user card")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, result)
}

func (c *CardController) FindCardByUserID(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	if err != nil {
		logrus.Errorf("Failed to get user id: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid user ID")
	}
	existingUser, _ := c.cardUsecase.FindCardByUserID(uint(userID))
	if existingUser == nil {
		logrus.Errorf("Failed to find card: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Card not found")
		return
	}

	logrus.Infof("Success to find card")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, existingUser)
}

func (c *CardController) FindCardByCardID(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("card_id"), 10, 64)
	if err != nil {
		logrus.Errorf("Failed to get card id: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid Card ID")
		return
	}

	existingUser, _ := c.cardUsecase.FindCardByCardID(uint(userID))
	if existingUser == nil {
		logrus.Errorf("Failed to find card id: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Card ID not found")
		return
	}

	logrus.Infof("Success to find Card ID")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, existingUser)
}

func (c *CardController) CreateCardID(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logrus.Errorf("Failed to get user ID: %v", exists)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Failed to get user ID")
		return
	}

	var newCard model.CardResponse
	err := ctx.BindJSON(&newCard)
	if err != nil {
		logrus.Errorf("Failed to request body: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := c.cardUsecase.Register(userID.(uint), &newCard)
	if err != nil {
		logrus.Errorf("failed to create card ID: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to create Card ID")
		return
	}

	logrus.Infof("Success to create card ID")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, result)
}

func (c *CardController) Edit(ctx *gin.Context) {
	cardID, err := strconv.ParseUint(ctx.Param("card_id"), 10, 64)
	if err != nil {
		logrus.Errorf("Failed to get card ID: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid Card ID")
		return
	}

	existingUser, _ := c.cardUsecase.FindCardByCardID(uint(cardID))
	if existingUser == nil {
		logrus.Errorf("Failed to find card: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Card not found")
		return
	}
	user := &model.Card{}
	if err := mapstructure.Decode(existingUser, user); err != nil {
		logrus.Errorf("Failed to edit card: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to edit Card")
		return
	}

	if err := ctx.BindJSON(user); err != nil {
		logrus.Errorf("Invalid input")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid input")
		return
	}
	updateCard := c.cardUsecase.Edit(user)

	logrus.Infof("Success to edit card")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, updateCard)
}

func (c *CardController) UnregAll(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	if err != nil {
		logrus.Errorf("Failed to get user ID: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid user ID")
	}

	user := &model.Card{
		UserID: uint(userID),
	}
	result := c.cardUsecase.UnregALL(user)

	logrus.Infof("Success to get user ID")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, result)
}

func (c *CardController) UnregByCardId(ctx *gin.Context) {
	cardID, err := strconv.ParseUint(ctx.Param("card_id"), 10, 64)
	if err != nil {
		logrus.Errorf("Failed to get valid account ID: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid account ID")
		return
	}

	err = c.cardUsecase.UnregByCardID(uint(cardID))
	if err != nil {
		logrus.Errorf("Failed to delete card ID: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to delete Card ID")
		return
	}

	logrus.Infof("Success to delete card ID")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, "Card ID deleted Successfully")
}

func NewCardController(u usecase.CardUsecase) *CardController {
	controller := CardController{
		cardUsecase: u,
	}

	return &controller
}
