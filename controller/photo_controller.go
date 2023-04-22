package controller

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/model/response"
	"github.com/ReygaFitra/inc-final-project.git/usecase"
	"github.com/gin-gonic/gin"
)

type PhotoController struct {
	photoUsecase usecase.PhotoUsecase
}

func (c *PhotoController) Upload(ctx *gin.Context) {
	userName := ctx.PostForm("username")
    if userName == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username"})
        return
    }
    file, header, err := ctx.Request.FormFile("photo")
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    defer file.Close()
	 // Validasi ekstensi file hanya png
	 ext := filepath.Ext(header.Filename)
	 if ext != ".png" {
		 ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only PNG files are allowed"})
		 return
	 }
    err = c.photoUsecase.Upload(ctx.Request.Context(), userName, file, header)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

func (c *PhotoController) Download(ctx *gin.Context) {
	userName := ctx.Param("username")
	if userName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid username"})
		return 
	}
	photo, err := c.photoUsecase.Download(userName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	file, err := os.Open(photo.Url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer file.Close()
	
    ctx.Header("Content-Type", "image/png")
	ctx.File(photo.Url)
	
}

func (c *PhotoController) Edit(ctx *gin.Context) {
	userName := ctx.Param("username")
	if userName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid username"})
		return 
	}
	file, header, err := ctx.Request.FormFile("photo")
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
	defer file.Close()
	// Validasi ekstensi file hanya png
	ext := filepath.Ext(header.Filename)
	if ext != ".png" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only PNG files are allowed"})
		return
	}
	err = c.photoUsecase.Edit(&model.PhotoUrl{}, userName, file, header)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Photo Update Succesfully"})
}

func (c *PhotoController) Remove(ctx *gin.Context) {
	userName := ctx.Param("username")
	if userName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid username"})
		return 
	}
	res := c.photoUsecase.Remove(userName)
	response.JSONSuccess(ctx.Writer, http.StatusOK, res)
}

func NewPhotoController(u usecase.PhotoUsecase) *PhotoController {
	controller := PhotoController{
		photoUsecase: u,
	}
	return &controller
}