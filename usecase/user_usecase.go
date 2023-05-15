package usecase

import (
	"fmt"
	"log"

	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/repository"
	"github.com/ReygaFitra/inc-final-project.git/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	Login(email string, password string) (*model.Credentials, error)
	FindUsers() any
	FindByUsername(username string) (*model.UserResponse, error)
	FindById(id uint) (*model.User, error)
	Register(user *model.UserCreate) (any, error)
	EditProfile(user *model.User) string
	EditEmailPassword(user *model.User) string
	Unreg(user *model.User) string
	FindByPhone(phoneNumber string) (*model.User, error)
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func (uc *userUseCase) Login(email string, password string) (*model.Credentials, error) {
	// Get the user by email and hashed password
	user, err := uc.userRepo.GetByEmailAndPassword(email, password)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	// Compare the provided password with the stored password hash
	err = utils.CheckPasswordHash(password, user.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials \n password = %s\n hased = %s", password, user.Password)
	}
	return &model.Credentials{Password: user.Password, Username: user.Username, UserID: user.UserID, Role: user.Role}, nil
}

func (uc *userUseCase) FindUsers() any {
	return uc.userRepo.GetAll()
}

func (uc *userUseCase) FindByUsername(username string) (*model.UserResponse, error) {
	return uc.userRepo.GetByUsername(username)
}
func (uc *userUseCase) FindByPhone(phone string) (*model.User, error) {
	return uc.userRepo.GetByPhone(phone)
}

func (uc *userUseCase) FindById(id uint) (*model.User, error) {
	return uc.userRepo.GetByiD(id)
}

func (uc *userUseCase) Register(user *model.UserCreate) (any, error) {
	return uc.userRepo.Create(user)
}

func (uc *userUseCase) EditEmailPassword(user *model.User) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error generating hashed password: %v", err)
	}
	user.Password = string(hashedPassword)
	return uc.userRepo.UpdateEmailPassword(user)
}

func (uc *userUseCase) EditProfile(user *model.User) string {

	return uc.userRepo.UpdateProfile(user)
}

func (uc *userUseCase) Unreg(user *model.User) string {
	return uc.userRepo.Delete(user)
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}
