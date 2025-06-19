package repository

import (
	"context"

	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uint64) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, updated *domain.User) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", updated.ID).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&user).Where("id = ?", updated.ID).Save(updated).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).First(&user, "id = ?", updated.ID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id uint64) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Delete(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
