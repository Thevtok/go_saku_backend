package repository

import (
	"github.com/ReygaFitra/inc-final-project.git/model"
)

type UserRepository interface {
	FindByID(id uint) (*model.Users, error)
	FindByEmail(email string) (*model.Users, error)
	FindByPhone(phone string) (*model.Users, error)
	Create(user *model.Users) error
	Update(user *model.Users) error
	Delete(user *model.Users) error
}
