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

	logrus.SetOutput(logger)
	// Body Form data user_id
	userID, err := strconv.Atoi(ctx.PostForm("user_id"))
	if err != nil {
		logrus.Errorf("Failed to get user id: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Failed to get user id")
		return
	}
	// Body Form data File
	file, err := ctx.FormFile("photo")
	if err != nil {
		logrus.Errorf("Failed to get file from request: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Failed to get file from request")
		return
	}
	// url photo location
	filename := file.Filename
	path := fmt.Sprintf(utils.DotEnv("FILE_LOCATION"), filename)
    // Cek apakah file dengan nama yang sama sudah ada
    if _, err := os.Stat(path); err == nil {
		logrus.Errorf("File with the same name already exists: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "File with the same name already exists")
		return
	}
	fileIn, err := file.Open()
	if err != nil {
		logrus.Errorf("Failed to open file: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to open file")
		return
	}
	defer fileIn.Close()
	// Validasi ekstensi file
	ext := filepath.Ext(filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		logrus.Errorf("Extension file is not image file: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Extension file is not image file")
		return
	}
	out, err := os.Create(path)
	if err != nil {
		logrus.Errorf("Failed to create file: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to create file")
		return
	}
	defer out.Close()
	_, err = io.Copy(out, fileIn)
	if err != nil {
		logrus.Errorf("Failed to write file: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to write file")
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
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to upload photo")
		return
	}
	logrus.Info("Photo uploaded succesfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusCreated, "Photo upload succesfully")
}

func (c *PhotoController) Download(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		logrus.Errorf("Failed to get user id: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Failed to get user id")
		return
	}
	photo, err := c.photoUsecase.Download(uint(userID))
	if err != nil {
		logrus.Errorf("Something went wrong when downloading file: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Something went wrong when downloading file")
		return
	}
	file, err := os.Open(photo.Url)
	if err != nil {
		logrus.Errorf("Failed to get photo: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to get photo")
		return
	}
	defer file.Close()
	// Validasi extensi file
	if filepath.Ext(photo.Url) != ".png" && filepath.Ext(photo.Url) != ".jpg" && filepath.Ext(photo.Url) != ".jpeg" {
		logrus.Errorf("Extension file is not image file: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Extension file is not image file")
		return
	}
	// Set response header sesuai format file
	contentType := "image/png"
	if filepath.Ext(photo.Url) == ".jpg" || filepath.Ext(photo.Url) == ".jpeg" {
		contentType = "image/jpeg"
	}
	ctx.Header("Content-Type", contentType)
	ctx.File(photo.Url)
	logrus.Info("Photo get succesfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, "Photo get succesfully")
}

func (c *PhotoController) Edit(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	// Body Form data user_id
	userID, err := strconv.Atoi(ctx.PostForm("user_id"))
	if err != nil {
		logrus.Errorf("Failed to get user id: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Failed to get user id")
		return
	}
	// Body Form data File
	file, err := ctx.FormFile("photo")
	if err != nil {
		logrus.Errorf("Failed to get file from request: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Failed to get file from request")
		return
	}
	// Validasi ekstensi file
	ext := filepath.Ext(file.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		logrus.Errorf("Extension file is not image file: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Extension file is not image file")
		return
	}
	// url photo location
	filename := file.Filename
	path := fmt.Sprintf(utils.DotEnv("FILE_LOCATION"), filename)
	out, err := os.Create(path)
	if err != nil {
		logrus.Errorf("Failed to edit file: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to edit file")
		return
	}
	defer out.Close()
	fileIn, err := file.Open()
	if err != nil {
		logrus.Errorf("Failed to open file: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to open file")
		return
	}
	defer fileIn.Close()
	
	
	// Simpan informasi file ke database
	oldPhoto, err := c.photoUsecase.Download(uint(userID))
	if err != nil {
		logrus.Errorf("Failed to download photo: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to download photo")
		return
	}
	photo := &model.PhotoUrl{
		UserID: uint(userID),
		Url:    path,
	}
	err = c.photoUsecase.Edit(photo)
	if err != nil {
		logrus.Errorf("Failed to upload photo: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to upload photo")
		return
	}
	if oldPhoto != nil {
		if err = os.Remove(oldPhoto.Url); err != nil {
			logrus.Errorf("Failed to delete file photo: %v", err)
		}
	}
	_, err = io.Copy(out, fileIn)
	if err != nil {
		logrus.Errorf("Failed to write file: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to write file")
		return
	}
	logrus.Info("Photo Edit succesfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, "Photo Edit succesfully")
}

func (c *PhotoController) Remove(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		logrus.Errorf("Failed to get user id: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Failed to get user id")
		return
	}

	photo, err := c.photoUsecase.Download(uint(userID))
	if err != nil {
		logrus.Errorf("Something went wrong when downloading file: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Something went wrong when downloading file")
		return
	}
	err = os.Remove(photo.Url)
	if err != nil {
		logrus.Errorf("Failed to get photo: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to remove photo")
		return
	}
	res := c.photoUsecase.Remove(photo.UserID)
	logrus.Info("Remove Photo Succesfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, res)
}

func NewPhotoController(u usecase.PhotoUsecase) *PhotoController {
	controller := PhotoController{
		photoUsecase: u,
	}
	return &controller
}
