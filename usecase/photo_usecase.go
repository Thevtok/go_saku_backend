package usecase

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/repository"
	"github.com/ReygaFitra/inc-final-project.git/utils"
)

type PhotoUsecase interface {
	Upload(ctx context.Context, userID uint, file multipart.File, header *multipart.FileHeader) error
	Download(id uint) (*model.PhotoUrl, error)
	Edit(photo *model.PhotoUrl, id uint, file multipart.File, header *multipart.FileHeader) error
	Remove(id uint) string
}

type photoUsecase struct {
	photoRepo repository.PhotoRepository
}

func NewPhotoUseCase(photoRepo repository.PhotoRepository) PhotoUsecase {
	return &photoUsecase{
		photoRepo: photoRepo,
	}
}

func (u *photoUsecase) Upload(ctx context.Context, userID uint, file multipart.File, header *multipart.FileHeader) error {
	// Simpan file ke server
    filename := header.Filename
    path := fmt.Sprintf(utils.DotEnv("FILE_LOCATION"), filename)
    out, err := os.Create(path)
    if err != nil {
        return err
    }
    defer out.Close()
    _, err = io.Copy(out, file)
    if err != nil {
        return err
    }
    // Simpan informasi file ke database
    photo := &model.PhotoUrl{
        User_ID: userID,
        Url:    path,
    }
    err = u.photoRepo.Create(photo)
    if err != nil {
        return err
    }
    return nil
}

func (u *photoUsecase) Download(id uint) (*model.PhotoUrl, error) {
	return u.photoRepo.GetByID(id)
}

func (u *photoUsecase) Edit(photo *model.PhotoUrl, id uint, file multipart.File, header *multipart.FileHeader) error {
	// Simpan file ke server
    filename := header.Filename
    path := fmt.Sprintf(utils.DotEnv("FILE_LOCATION"), filename)
    out, err := os.Create(path)
    if err != nil {
        return err
    }
    defer out.Close()
    _, err = io.Copy(out, file)
    if err != nil {
        return err
    }
    // Simpan informasi file ke database
    photoUrl := &model.PhotoUrl{
        User_ID: id,
        Url:    path,
    }
    return u.photoRepo.Update(photoUrl)
}

func (u *photoUsecase) Remove(id uint) string {
	return u.photoRepo.Delete(id)
}
