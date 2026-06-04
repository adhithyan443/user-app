package usecase

import (
	"errors"
	"regexp"
	"user-app/internal/domain"
	"user-app/internal/repository/interfaces"
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