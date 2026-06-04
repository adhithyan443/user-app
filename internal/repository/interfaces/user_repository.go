package interfaces

import "user-app/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error

	FindByEmail(email string) (*domain.User, error)

	FindByID(id uint) (*domain.User, error)

	Update(user *domain.User) error

	Delete(id uint) error

	GetAll() ([]domain.User, error)
}
