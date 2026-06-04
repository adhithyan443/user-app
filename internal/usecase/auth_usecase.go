package usecase

import "user-app/internal/repository/interfaces"

type AuthUsecase struct {
	userRepo interfaces.UserRepository
}

func NewAuthUsecase(
	userRepo interfaces.UserRepository,
) *AuthUsecase {

	return &AuthUsecase{
		userRepo: userRepo,
	}
}

type SignupRequest struct {
	Name     string
	Email    string
	Password string
}

func (u *AuthUsecase) Signup(
	req SignupRequest,
) error {

	return nil
}