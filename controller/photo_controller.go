package controller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/model/response"
	"github.com/ReygaFitra/inc-final-project.git/usecase"
	"github.com/ReygaFitra/inc-final-project.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PhotoController struct {
	photoUsecase usecase.PhotoUsecase
}

func (c *PhotoController) Upload(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logger.Close()
	logrus.SetOutput(logger)
	// Body Form data user_id
	userID, err := strconv.Atoi(ctx.PostForm("user_id"))
	if err != nil {
		logrus.Errorf("Failed to get user id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}
	// Body Form data File
	file, err := ctx.FormFile("photo")
	if err != nil {
		logrus.Errorf("Failed to get file from request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
		return
	}
	// url photo location
	filename := file.Filename
	path := fmt.Sprintf(utils.DotEnv("FILE_LOCATION"), filename)
	out, err := os.Create(path)
	if err != nil {
		logrus.Errorf("Failed to create file: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer out.Close()
	// Validasi ekstensi file
	ext := filepath.Ext(filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		logrus.Errorf("Extension file is not image file: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only Image files are allowed"})
		return
	}
	fileIn, err := file.Open()
	if err != nil {
		logrus.Errorf("Failed to open file: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer fileIn.Close()
	_, err = io.Copy(out, fileIn)
	if err != nil {
		logrus.Errorf("Failed to write file: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file"})
		return
	}
	// Simpan informasi file ke database
	photo := &model.PhotoUrl{
		UserID: uint(userID),
		Url:    path,
	}
	err = c.photoUsecase.Upload(photo)
	if err != nil {
		logrus.Errorf("Failed to upload photo: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload photo"})
		return
	}
	logrus.Info("Photo uploaded succesfully")
	ctx.JSON(http.StatusCreated, gin.H{"message": "photo uploaded successfully"})
}

func (c *PhotoController) Download(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logger.Close()
	logrus.SetOutput(logger)
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		logrus.Errorf("Failed to get user id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user_id"})
		return
	}
	photo, err := c.photoUsecase.Download(uint(userID))
	if err != nil {
		logrus.Errorf("Something went wrong when downloading file: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	// Validasi extensi file
	if filepath.Ext(photo.Url) != ".png" && filepath.Ext(photo.Url) != ".jpg" && filepath.Ext(photo.Url) != ".jpeg" {
		logrus.Errorf("Extension file is not image file: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only Image files are allowed"})
		return
	}
	file, err := os.Open(photo.Url)
	if err != nil {
		logrus.Errorf("Failed to get photo: %v", err)
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
	logrus.Info("Photo get succesfully")
	ctx.JSON(http.StatusOK, gin.H{"message": "Photo get successfully"})
}

func (c *PhotoController) Edit(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logger.Close()
	logrus.SetOutput(logger)
	// Body Form data user_id
	userID, err := strconv.Atoi(ctx.PostForm("user_id"))
	if err != nil {
		logrus.Errorf("Failed to get user id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}
	// Body Form data File
	file, err := ctx.FormFile("photo")
	if err != nil {
		logrus.Errorf("Failed to get file from request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
		return
	}
	// url photo location
	filename := file.Filename
	path := fmt.Sprintf(utils.DotEnv("FILE_LOCATION"), filename)
	out, err := os.Create(path)
	if err != nil {
		logrus.Errorf("Failed to create file: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer out.Close()
	// Validasi ekstensi file
	ext := filepath.Ext(filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		logrus.Errorf("Extension file is not image file: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only Image files are allowed"})
		return
	}
	fileIn, err := file.Open()
	if err != nil {
		logrus.Errorf("Failed to open file: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer fileIn.Close()
	_, err = io.Copy(out, fileIn)
	if err != nil {
		logrus.Errorf("Failed to write file: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file"})
		return
	}
	// Simpan informasi file ke database
	photo := &model.PhotoUrl{
		UserID: uint(userID),
		Url:    path,
	}
	err = c.photoUsecase.Edit(photo)
	if err != nil {
		logrus.Errorf("Failed to upload photo: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload photo"})
		return
	}
	logrus.Info("Photo Edit succesfully")
	ctx.JSON(http.StatusOK, gin.H{"message": "photo Edit successfully"})
}

func (c *PhotoController) Remove(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logger.Close()
	logrus.SetOutput(logger)
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		logrus.Errorf("Failed to get user id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user_id"})
		return
	}
	res := c.photoUsecase.Remove(uint(userID))
	logrus.Info("Remove Photo Succesfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, res)
}

func NewPhotoController(u usecase.PhotoUsecase) *PhotoController {
	controller := PhotoController{
		photoUsecase: u,
	}
	return &controller
}
