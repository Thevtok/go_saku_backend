package usecase

import (
	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/repository"
)

type CardUsecase interface {
	FindAllCard() any
	FindCardByUsername(username string) (any, error)
	FindCardByCardID(id uint) (*model.Card, error)
	Register(username string, newCard *model.CardResponse) (any, error)
	Edit(card *model.Card) string
	UnregALL(card *model.Card) string
	UnregByCardId(cardID uint) error
}

type cardUsecase struct {
	cardRepo repository.CardRepository
}

func (u *cardUsecase) FindAllCard() any {
	return u.cardRepo.GetAll()
}

func (u *cardUsecase) FindCardByUsername(username string) (any, error) {
	return u.cardRepo.GetByUsername(username)
}

func (u *cardUsecase) FindCardByCardID(id uint) (*model.Card, error) {
	return u.cardRepo.GetByCardID(id)
}

func (u *cardUsecase) Register(username string, newCard *model.CardResponse) (any, error) {
	return u.cardRepo.Create(username, newCard)
}

func (u *cardUsecase) Edit(card *model.Card) string {
	return u.cardRepo.Update(card)
}

func (u *cardUsecase) UnregALL(card *model.Card) string {
	return u.cardRepo.DeleteByUsername(card.Username)
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
