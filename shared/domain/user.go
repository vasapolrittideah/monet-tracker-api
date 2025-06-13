package domain

import (
	"time"

	userv1 "github.com/vasapolrittideah/money-tracker-api/protogen/user/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type User struct {
	ID           uint64    `json:"id"         gorm:"primaryKey;autoIncrement"   example:"1"                    extensions:"x-order=1"`
	FullName     string    `json:"full_name"  gorm:"not null;type:varchar(100)" example:"John Doe"             extensions:"x-order=2"`
	Email        string    `json:"email"      gorm:"not null;uniqueIndex"       example:"john@example.com"     extensions:"x-order=3"`
	Verified     bool      `json:"verified"   gorm:"not null;default:false"     example:"true"                 extensions:"x-order=4"`
	Password     string    `json:"-"          gorm:"not null"`
	RefreshToken string    `json:"-"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"             example:"2022-01-01T00:00:00Z" extensions:"x-order=5"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"             example:"2022-01-01T00:00:00Z" extensions:"x-order=6"`
}

func (u *User) ToProto() *userv1.User {
	return &userv1.User{
		Id:        u.ID,
		FullName:  u.FullName,
		Email:     u.Email,
		Verified:  u.Verified,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}

func NewUserFromProto(user *userv1.User) *User {
	return &User{
		ID:        user.Id,
		FullName:  user.FullName,
		Email:     user.Email,
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt.AsTime(),
		UpdatedAt: user.UpdatedAt.AsTime(),
	}
}

type CreateUserRequest struct {
	FullName string `json:"full_name" example:"John Doe"         extensions:"x-order=1"`
	Email    string `json:"email"     example:"john@example.com" extensions:"x-order=2"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	FullName *string `json:"full_name,omitempty" example:"John Doe"         extensions:"x-order=1"`
	Email    *string `json:"email,omitempty"     example:"john@example.com" extensions:"x-order=2"`
	Verified *bool   `json:"verified,omitempty"  example:"true"             extensions:"x-order=3"`
} // @ignore

type UserRepository interface {
	GetAllUsers() ([]*User, error)
	GetUserByID(id uint64) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(user *User) (*User, error)
	UpdateUser(id uint64, user *User) (*User, error)
	DeleteUser(id uint64) (*User, error)
}

type UserUsecase interface {
	GetAllUsers() ([]*User, error)
	GetUserByID(id uint64) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(user *User) (*User, error)
	UpdateUser(id uint64, user *User) (*User, error)
	DeleteUser(id uint64) (*User, error)
}
