package usecase

import (
	"errors"
	"regexp"
	"user-app/internal/domain"
	"user-app/internal/repository/interfaces"
	"user-app/utils"

	"golang.org/x/crypto/bcrypt"
)

type AdminUsecase struct {
	userRepo interfaces.UserRepository
}

type UpdateUserRequest struct {
	ID    uint
	Name  string
	Email string
	Role  string
}

type CreateUserRequest struct {
	Name     string
	Email    string
	Password string
	Role     string
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

func (u *AdminUsecase) DeleteUser(
	id uint,
) error {

	return u.userRepo.Delete(id)
}

func (u *AdminUsecase) GetUserByID(
	id uint,
) (*domain.User, error) {

	return u.userRepo.FindByID(id)
}

func (u *AdminUsecase) UpdateUser(
	req UpdateUserRequest,
) error {

	if req.Name == "" ||
		req.Email == "" ||
		req.Role == "" {

		return errors.New(
			"all fields are required",
		)
	}

	if len(req.Name) < 3 {
		return errors.New(
			"name must be at least 3 characters",
		)
	}

	nameRegex := regexp.MustCompile(
		`^[a-zA-Z ]+$`,
	)

	if !nameRegex.MatchString(req.Name) {
		return errors.New(
			"name should contain only letters",
		)
	}

	if req.Role != "admin" &&
		req.Role != "user" {

		return errors.New(
			"invalid role selected",
		)
	}

	user, err := u.userRepo.FindByID(
		req.ID,
	)

	if err != nil {
		return err
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Role = req.Role

	return u.userRepo.Update(user)
}

func (u *AdminUsecase) CreateUser(req CreateUserRequest) error {

	if req.Name == "" || req.Email == "" || req.Password == "" || req.Role == "" {
		return errors.New("all fields are required")
	}

	if len(req.Name) < 3 {
		return errors.New("name must be at least 3 characters")
	}

	nameRegex := regexp.MustCompile(`^[a-zA-Z ]+$`)
	if !nameRegex.MatchString(req.Name) {
		return errors.New("name should contain only letters")
	}

	if !utils.IsStrongPassword(req.Password) {
		return errors.New("password must contain uppercase, lowercase, number, and special character")
	}

	if req.Role != "admin" && req.Role != "user" {
		return errors.New("invalid role selected")
	}

	// Check if email already exists
	_, err := u.userRepo.FindByEmail(req.Email)
	if err == nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	return u.userRepo.Create(user)
}


// UpdateUserPassword - Admin changes any user's password
func (u *AdminUsecase) UpdateUserPassword(id uint, newPassword string) error {

	if !utils.IsStrongPassword(newPassword) {
		return errors.New("password must contain uppercase, lowercase, number, and special character")
	}

	if len(newPassword) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return u.userRepo.UpdatePassword(id, string(hashedPassword))
}

func (u *AdminUsecase) GetUserCount() (
	int64,
	error,
) {
	return u.userRepo.CountUsers()
}