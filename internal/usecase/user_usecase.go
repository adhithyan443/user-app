package usecase

import "user-app/internal/repository/interfaces"

type UserUsecase struct {
	userRepo interfaces.UserRepository
}

func NewUserUsecase(
	userRepo interfaces.UserRepository,
) *UserUsecase {

	return &UserUsecase{
		userRepo: userRepo,
	}
}

type UpdateProfileRequest struct {
	ID    uint
	Name  string
	Email string
}

func (u *UserUsecase) GetProfile(
	userID uint,
) error {

	return nil
}

func (u *UserUsecase) UpdateProfile(
	req UpdateProfileRequest,
) error {

	return nil
}

func (u *UserUsecase) ChangePassword(
	userID uint,
	currentPassword string,
	newPassword string,
) error {

	return nil
}