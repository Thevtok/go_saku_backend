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
	Upload(ctx context.Context, username string, file multipart.File, header *multipart.FileHeader) error
	Download(username string) (*model.PhotoUrl, error)
	Edit(photo *model.PhotoUrl, username string, file multipart.File, header *multipart.FileHeader) error
	Remove(username string) string
}

type photoUsecase struct {
	photoRepo repository.PhotoRepository
}

func NewPhotoUseCase(photoRepo repository.PhotoRepository) PhotoUsecase {
	return &photoUsecase{
		photoRepo: photoRepo,
	}
}

func (u *photoUsecase) Upload(ctx context.Context, username string, file multipart.File, header *multipart.FileHeader) error {
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
        Url:    path,
        Username: username,
    }
    err = u.photoRepo.Create(photo)
    if err != nil {
        return err
    }
    return nil
}

func (u *photoUsecase) Download(username string) (*model.PhotoUrl, error) {
	return u.photoRepo.GetByID(username)
}

func (u *photoUsecase) Edit(photo *model.PhotoUrl, username string, file multipart.File, header *multipart.FileHeader) error {
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
        Url:    path,
        Username: username,
    }
    return u.photoRepo.Update(photoUrl)
}

func (u *photoUsecase) Remove(username string) string {
	return u.photoRepo.Delete(username)
}
