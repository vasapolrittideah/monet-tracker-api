package usecase

import (
	"context"
	"errors"
	"strings"

	user "github.com/vasapolrittideah/money-tracker-api/services/user/internal"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"github.com/vasapolrittideah/money-tracker-api/shared/errors/apperror"
	"github.com/vasapolrittideah/money-tracker-api/shared/utils/hashutil"
	"gorm.io/gorm"
)

type userUsecase struct {
	repository user.UserRepository
	config     *config.Config
}

func NewUserUsecase(repository user.UserRepository, config *config.Config) user.UserUsecase {
	return &userUsecase{repository: repository, config: config}
}

func (u *userUsecase) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := u.repository.GetAllUsers(ctx)
	if err != nil {
		return nil, apperror.NewError(apperror.CodeInternal, "failed to get users")
	}
	if len(users) == 0 {
		return nil, apperror.NewError(apperror.CodeNotFound, "no users found")
	}

	return users, nil
}

func (u *userUsecase) GetUserByID(ctx context.Context, id uint64) (*domain.User, error) {
	user, err := u.repository.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewError(apperror.CodeNotFound, "user not found")
		}

		return nil, apperror.NewError(apperror.CodeInternal, "failed to get user")
	}

	return user, nil
}

func (u *userUsecase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.repository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewError(apperror.CodeNotFound, "user not found")
		}

		return nil, apperror.NewError(apperror.CodeInternal, "failed to get user")
	}

	return user, nil
}

func (u *userUsecase) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	hashedPassword, err := hashutil.Hash(user.Password)
	if err != nil {
		return nil, apperror.NewError(apperror.CodeInternal, "failed to hash password")
	}

	user.Password = hashedPassword

	created, err := u.repository.CreateUser(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, apperror.NewError(apperror.CodeAlreadyExists, "user already exists")
		}

		return nil, apperror.NewError(apperror.CodeInternal, "failed to create user")
	}

	return created, nil
}

func (u *userUsecase) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	existing, err := u.GetUserByID(ctx, user.ID)
	if err != nil {
		return nil, apperror.NewError(apperror.CodeInternal, "failed to get user")
	}

	if existing.FullName == user.FullName &&
		existing.Email == user.Email &&
		existing.Password == user.Password &&
		existing.Verified == user.Verified &&
		existing.Registered == user.Registered {
		return nil, apperror.NewError(apperror.CodeInvalidArgument, "no changes detected")
	}

	updated, err := u.repository.UpdateUser(ctx, user)
	if err != nil {
		return nil, apperror.NewError(apperror.CodeInternal, "failed to update user")
	}

	return updated, nil
}

func (u *userUsecase) DeleteUser(ctx context.Context, id uint64) (*domain.User, error) {
	deleted, err := u.repository.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewError(apperror.CodeNotFound, "user not found")
		}

		return nil, apperror.NewError(apperror.CodeInternal, "failed to delete user")
	}

	return deleted, nil
}
