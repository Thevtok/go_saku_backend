package controller

import (
	"net/http"
	"strconv"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/model/response"
	"github.com/ReygaFitra/inc-final-project.git/usecase"
	"github.com/gin-gonic/gin"
)

type RewardController struct {
	rewardUsecase usecase.RewardUseCase
}

func (c *RewardController) FindPoints(ctx *gin.Context) {
	res := c.rewardUsecase.FindPoints()
	if res == nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusInternalServerError, "Failed to get user point")
		return
	}
	response.JSONSuccess(ctx.Writer, http.StatusOK, res)
}

func (c *RewardController) FindPointByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid point ID")
		return
	}
	res := c.rewardUsecase.FindPointByID(uint(id))
	if res == nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusNotFound, "Point not found")
		return
	}
	response.JSONSuccess(ctx.Writer, http.StatusOK, res)
}

func (c *RewardController) Register(ctx *gin.Context) {
	var newUser model.Reward
	if err := ctx.BindJSON(&newUser); err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Input")
		return
	}
	res := c.rewardUsecase.Register(&newUser)
	response.JSONSuccess(ctx.Writer, http.StatusCreated, res)
}

func (c *RewardController) Edit(ctx *gin.Context) {
	var user model.Reward
	if err := ctx.BindJSON(&user); err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid Input")
		return
	}
	res := c.rewardUsecase.Edit(&user)
	response.JSONSuccess(ctx.Writer, http.StatusOK, res)
}

func (c *RewardController) Unreg(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.JSONErrorResponse(ctx.Writer, http.StatusBadRequest, "Invalid point ID")
		return
	}
	res := c.rewardUsecase.Unreg(uint(id))
	response.JSONSuccess(ctx.Writer, http.StatusOK, res)
}

func NewRewardController(u usecase.RewardUseCase) *RewardController {
	controller := RewardController{
		rewardUsecase: u,
	}
	return &controller
}