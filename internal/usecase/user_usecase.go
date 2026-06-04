package usecase

import (
	"errors"
	"regexp"
	"user-app/internal/domain"
	"user-app/internal/repository/interfaces"
)

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
) (*domain.User, error) {

	return u.userRepo.FindByID(userID)
}

func (u *UserUsecase) UpdateProfile(
	req UpdateProfileRequest,
) error {

	if req.Name == "" || req.Email == "" {
		return errors.New("all fields are required")
	}

	if len(req.Name) < 3 {
		return errors.New("name must be at least 3 characters")
	}

	nameRegex := regexp.MustCompile(`^[a-zA-Z ]+$`) 

	if !nameRegex.MatchString(req.Name) {
		return errors.New("name should contain only letters")
	}

	return u.userRepo.UpdateProfile(
		req.ID,
		req.Name,
		req.Email,
	)
}

func (u *UserUsecase) ChangePassword(
	userID uint,
	currentPassword string,
	newPassword string,
) error {

	return nil
}
