package repository

import (
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetAllUsers() ([]*domain.User, error) {
	var users []*domain.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetUserByID(id uint64) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CreateUser(user *domain.User) (*domain.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) UpdateUser(id uint64, updated *domain.User) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&user).Where("id = ?", id).Save(updated).Error; err != nil {
		return nil, err
	}
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) DeleteUser(id uint64) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := r.db.Delete(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
