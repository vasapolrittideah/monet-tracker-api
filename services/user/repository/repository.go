package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"vasapolrittideah/money-tracker-api/shared/domain/apperror"
	"vasapolrittideah/money-tracker-api/shared/domain/entity"
	"vasapolrittideah/money-tracker-api/shared/utils/errorutil"
)

type UserRepository interface {
	GetAllUsers() ([]*entity.User, *apperror.Error)
	GetUserById(id uuid.UUID) (*entity.User, *apperror.Error)
	GetUserByEmail(email string) (*entity.User, *apperror.Error)
	CreateUser(user *entity.User) (*entity.User, *apperror.Error)
	UpdateUser(user *entity.User) (*entity.User, *apperror.Error)
	DeleteUser(id uuid.UUID) (*entity.User, *apperror.Error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetAllUsers() ([]*entity.User, *apperror.Error) {
	var users []*entity.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, errorutil.HandleRecordNotFoundError(err)
	}

	return users, nil
}

func (r *userRepository) GetUserById(id uuid.UUID) (*entity.User, *apperror.Error) {
	var user *entity.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, errorutil.HandleRecordNotFoundError(err)
	}

	return user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*entity.User, *apperror.Error) {
	var user *entity.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, errorutil.HandleRecordNotFoundError(err)
	}

	return user, nil
}

func (r *userRepository) CreateUser(user *entity.User) (*entity.User, *apperror.Error) {
	if err := r.db.Create(&user).Error; err != nil {
		return nil, errorutil.HandleUnqiueConstraintError(err)
	}

	return user, nil
}

func (r *userRepository) UpdateUser(newUserData *entity.User) (*entity.User, *apperror.Error) {
	var user *entity.User
	if err := r.db.First(&user, "id = ?", newUserData.Id).Error; err != nil {
		return nil, errorutil.HandleRecordNotFoundError(err)
	}

	if err := r.db.Model(&user).Updates(newUserData.Id).Error; err != nil {
		return nil, errorutil.HandleUnknownDatabaseError(err)
	}

	return user, nil
}

func (r *userRepository) DeleteUser(id uuid.UUID) (*entity.User, *apperror.Error) {
	var user *entity.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, errorutil.HandleRecordNotFoundError(err)
	}

	if err := r.db.Delete(&user).Error; err != nil {
		return nil, errorutil.HandleUnknownDatabaseError(err)
	}

	return user, nil
}
