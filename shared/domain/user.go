package domain

import (
	"time"

	userv1 "github.com/vasapolrittideah/money-tracker-api/protogen/user/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type User struct {
	ID           uint64    `json:"id"         gorm:"primarykey;type:uuid;autoIncrement"`
	FullName     string    `json:"name"       gorm:"not null;type:varchar(100)"`
	Email        string    `json:"email"      gorm:"not null;uniqueIndex"`
	Verified     bool      `json:"verified"   gorm:"not null;default:false"`
	Password     string    `json:"-"          gorm:"not null"`
	RefreshToken string    `json:"-"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
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
