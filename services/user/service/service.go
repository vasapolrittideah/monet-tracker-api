package service

import (
	"github.com/google/uuid"
	"vasapolrittideah/money-tracker-api/services/user/repository"
	"vasapolrittideah/money-tracker-api/shared/config"
	"vasapolrittideah/money-tracker-api/shared/domain/apperror"
	"vasapolrittideah/money-tracker-api/shared/domain/entity"
)

type UserService interface {
	GetAllUsers() ([]*entity.User, *apperror.Error)
	GetUserById(id uuid.UUID) (*entity.User, *apperror.Error)
	GetUserByEmail(email string) (*entity.User, *apperror.Error)
	CreateUser(user *entity.User) (*entity.User, *apperror.Error)
	UpdateUser(user *entity.User) (*entity.User, *apperror.Error)
	DeleteUser(id uuid.UUID) (*entity.User, *apperror.Error)
}

type userService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewUserService(userRepo repository.UserRepository, cfg *config.Config) UserService {
	return &userService{userRepo, cfg}
}

func (s *userService) GetAllUsers() ([]*entity.User, *apperror.Error) {
	return s.userRepo.GetAllUsers()
}

func (s *userService) GetUserById(id uuid.UUID) (*entity.User, *apperror.Error) {
	return s.userRepo.GetUserById(id)
}

func (s *userService) GetUserByEmail(email string) (*entity.User, *apperror.Error) {
	return s.userRepo.GetUserByEmail(email)
}

func (s *userService) CreateUser(user *entity.User) (*entity.User, *apperror.Error) {
	return s.userRepo.CreateUser(user)
}

func (s *userService) UpdateUser(user *entity.User) (*entity.User, *apperror.Error) {
	return s.userRepo.UpdateUser(user)
}

func (s *userService) DeleteUser(id uuid.UUID) (*entity.User, *apperror.Error) {
	return s.userRepo.DeleteUser(id)
}
