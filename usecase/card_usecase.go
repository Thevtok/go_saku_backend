package usecase

import (
	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/repository"
)

type CardUsecase interface {
	FindAllCard() any
	FindCardByUserID(id uint) ([]*model.CardResponse, error)
	FindCardByCardID(id uint) (*model.Card, error)
	Register(id uint, newCard *model.CreateCard) (any, error)
	Edit(card *model.Card) string
	UnregALL(userID uint) string
	UnregByCardID(cardID uint) error
}

type cardUsecase struct {
	cardRepo repository.CardRepository
}

func (u *cardUsecase) FindAllCard() any {
	return u.cardRepo.GetAll()
}

func (u *cardUsecase) FindCardByUserID(id uint) ([]*model.CardResponse, error) {
	return u.cardRepo.GetByUserID(id)
}

func (u *cardUsecase) FindCardByCardID(id uint) (*model.Card, error) {
	return u.cardRepo.GetByCardID(id)
}

func (u *cardUsecase) Register(id uint, newCard *model.CreateCard) (any, error) {
	newCard.UserID = id
	return u.cardRepo.Create(id, newCard)
}

func (u *cardUsecase) Edit(card *model.Card) string {
	return u.cardRepo.Update(card)
}

func (u *cardUsecase) UnregALL(userID uint) string {
	return u.cardRepo.DeleteByUserID(userID)
}

func (u *cardUsecase) UnregByCardID(cardID uint) error {
	err := u.cardRepo.DeleteByCardID(cardID)
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
