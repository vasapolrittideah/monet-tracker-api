package service

import (
	"github.com/google/uuid"
	"github.com/vasapolrittideah/money-tracker-api/services/user/repository"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/model/apperror"
	"github.com/vasapolrittideah/money-tracker-api/shared/model/domain"
)

type UserService interface {
	GetAllUsers() ([]*domain.User, *apperror.Error)
	GetUserById(id uuid.UUID) (*domain.User, *apperror.Error)
	GetUserByEmail(email string) (*domain.User, *apperror.Error)
	CreateUser(user *domain.User) (*domain.User, *apperror.Error)
	UpdateUser(id uuid.UUID, newUserData *domain.User) (*domain.User, *apperror.Error)
	DeleteUser(id uuid.UUID) (*domain.User, *apperror.Error)
}

type userService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewUserService(userRepo repository.UserRepository, cfg *config.Config) UserService {
	return &userService{userRepo, cfg}
}

func (s *userService) GetAllUsers() ([]*domain.User, *apperror.Error) {
	return s.userRepo.GetAllUsers()
}

func (s *userService) GetUserById(id uuid.UUID) (*domain.User, *apperror.Error) {
	return s.userRepo.GetUserById(id)
}

func (s *userService) GetUserByEmail(email string) (*domain.User, *apperror.Error) {
	return s.userRepo.GetUserByEmail(email)
}

func (s *userService) CreateUser(user *domain.User) (*domain.User, *apperror.Error) {
	return s.userRepo.CreateUser(user)
}

func (s *userService) UpdateUser(id uuid.UUID, newUserData *domain.User) (*domain.User, *apperror.Error) {
	return s.userRepo.UpdateUser(id, newUserData)
}

func (s *userService) DeleteUser(id uuid.UUID) (*domain.User, *apperror.Error) {
	return s.userRepo.DeleteUser(id)
}
