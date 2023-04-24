package controller

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/model/response"
	"github.com/ReygaFitra/inc-final-project.git/usecase"
	"github.com/gin-gonic/gin"
)

type PhotoController struct {
	photoUsecase usecase.PhotoUsecase
}

func (c *PhotoController) Upload(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.PostForm("user_id")) 
	if err != nil {
		log.Printf("Failed to get user id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}
	file, header, err := ctx.Request.FormFile("photo")
	if err != nil {
		log.Printf("Something went wrong at Form File Key: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()
	// Validasi ekstensi file
	ext := filepath.Ext(header.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		log.Printf("Extension file is not image file: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only Image files are allowed"})
		return
	}
	err = c.photoUsecase.Upload(ctx.Request.Context(), uint(userID), file, header)
	if err != nil {
		log.Printf("Something went wrong when uploading file: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

func (c *PhotoController) Download(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("user_id")) 
	if err != nil {
		log.Printf("Failed to get user id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user_id"})
		return
	}
	photo, err := c.photoUsecase.Download(uint(userID))
	if err != nil {
		log.Printf("Something went wrong when downloading file: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	// Validasi extensi file
	if filepath.Ext(photo.Url) != ".png" && filepath.Ext(photo.Url) != ".jpg" && filepath.Ext(photo.Url) != ".jpeg" {
		log.Printf("Extension file is not image file: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only Image files are allowed"})
		return
	}
	file, err := os.Open(photo.Url)
	if err != nil {
		log.Printf("Failed to get photo: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer file.Close()
	// Set response header sesuai format file
	contentType := "image/png"
	if filepath.Ext(photo.Url) == ".jpg" || filepath.Ext(photo.Url) == ".jpeg" {
		contentType = "image/jpeg"
	}
	ctx.Header("Content-Type", contentType)
	ctx.File(photo.Url)

}

func (c *PhotoController) Edit(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("user_id")) 
	if err != nil {
		log.Printf("Failed to get user id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user_id"})
		return
	}
	file, header, err := ctx.Request.FormFile("photo")
	if err != nil {
		log.Printf("Something went wrong at Form File Key: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()
	// Validasi ekstensi file hanya png
	ext := filepath.Ext(header.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		log.Printf("Extension file is not image file: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only Image files are allowed"})
		return
	}
	err = c.photoUsecase.Edit(&model.PhotoUrl{}, uint(userID), file, header)
	if err != nil {
		log.Printf("Something went wrong when editing file: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Photo Update Succesfully"})
}

func (c *PhotoController) Remove(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("user_id")) 
	if err != nil {
		log.Printf("Failed to get user id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user_id"})
		return
	}
	res := c.photoUsecase.Remove(uint(userID))
	response.JSONSuccess(ctx.Writer, http.StatusOK, res)
}

func NewPhotoController(u usecase.PhotoUsecase) *PhotoController {
	controller := PhotoController{
		photoUsecase: u,
	}
	return &controller
}
