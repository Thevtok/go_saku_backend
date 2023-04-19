package usecase

import (
	"github.com/ReygaFitra/inc-final-project.git/model"
	"github.com/ReygaFitra/inc-final-project.git/repository"
)

type RewardUseCase interface {
	FindPoints() any
	FindPointByID(id uint) any
	Register(newPoint *model.Reward) string
	// Edit(point *model.Reward) string
	Unreg(id uint) string
}

type rewardUsecase struct {
	rewardRepo repository.RewardRepository
}

func (u *rewardUsecase) FindPoints() any {
	return u.rewardRepo.GetAll()
}

func (u *rewardUsecase) FindPointByID(id uint) any {
	return u.rewardRepo.GetByID(id)
}

func (u *rewardUsecase) Register(newPoint *model.Reward) string {
	return u.rewardRepo.Create(newPoint)
}

// func (u *rewardUsecase) Edit(point *model.Reward) string {
// 	return u.rewardRepo.Update(point)
// }

func (u *rewardUsecase) Unreg(id uint) string {
	return u.rewardRepo.Delete(id)
}

func NewRewardUseCase(rewardRepo repository.RewardRepository) RewardUseCase {
	return &rewardUsecase{
		rewardRepo: rewardRepo,
	}
}