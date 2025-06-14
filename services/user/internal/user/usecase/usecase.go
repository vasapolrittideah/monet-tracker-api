package usecase

import (
	"errors"
	"strings"

	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"github.com/vasapolrittideah/money-tracker-api/shared/utils/hashutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type userUsecase struct {
	repository domain.UserRepository
	config     *config.Config
}

func NewUserUsecase(repository domain.UserRepository, config *config.Config) domain.UserUsecase {
	return &userUsecase{repository: repository, config: config}
}

func (u *userUsecase) GetAllUsers() ([]*domain.User, error) {
	users, err := u.repository.GetAllUsers()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}
	if len(users) == 0 {
		return nil, status.Errorf(codes.NotFound, "no users found")
	}

	return users, nil
}

func (u *userUsecase) GetUserByID(id uint64) (*domain.User, error) {
	user, err := u.repository.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}

		return nil, status.Errorf(codes.Internal, "failed to get user: %s", err)
	}

	return user, nil
}

func (u *userUsecase) GetUserByEmail(email string) (*domain.User, error) {
	user, err := u.repository.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}

		return nil, status.Errorf(codes.Internal, "failed to get user: %s", err)
	}

	return user, nil
}

func (u *userUsecase) CreateUser(user *domain.User) (*domain.User, error) {
	hashedPassword, err := hashutil.Hash(user.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	user.Password = hashedPassword

	created, err := u.repository.CreateUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return created, nil
}

func (u *userUsecase) UpdateUser(user *domain.User) (*domain.User, error) {
	existing, err := u.GetUserByID(user.ID)
	if err != nil {
		return nil, err
	}

	if existing.FullName == user.FullName &&
		existing.Email == user.Email &&
		existing.Verified == user.Verified &&
		existing.Password == user.Password &&
		existing.RefreshToken == user.RefreshToken {
		return nil, status.Errorf(codes.InvalidArgument, "no changes detected")
	}

	updated, err := u.repository.UpdateUser(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return updated, nil
}

func (u *userUsecase) DeleteUser(id uint64) (*domain.User, error) {
	deleted, err := u.repository.DeleteUser(id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return deleted, nil
}
