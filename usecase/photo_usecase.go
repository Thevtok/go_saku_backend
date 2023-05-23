package usecase

import (
	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/repository"
)

type PhotoUsecase interface {
	Upload(photo *model.PhotoUrl) error
	Download(id string) (*model.PhotoUrl, error)
	Edit(photo *model.PhotoUrl) error
	Remove(id string) string
}

type photoUsecase struct {
	photoRepo repository.PhotoRepository
}

func (u *photoUsecase) Upload(photo *model.PhotoUrl) error {
	return u.photoRepo.Create(photo)
}

func (u *photoUsecase) Download(id string) (*model.PhotoUrl, error) {
	return u.photoRepo.GetByID(id)
}

func (u *photoUsecase) Edit(photo *model.PhotoUrl) error {
	return u.photoRepo.Update(photo)
}

func (u *photoUsecase) Remove(id string) string {
	return u.photoRepo.Delete(id)
}

func NewPhotoUseCase(photoRepo repository.PhotoRepository) PhotoUsecase {
	return &photoUsecase{
		photoRepo: photoRepo,
	}
}
