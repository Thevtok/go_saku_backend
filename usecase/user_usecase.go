package usecase

import (
	"errors"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	Login(email string, password string) (*model.User, error)
	FindUsers() any
	FindByID(id uint) any

	Register(user *model.User) string
	Edit(user *model.User) string
	Unreg(id uint) string
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (uc *userUseCase) Login(email string, password string) (*model.User, error) {
	// Hash the provided password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Get the user by email and hashed password
	user, err := uc.userRepo.GetByUsernameAndPassword(email, string(hashedPassword))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Return the user if the passwords match
	return user, nil
}

func (uc *userUseCase) FindUsers() any {
	return uc.userRepo.GetAll()
}

func (uc *userUseCase) FindByID(userID uint) any {
	return uc.userRepo.GetByID(userID)
}

func (uc *userUseCase) Register(user *model.User) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		print(err)
	}
	user.Password = string(hashedPassword)
	return uc.userRepo.Create(user)
}

func (uc *userUseCase) Edit(user *model.User) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		print(err)
	}
	user.Password = string(hashedPassword)
	return uc.userRepo.Update(user)
}

func (uc *userUseCase) Unreg(userID uint) string {
	return uc.userRepo.Delete(userID)
}
