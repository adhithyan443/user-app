package usecase

import (
	"user-app/internal/domain"
	"user-app/internal/repository/interfaces"
)

type AdminUsecase struct {
	userRepo interfaces.UserRepository
}

func NewAdminUsecase(
	userRepo interfaces.UserRepository,
) *AdminUsecase {

	return &AdminUsecase{
		userRepo: userRepo,
	}
}

func (u *AdminUsecase) GetAllUsers() (
	[]domain.User,
	error,
) {
	return u.userRepo.GetAll()
}