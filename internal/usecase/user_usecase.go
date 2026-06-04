package usecase

import (
	"errors"
	"regexp"
	"user-app/internal/domain"
	"user-app/internal/repository/interfaces"
	"user-app/utils"

	"golang.org/x/crypto/bcrypt"
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

type ChangePasswordRequest struct {
	ID              uint
	OldPassword     string
	NewPassword     string
	ConfirmPassword string
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
	req ChangePasswordRequest,
) error {

	if !utils.IsStrongPassword(req.NewPassword) { 
		return errors.New(
			"password must contain uppercase, lowercase, number, and special character",
		)
	}

	if req.NewPassword != req.ConfirmPassword {
		return errors.New(
			"passwords do not match",
		)
	}

	user, err := u.userRepo.FindByID(req.ID)

	if err != nil {
		return errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.OldPassword),
	)

	if err != nil {
		return errors.New(
			"current password is incorrect",
		)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	return u.userRepo.UpdatePassword(
		req.ID,
		string(hashedPassword),
	)
}
