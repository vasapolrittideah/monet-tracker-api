package usecase

import (
	"errors"
	"strings"

	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"github.com/vasapolrittideah/money-tracker-api/shared/errors/apperror"
	"github.com/vasapolrittideah/money-tracker-api/shared/utils/hashutil"
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
		return nil, apperror.NewError(apperror.ErrInternal, err.Error())
	}
	if len(users) == 0 {
		return nil, apperror.NewError(apperror.ErrNotFound, "no users found")
	}

	return users, nil
}

func (u *userUsecase) GetUserByID(id uint64) (*domain.User, error) {
	user, err := u.repository.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewError(apperror.ErrNotFound, "user not found")
		}

		return nil, apperror.NewError(apperror.ErrInternal, err.Error())
	}

	return user, nil
}

func (u *userUsecase) GetUserByEmail(email string) (*domain.User, error) {
	user, err := u.repository.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewError(apperror.ErrNotFound, "user not found")
		}

		return nil, apperror.NewError(apperror.ErrInternal, err.Error())
	}

	return user, nil
}

func (u *userUsecase) CreateUser(user *domain.User) (*domain.User, error) {
	hashedPassword, err := hashutil.Hash(user.Password)
	if err != nil {
		return nil, apperror.NewError(apperror.ErrInternal, err.Error())
	}

	user.Password = hashedPassword

	created, err := u.repository.CreateUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, apperror.NewError(apperror.ErrAlreadyExists, "user already exists")
		}

		return nil, apperror.NewError(apperror.ErrInternal, err.Error())
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
		return nil, apperror.NewError(apperror.ErrInvalidArgument, "no changes detected")
	}

	updated, err := u.repository.UpdateUser(user)
	if err != nil {
		return nil, apperror.NewError(apperror.ErrInternal, err.Error())
	}

	return updated, nil
}

func (u *userUsecase) DeleteUser(id uint64) (*domain.User, error) {
	deleted, err := u.repository.DeleteUser(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewError(apperror.ErrNotFound, "user not found")
		}

		return nil, apperror.NewError(apperror.ErrInternal, err.Error())
	}

	return deleted, nil
}
