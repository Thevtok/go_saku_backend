package usecase

import (
	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/repository"
)

type CardUsecase interface {
	FindByID(id uint) (*model.Card, error)
	FindByUserID(userID uint) ([]*model.Card, error)
	Register(card *model.Card) error
	Edit(card *model.Card) error
	Unreg(card *model.Card) error
}

type cardUsecase struct {
	cardRepo repository.CardRepository
}

func (u *cardUsecase) FindByID(id uint) (*model.Card, error) {
	return u.cardRepo.GetByID(id)
}

func (u *cardUsecase) FindByUserID(userID uint) ([]*model.Card, error) {
	return u.cardRepo.GetByUserID(userID)
}

func (u *cardUsecase) Register(card *model.Card) error {
	return u.cardRepo.Create(card)
}

func (u *cardUsecase) Edit(card *model.Card) error {
	return u.cardRepo.Update(card)
}

func (u *cardUsecase) Unreg(card *model.Card) error {
	return u.cardRepo.Delete(card)
}

func NewCardUsecase(cardRepo repository.CardRepository) CardUsecase {
	return &cardUsecase{
		cardRepo: cardRepo,
	}
}
