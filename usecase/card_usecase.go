package usecase

import (
	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/repository"
)

type CardUsecase interface {
	FindCardByID(id uint) (any, error)
	FindCardByCardID(id uint) (*model.Card, error)
	Register(userID uint, newCard *model.CardResponse) (any, error)
	Edit(card *model.Card) string
	UnregALL(card *model.Card) string
	UnregByCardId(cardID uint) error
}

type cardUsecase struct {
	cardRepo repository.CardRepository
}

func (u *cardUsecase) FindCardByID(id uint) (any, error) {
	return u.cardRepo.GetByID(id)
}

func (u *cardUsecase) FindCardByCardID(id uint) (*model.Card, error) {
	return u.cardRepo.GetByCardID(id)
}

func (u *cardUsecase) Register(userID uint, newCard *model.CardResponse) (any, error) {
	return u.cardRepo.Create(userID, newCard)
}

func (u *cardUsecase) Edit(card *model.Card) string {
	return u.cardRepo.Update(card)
}

func (u *cardUsecase) UnregALL(card *model.Card) string {
	return u.cardRepo.DeleteByUserID(card.UserID)
}

func (u *cardUsecase) UnregByCardId(cardID uint) error {
	err := u.cardRepo.DeleteByCardId(cardID)
	if err != nil {
		return err
	}
	return nil
}

func NewCardUsecase(cardRepo repository.CardRepository) CardUsecase {
	return &cardUsecase{
		cardRepo: cardRepo,
	}
}
