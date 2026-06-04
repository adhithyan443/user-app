package usecase

import (
	"errors"
	"regexp"
	"user-app/internal/domain"
	"user-app/internal/repository/interfaces"
	"user-app/internal/utils"

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

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	User  *domain.User
	Token string
}

type ResetPasswordRequest struct {
	UserID          uint
	NewPassword     string
	ConfirmPassword string
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

func (u *AuthUsecase) Login(
	req LoginRequest,
) (*LoginResponse, error) {

	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := u.userRepo.FindByEmail(req.Email)

	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(
		user.ID,
		user.Email,
		user.Role,
	)

	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (u *AuthUsecase) FindUserByEmail(
	email string,
) (*domain.User, error) {

	if email == "" {
		return nil, errors.New("email is required")
	}

	return u.userRepo.FindByEmail(email)
}


func (u *AuthUsecase) ResetPassword(
	req ResetPasswordRequest,
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

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	return u.userRepo.UpdatePassword(
		req.UserID,
		string(hashedPassword),
	)
}