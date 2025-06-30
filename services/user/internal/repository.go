package user

import (
	"context"

	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	GetUserByID(ctx context.Context, id uint64) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateUser(ctx context.Context, id uint64, updates map[string]any) (*domain.User, error)
	DeleteUser(ctx context.Context, id uint64) (*domain.User, error)
}
