package usecase

import (
	"errors"
	"regexp"
	"user-app/internal/domain"
	"user-app/internal/repository/interfaces"
	"user-app/utils"

	"golang.org/x/crypto/bcrypt"
)

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



func (u *AuthUsecase) Signup(req SignupRequest) error {

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return errors.New("all fields are required")
		
	}

	if len(req.Name) < 3 {
		return errors.New("name must be at least 3 characters")
	}

	var nameRegex = regexp.MustCompile(`^[a-zA-Z ]+$`)

	if !nameRegex.MatchString(req.Name) {
		return errors.New("name should contain only letters")
	}

	if !utils.IsStrongPassword(req.Password) {
		
		return errors.New(
			"password must contain uppercase, lowercase, number, and special character",
		)
	}

	_, err := u.userRepo.FindByEmail(req.Email)

	if err == nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	user := domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	return u.userRepo.Create(&user)
}
